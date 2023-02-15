package main

import (
	"archive/zip"
	"flag"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

var language = "dotnet"
var missingDebugFolder bool = false

// reused some code from https://github.com/fw10/veracode-javascript-packager/
// Only tested on a specific .NET project and not .NET solution 
func main() {
	// parse all the command line flags
	sourcePtr := flag.String("source", "", "The path of the .NET project app you want to package")
	targetPtr := flag.String("target", ".", "The path where you want the vc-output.zip to be stored to")

	flag.Parse()

	log.Info("#################################################")
	log.Info("#                                               #")
	log.Info("#   Veracode ", language, " Packager (Unofficial)       #")
	log.Info("#                                               #")
	log.Info("#################################################" + "\n\n")

	var binPath string
	var testsPath string

	// add the current date to the output zip name, like e.g. "2023-Jan-04"
	currentTime := time.Now()
	outputZipPath := filepath.Join(*targetPtr, "vc-output_"+currentTime.Format("2006-Jan-02")+".zip")

	// echo the provided flags
	log.Info("Provided Flags:")
	log.Info("\t`-source` directory to zip up: ", *sourcePtr)
	log.Info("\t`-target` directory for the output: ", *targetPtr)

	log.Info("Checking for 'smells' that indicate packaging issues - Started...")
	checkForPotentialSmells(*sourcePtr)
	log.Info("'Smells' Check - Done\n\n")

	debugPath, dotnetVersion := checkForDotNetVersion(*sourcePtr)
	log.Info("Found Debug Path = ", debugPath)
	log.Info("Found .NET version = ", dotnetVersion)
	
	binPath = debugPath

	publishFolderPath := checkForPublishFolder(debugPath, dotnetVersion)
	if publishFolderPath != ""{
		binPath = publishFolderPath
		log.Info("Found .NET Publish Folder ", publishFolderPath)
	}

	if dotnetVersion != "" {
		binPath = debugPath + string(os.PathSeparator) + dotnetVersion
	}

	log.Info("Creating a Zip while omitting non-required files - Started...")
	if err := zipSource(binPath, outputZipPath, testsPath); err != nil {
		log.Error(err)
	}

	log.Info("Zip Process - Done")
	log.Info("Wrote archive to: ", outputZipPath)
	log.Info("Please upload this archive to the Veracode Platform")
}

// default to the latest dotnet version folder found
func checkForDotNetVersion(source string) (debugPath string, version string) {

	// https://learn.microsoft.com/en-us/dotnet/standard/frameworks
	var dotnet_version_list = [35]string{"net7.0","net6.0","net5.0","netcoreapp3.1","netcoreapp3.0","netcoreapp2.2","netcoreapp2.1","netcoreapp2.0","netcoreapp1.1","netcoreapp1.0","netcoreapp1.0","netstandard2.1","netstandard2.0","netstandard1.6","netstandard1.5","netstandard1.4","netstandard1.3","netstandard1.2","netstandard1.1","netstandard1.0","net48","net472","net471","net47","net462","net461","net46","net452","net451","net45","net403","net40","net35","net20","net11"}

	err := filepath.WalkDir(source, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		if d.IsDir() && strings.HasSuffix(path, "bin" + string(os.PathSeparator) + "Debug"){
			debugPath = path
		}
		
		if debugPath != "" {
			for k := range dotnet_version_list {
				if d.IsDir() && strings.HasSuffix(path, debugPath + string(os.PathSeparator) + dotnet_version_list[k]) {
					version = dotnet_version_list[k]
					break
				}
			}
		}

		return nil
	})

	if err != nil {
		log.Error(err)
	} 

	return debugPath, version
}

func checkForPublishFolder(source string, version string) string {
	var publishPath string
	err := filepath.WalkDir(source, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		if version != "" {
			if d.IsDir() && strings.HasSuffix(path, string(os.PathSeparator) + version + string(os.PathSeparator) + "publish") {
				publishPath = path
			}
		} else {
			if d.IsDir() && strings.HasSuffix(path, string(os.PathSeparator) + "publish") {
				publishPath = path
			}
		}
		return nil
	})

	if err != nil {
		log.Error(err)
	} 

	return publishPath

}

func checkForPotentialSmells(source string) {

	err := filepath.WalkDir(source, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		if d.IsDir() && strings.HasSuffix(path, string(os.PathSeparator) + "bin" + string(os.PathSeparator) + "Debug") {
			missingDebugFolder = false
		}

		return nil
	})

	if err != nil {
		log.Error(err)
	}

	if missingDebugFolder == true {
		log.Info("Debug folder not found!")
	}

}

func zipSource(source string, target string, testsPath string) error {
	// 1. Create a ZIP file and zip.Writer
	f, err := os.Create(target)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := zip.NewWriter(f)
	defer writer.Close()

	// 2. Go through all the files of the source
	return filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

  	// avoids processing the created zip...
		// 	- Say the tool is finished and an `/vc-output_2023-Jan-05.zip` is created...
		//  - In this case, the analysis may restart with this zip as `path`
		// 		- This edge case was observed when running the tool within a sample JS app..
		//		- ... i.e., `veracode-js-packager -source . -target .`
		if strings.HasSuffix(path, ".zip") {
			return nil
		}
	
		// 3. Create a local file header
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// set compression
		header.Method = zip.Deflate

		// 4. Set relative path of a file as the header name
		// 	-> We want the following:
		//		- Say `-source some/path/my-js-project` is provided...
		//			- Now, say we have a path `some/path/my-js-project/build/some.js`....
		//		- In this scenario, we want `header.Name` to be `build/some.js`
		header.Name, err = filepath.Rel(source, path)
		if err != nil {
			return err
		}

		// avoids the `./` folder in the root of the output zip
		if header.Name == "." {
			return nil
		}

		// prepends the `/` we want before e.g. `build/some.js`
		headerNameWithSlash := string(os.PathSeparator) + header.Name

		// check if the path is required for the upload (otherwise, it will be omitted)
		if !isRequired(headerNameWithSlash, testsPath) {
			return nil
		}

		if info.IsDir() {
			// add e.g. a `/` if the current path is a directory
			header.Name += string(os.PathSeparator)
		}

		// 5. Create writer for the file header and save content of the file
		headerWriter, err := writer.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(headerWriter, f)
		return err
	})
}

func isRequired(path string, testsPath string) bool {
	return !IsRoslynFolder(path) &&
		!IsRuntimeFolder(path) &&
		!IsRuntimeIdentifierFolder(path) &&
		!IsWebRootFolder(path) &&
		!IsSatelliteLanguageFolder(path) &&
		!IsImage(path) &&
		!IsDocument(path) &&
		!IsVideo(path) &&
		!IsFont(path) &&
		!IsNestedArchive(path)
}
