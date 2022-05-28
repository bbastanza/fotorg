package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// TODO Make the parent folder with format of year_month_dayGivenName ie 22_5_27_TestFolderName
// ---- Get Date + Arg1... if no arg 1 ask for the name
// TODO Ability for config paths to have / at the end or not. Just a little smarter
// TODO break project into files that make sense
// TODO add config items for naming types and add ability to add --option to replace config with

type Config struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
}

func main() {
	config, err := getConfig()

	if err != nil {
		fmt.Println("Error reading config file...")
		return
	}

	// get parent folder and source folder from config
	sourcePath := config.Source
	destinationPath := config.Destination

	args := os.Args
	folderName := buildParentFolderName(args)

	// Get Files in watched directory
	files, err1 := ioutil.ReadDir(sourcePath)

	if err1 != nil {
		fmt.Println("Error reading files from source directory...")
		return
	}

	// Get fileTypes for directories and make in destination directory
	fileTypes := getExtensionsFound(files)

	err2 := makeNeededDirectories(destinationPath+"/", fileTypes, folderName)

	if err2 != nil {
		fmt.Println("Error making directories...")
		return
	}

	for _, sourceFile := range files {

		// check that the file is a regular file
		mode := sourceFile.Mode()
		if mode.IsRegular() {
			dirName, _ := getTypeNameFromExtension(filepath.Ext(sourceFile.Name()))

			oldPath := sourcePath + "/" + sourceFile.Name()

			newPath := destinationPath + "/" + folderName + "/" + dirName + "/" + sourceFile.Name()
			fmt.Println("writing " + folderName + "/" + dirName + "/" + sourceFile.Name())

			sourceContents, err := ioutil.ReadFile(oldPath)

			if err != nil {
				fmt.Println("Error reading file contents...", sourceFile.Name())
				continue
			}

			err = ioutil.WriteFile(newPath, sourceContents, os.ModePerm)

			if err != nil {
				fmt.Println("Error writing file contents...", sourceFile.Name())
				continue
			}
		}
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

func makeNeededDirectories(destPath string, extensionNames []string, folderName string) error {

	destPath = destPath + folderName + "/"

	os.Mkdir(destPath, os.ModePerm)

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
