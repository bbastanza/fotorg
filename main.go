package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// TODO handle errors by showing dialog and write to log file
// TODO status bar
// TODO write unit tests/ figure out interfaces and dependency injection

// TODO what do the *Config and the &Config mean ... Pointers???
// type Config struct {
//   Enabled      bool
//   DatabasePath string
//   Port         string
// }

// func NewConfig() *Config {
//   return &Config{
//     Enabled:      true,
//     DatabasePath: "./example.db",
//     Port:         "8000",
//   }
// }

func main() {
	args := os.Args

	config, _, err := getConfig()

	if err != nil {
		fmt.Println("Error reading config file...")
		return
	}

	if contains(args, "--no-window") {
		destFolderName := buildParentFolderName(args)
		doTheThing(destFolderName, config)
	} else {
		runGuiApplication(config)
	}
}

func doTheThing(destFolderName string, config Config) {

	sourcePath, destinationPath := getSourceAndDestPathFromConfig(config)

	separator := getSeparator()

	files := readFiles(sourcePath)

	makeDirectories(destinationPath, files, destFolderName)

	for _, sourceFile := range files {

		if !sourceFile.Mode().IsRegular() {
			continue
		}
		// check that the file is a regular file
		dirName, _ := removeDotSafely(filepath.Ext(sourceFile.Name()))

		sourceFilePath :=
			sourcePath +
				separator +
				sourceFile.Name()

		destFilePath :=
			destinationPath +
				separator +
				destFolderName +
				separator +
				dirName +
				separator +
				sourceFile.Name()

		sourceContents := readFile(sourceFilePath)

		writeFiles(destFilePath, sourceContents)
	}
}

func buildParentFolderName(args []string) string {
	date := time.Now().Format("06_01_02_")

	if len(args) > 1 {
		return date + args[1]
	}

	var name string

	for {
		fmt.Print("Enter a folder name: ")
		fmt.Scanf("%s", &name)

		name = strings.Trim(name, " ")

		if len(name) > 0 {
			return date + name
		}
	}
}

func makeDirectories(destPath string, files []fs.FileInfo, folderName string) {
	extensionNames := getExtensionsFound(files)

	destPath = buildDestPath(destPath, folderName)

	os.Mkdir(destPath, os.ModePerm)

	// initialize empty slice to hold directory names;
	currentDirList := make([]string, 0)

	// initialize empty slice to hold non directory names;
	currentNonDirList := make([]string, 0)

	// loop through current files
	for _, file := range files {
		if file.IsDir() {
			// if the file is a directory add to the currentDirList slice
			currentDirList = append(currentDirList, file.Name())
		} else {
			// else append to the current non dir list... we may need to check this still for errors
			currentNonDirList = append(currentNonDirList, file.Name())
		}
	}

	directoriesToMake := make([]string, 0)

	for _, extension := range extensionNames {
		if !contains(currentDirList, extension) {
			directoriesToMake = append(directoriesToMake, extension)
		}
	}

	for _, dirName := range directoriesToMake {
		// check for a file with the same name as the extension
		if contains(currentNonDirList, dirName) {
			fmt.Println("File already exists with name of proposed directory " + dirName)
		} else {
			os.Mkdir(destPath+dirName, os.ModePerm)
		}
	}
}

func getSourceAndDestPathFromConfig(config Config) (string, string) {
	sourcePath := config.Source
	destinationPath := config.Destination

	OS := runtime.GOOS

	if OS == "windows" {
		sourcePath = filepath.FromSlash(sourcePath)
		destinationPath = filepath.FromSlash(destinationPath)
	}

	return sourcePath, destinationPath
}

func buildDestPath(parent string, folderName string) string {
	destPath := parent + "/" + folderName + "/"

	OS := runtime.GOOS

	if OS == "windows" {
		destPath = filepath.FromSlash(destPath)
	}

	return destPath
}
