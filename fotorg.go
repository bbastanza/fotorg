package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func main() {
	watchPath := getWatchPath()
	toPath := getToPath()

	_, err := makeNeededDirectories(toPath, getExtensionsFound(watchPath))

	if err != nil {
		fmt.Println(err)
		fmt.Println(err)
		return
	}

}

func getExtensionsFound(watchPath string) []string {
	files, _ := ioutil.ReadDir(watchPath)

	fileExtensionsFound := make([]string, 0)

	for _, item := range files {
		fullExtension := filepath.Ext(item.Name())

		if fullExtension == "" || len(fullExtension) < 2 {
			continue
		}

		extension := fullExtension[1:]

		if !contains(fileExtensionsFound, extension) {
			fileExtensionsFound = append(fileExtensionsFound, extension)
		}
	}

	return fileExtensionsFound
}

func contains(arr []string, value string) bool {
	for _, v := range arr {
		if v == value {
			return true
		}
	}
	return false
}

func getWatchPath() string {
	homeDir, _ := os.UserHomeDir()
	relativePath := "/test-go/"
	return homeDir + relativePath
}

func getToPath() string {
	homeDir, _ := os.UserHomeDir()
	relativePath := "/test-go-move-to/"
	return homeDir + relativePath
}

func makeNeededDirectories(dirPath string, extensionNames []string) (bool, error) {
	// get all files in the directory we are moving to
	files, _ := ioutil.ReadDir(dirPath)

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
		if contains(currentNonDirList, dirName) {
			return false, errors.New("File already exists with name of proposed directory " + dirName)
		} else {
			os.Mkdir(dirPath+dirName, os.ModePerm)
			fmt.Println("Created directory " + dirName)
			// make the directory
		}
	}
	return true, nil
}
