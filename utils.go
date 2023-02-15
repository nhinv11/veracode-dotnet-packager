package main

import (
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

// reused some code from https://github.com/fw10/veracode-javascript-packager/

var didPrintRoslynFolderMsg bool = false
var didPrintRuntimeFolderMsg bool = false
var didPrintRIDFolderMsg bool = false
var didPrintWebRootFolderMsg bool = false
var didPrintSatelliteLanguageFolderMsg bool = false
var didPrintImagesMsg bool = false
var didPrintDocumentsMsg bool = false
var didPrintVideoMsg bool = false
var didPrintFontsMsg bool = false

func IsRoslynFolder(path string) bool {
	if strings.Contains(path, string(os.PathSeparator)+"Roslyn") {
		if !didPrintRoslynFolderMsg {
			log.Info("\tIgnoring the entire `Roslyn` folder")
			didPrintRoslynFolderMsg = true
		}
		return true
	}

	return false
}

func IsRuntimeFolder(path string) bool {
	if strings.Contains(path, string(os.PathSeparator)+"runtimes") {
		if !didPrintRuntimeFolderMsg {
			log.Info("\tIgnoring the entire `Runtimes` folder")
			didPrintRuntimeFolderMsg = true
		}
		return true
	}

	return false
}

func IsRuntimeIdentifierFolder(path string) bool {
	// Might Need to add more
	// https://learn.microsoft.com/en-us/dotnet/core/rid-catalog
	var rid_folders = [2]string{"win-x64","win-x86"}

	for _, rid_folder := range rid_folders {
		ridFolderPath :=  string(os.PathSeparator) + rid_folder
		fileInRIDFolderPath :=  ridFolderPath + string(os.PathSeparator)

		if strings.HasSuffix(path, ridFolderPath) || strings.Contains(path, fileInRIDFolderPath) {
			if !didPrintRIDFolderMsg {
				log.Info("\tIgnoring the `language` folders")
				didPrintRIDFolderMsg = true
			}
			return true
		}
	}

	return false
}

func IsWebRootFolder(path string) bool {
	if strings.Contains(path, string(os.PathSeparator)+"wwwroot") {
		if !didPrintWebRootFolderMsg {
			log.Info("\tIgnoring the entire `wwwroot` folder")
			didPrintWebRootFolderMsg = true
		}
		return true
	}

	return false
}

func IsSatelliteLanguageFolder(path string) bool {
	// Might Need to add more
	// https://lonewolfonline.net/list-net-culture-country-codes/
	languageFolders := [25]string{"cs","da","de","es","es-MX","fa","fi","fr","it","ja","ko","nb","nl","pl","pt","pt-BR","ro","ru","ru-ru","sl","sv","tr","uk","zh-Hans","zh-Hant"}

	for _, languageFolder := range languageFolders {
		languageFolderPath :=  string(os.PathSeparator) + languageFolder
		fileInLanguageFolderPath :=  languageFolderPath + string(os.PathSeparator)

		if strings.HasSuffix(path, languageFolderPath) || strings.Contains(path, fileInLanguageFolderPath) {
			if !didPrintSatelliteLanguageFolderMsg {
				log.Info("\tIgnoring the `language` folders")
				didPrintSatelliteLanguageFolderMsg = true
			}
			return true
		}
	}

	return false
}

// check for images (like .jpg, .png, .jpeg)
func IsImage(path string) bool {
	imageExtensions := [8]string{".jpg", ".png", ".jpeg", ".gif", ".svg", ".bmp", ".ico", ".icns"}

	for _, element := range imageExtensions {
		if strings.HasSuffix(path, element) {
			if !didPrintImagesMsg {
				log.Info("\tIgnoring images (such as `.jpg`)")
				didPrintImagesMsg = true
			}

			return true
		}
	}

	return false
}

// check for documents (like .pdf, .md)
func IsDocument(path string) bool {
	// inspired by https://en.wikipedia.org/wiki/List_of_Microsoft_Office_filename_extensions (and additionally `.md`)
	documentExtensions := [39]string{
		".pdf",
		".md",
		".doc", ".dot", ".wbk", ".docx", ".docm", ".dotx", ".dotm", ".docb", ".wll", ".wwl",
		".xls", ".xlt", ".xlm", ".xll_", ".xla_", ".xla5", ".xla8",
		".xlsx", ".xlsm", ".xltx", ".xltm",
		".ppt", ".pot", ".pps", ".pptx", ".pptm", ".potx", ".potm",
		".one", ".ecf",
		".ACCDA", ".ACCDB", ".ACCDE", ".ACCDT", ".MDA", ".MDE",".xml"}

	for _, element := range documentExtensions {
		if strings.HasSuffix(path, element) {
			if !didPrintDocumentsMsg {
				log.Info("\tIgnoring documents (such as `.pdf`, `.docx`)")
				didPrintDocumentsMsg = true
			}

			return true
		}
	}

	return false
}

// check for video files
func IsVideo(path string) bool {
	// inspired by this list: https://en.wikipedia.org/wiki/Video_file_format
	videoExtensions := [18]string{
		".mp4", ".webm", ".mkv", ".flv", ".vob", ".ogv", ".drc", ".gifv", ".mng", ".avi", ".mov", ".qt", ".mts", ".wmv", ".amv",
		".svi", ".m4v", ".mpg",
	}

	for _, element := range videoExtensions {
		if strings.HasSuffix(path, element) {
			if !didPrintVideoMsg {
				log.Info("\tIgnoring videos (such as `.mp4`)")
				didPrintVideoMsg = true
			}

			return true
		}
	}

	return false
}

// check for fonts (like .woff)
func IsFont(path string) bool {
	fontExtensions := [4]string{".ttf", ".otf", ".woff", ".woff2"}

	for _, element := range fontExtensions {
		if strings.HasSuffix(path, element) {
			if !didPrintFontsMsg {
				log.Info("\tIgnoring fonts (such as `.woff`)")
				didPrintFontsMsg = true
			}

			return true
		}
	}

	return false
}