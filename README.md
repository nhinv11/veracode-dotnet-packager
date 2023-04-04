# veracode-dotnet-packager ⚡

Please note that this is not an official Veracode project, not supported by Veracode in any form, and comes with no warranty whatsoever. It is simply a little pet project of mine trying to make the life of Veracode's DotNet customers a bit easier. Use at your own risk.

Please feel free to extend the existing functionality, followed by a `Merge Request` ❤️.

# Release v0.0.1
You can grab the pre-built distributions here:
https://github.com/nhinv11/veracode-dotnet-packager/releases/tag/v0.0.1

# Built-in Help 🆘

Help is built-in!

- `veracode-dotnet-packager --help` - outputs the help.

# How to Use ⚙

App needs to be compiled before running packager. 
To compile run:

dotnet publish -c Debug /p:UseAppHost=false /p:SatelliteResourceLanguages="en"

```text
Usage:
    veracode-dotnet-packager [flags]

Flags:
  -source string     The path to the dotnet project you want to package (required)
  -target string     The path where you want the vc-output.zip to be stored to (default ".")

Example:
    ./veracode-dotnet-packager -source dotnet-project -target . 
```

# Bug Reports 🐞

If you find a bug, please file an Issue right here in GitHub, and I will try to resolve it in a timely manner.
