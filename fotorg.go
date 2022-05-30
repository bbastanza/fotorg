package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// TODO break project into files that make sense
// TODO add config items for naming types and add ability to add --option to replace config with
// TODO Ability for config paths to have / at the end or not. Just a little smarter

type Config struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
}

func main() {
	args := os.Args

	config, _, err := getConfig()

	if err != nil {
		fmt.Println("Error reading config file...")
		return
	}

	// TODO how to handle args more elegantly
	if contains(args, "--no-window") {
		destFolderName := buildParentFolderName(args)
		organizeFiles(destFolderName, config)
	} else {
		runGuiApplication(config)
	}
}

func writeConfig(path string, propertyName string) {
	config, configPath, err := getConfig()

	if err != nil {
		fmt.Println("Error getting config in writeConfig function")
		return
	}

	if propertyName == "source" {
		config.Source = path
	} else {
		config.Destination = path
	}

	encodedConfig, _ := json.Marshal(config)

	err = ioutil.WriteFile(configPath, encodedConfig, os.ModePerm)

	if err != nil {
		fmt.Println("Error getting config in writeConfig function")
		return
	}
}

func runGuiApplication(initialConfig Config) {
	a := app.New()
	w := a.NewWindow("Fotorg")
	w.Resize(fyne.NewSize(800, 800))

	// Create source element
	sourceLabel := widget.NewLabel(initialConfig.Source)

	sourceBtn := widget.NewButton("source", func() {
		openPathDialog(w, "source",
			func(uri string) {
				sourceLabel.SetText(uri)
			})
	})

	sourceBtn.Alignment = widget.ButtonAlign(fyne.TextAlignCenter)

	sourceContainer := fyne.NewContainer(sourceLabel, sourceBtn)

	sourceContainer.Layout = layout.NewVBoxLayout()

	// Create destination element
	destLabel := widget.NewLabel(initialConfig.Destination)

	destBtn := widget.NewButton("destinaiton", func() {
		openPathDialog(w, "destination",
			func(uri string) {
				destLabel.SetText(uri)
			})
	})

	destBtn.Alignment = widget.ButtonAlign(fyne.TextAlignCenter)

	destinationContainer := fyne.NewContainer(destLabel, destBtn)

	destinationContainer.Layout = layout.NewVBoxLayout()

	// Created split with source on left and destination on right
	folderNameInput := widget.NewEntry()
	folderNameInput.SetText(time.Now().Format("06_01_02_"))

	form := &widget.Form{
		Items: []*widget.FormItem{ // we can specify items in the constructor
			{Text: "Folder Name", Widget: folderNameInput}},
		OnSubmit: func() { // optional, handle form submission
			config, _, _ := getConfig()
			organizeFiles(folderNameInput.Text, config)
		},
	}

	content := container.NewVBox(
		sourceContainer,
		destinationContainer,
		form,
	)

	w.SetContent(content)

	w.ShowAndRun()
}

func openPathDialog(w fyne.Window, configProperty string, callback func(uri string)) {
	d := dialog.FileDialog(
		*dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil {
				fmt.Println(err.Error())
			} else if uri == nil {
				return
			} else {
				writeConfig(uri.Path(), configProperty)
				callback(uri.Path())
			}
		}, w))

	d.Resize(w.Content().Size())

	d.Show()
}

func organizeFiles(destFolderName string, config Config) {
	// Get parent folder and source folder from config
	sourcePath := config.Source
	destinationPath := config.Destination

	OS := runtime.GOOS

	separator := "/"

	if OS == "windows" {
		separator = "\""
		sourcePath = filepath.FromSlash(sourcePath)
		destinationPath = filepath.FromSlash(destinationPath)
	}

	// Get Files in watched directory
	files, err := ioutil.ReadDir(sourcePath)

	if err != nil {
		fmt.Println("")
		return
	}

	// Get fileTypes for directories and make in destination directory
	fileTypes := getExtensionsFound(files)

	err = makeNeededDirectories(destinationPath, fileTypes, destFolderName)

	if err != nil {
		fmt.Println("Error making directories...")
		return
	}

	for _, sourceFile := range files {

		// check that the file is a regular file
		mode := sourceFile.Mode()
		if mode.IsRegular() {
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

			sourceContents, err := ioutil.ReadFile(sourceFilePath)

			if err != nil {
				fmt.Println("Error reading file contents...", sourceFile.Name())
				continue
			}

			err = ioutil.WriteFile(destFilePath, sourceContents, os.ModePerm)

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

func removeDotSafely(ext string) (string, error) {
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

		ext, _ := removeDotSafely(fullExtension)

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
	destPath = destPath + "/" + folderName + "/"

	OS := runtime.GOOS

	if OS == "windows" {
		destPath = filepath.FromSlash(destPath)
	}

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
		// check for a file with the same name as the extension
		if contains(currentNonDirList, dirName) {
			return errors.New("File already exists with name of proposed directory " + dirName)
		} else {
			os.Mkdir(destPath+dirName, os.ModePerm)
		}
	}

	return nil
}

func getConfig() (Config, string, error) {
	OS := runtime.GOOS

	homeDir, _ := os.UserHomeDir()

	configPath := homeDir + "/.config/fotorg/config.json"

	if OS == "windows" {
		configPath = filepath.FromSlash(configPath)
	}

	config, err := ioutil.ReadFile(configPath)

	if err != nil {
		return Config{}, "", err
	}

	data := Config{}

	_ = json.Unmarshal([]byte(config), &data)

	return data, configPath, nil
}
