package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Config struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
}

func main() {
	config, err := getConfig()

	if err != nil {
		fmt.Println(err)
		return
	}

	sourcePath := config.Source
	destinationPath := config.Destination

	// Get Files in watched directory
	files, err1 := ioutil.ReadDir(sourcePath)

	if err1 != nil {
		fmt.Println(err1)
		return
	}

	// Get fileTypes for directories and make in destination directory
	fileTypes := getExtensionsFound(files)

	err2 := makeNeededDirectories(destinationPath, fileTypes)

	if err2 != nil {
		fmt.Println(err2)
		return
	}

	for _, sourceFile := range files {

		// check that the file is a regular file
		mode := sourceFile.Mode()
		if mode.IsRegular() {
			dirName, _ := getTypeNameFromExtension(filepath.Ext(sourceFile.Name()))

			oldPath := sourcePath + "/" + sourceFile.Name()

			newPath := destinationPath + dirName + "/" + sourceFile.Name()

			sourceContents, err := ioutil.ReadFile(oldPath)

			if err != nil {
				fmt.Println(err2)
				continue
			}

			err = ioutil.WriteFile(newPath, sourceContents, os.ModePerm)

			if err != nil {
				fmt.Println(err2)
				continue
			}

		}
	}
}

func getTypeNameFromExtension(ext string) (string, error) {
	if len(ext) < 2 {
		return ext, errors.New("Extension too short " + ext)
	}

	return ext[1:], nil
}

func getExtensionsFound(files []fs.FileInfo) []string {
	fileExtensionsFound := make([]string, 0)

	for _, item := range files {
		fullExtension := filepath.Ext(item.Name())

		if fullExtension == "" || len(fullExtension) < 2 {
			continue
		}

		ext, _ := getTypeNameFromExtension(fullExtension)

		if !contains(fileExtensionsFound, ext) {
			fileExtensionsFound = append(fileExtensionsFound, ext)
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

func makeNeededDirectories(destPath string, extensionNames []string) error {
	// get all files in the directory we are moving to
	files, _ := ioutil.ReadDir(destPath)

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
			return errors.New("File already exists with name of proposed directory " + dirName)
		} else {
			os.Mkdir(destPath+dirName, os.ModePerm)
			fmt.Println("Created directory " + dirName)
			// make the directory
		}
	}

	return nil
}

func getConfig() (Config, error) {
	homeDir, _ := os.UserHomeDir()

	configPath := homeDir + "/.config/fotorg/config.json"

	config, err := ioutil.ReadFile(configPath)

	if err != nil {
		return Config{}, err
	}

	data := Config{}

	_ = json.Unmarshal([]byte(config), &data)

	return data, nil
}

// func sourcePath() string {
// 	homeDir, _ := os.UserHomeDir()
// 	relativePath := "/test-go/"
// 	return homeDir + relativePath
// }

// func destinationPath() string {
// 	homeDir, _ := os.UserHomeDir()
// 	relativePath := "/test-go-move-to/"
// 	return homeDir + relativePath
// }
