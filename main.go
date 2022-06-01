package main

import (
	"fmt"
	c "fotorg/project/config"
	f "fotorg/project/files"
	h "fotorg/project/helpers"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/container"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
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

	config, _, err := c.GetConfig()

	if err != nil {
		fmt.Println("Error reading config file...")
		return
	}

	if h.Contains(args, "--no-window") {
		destFolderName := buildParentFolderName(args)
		DoTheThing(destFolderName, config)
	} else {
		runGuiApplication(config)
	}
}

func DoTheThing(destFolderName string, config c.Config) {

	sourcePath, destinationPath := getSourceAndDestPathFromConfig(config)

	separator := h.GetSeparator()

	files := f.ReadFiles(sourcePath)

	makeDirectories(destinationPath, files, destFolderName)

	for _, sourceFile := range files {

		if !sourceFile.Mode().IsRegular() {
			continue
		}
		// check that the file is a regular file
		dirName, _ := h.RemoveDotSafely(filepath.Ext(sourceFile.Name()))

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

		sourceContents := f.ReadFile(sourceFilePath)

		f.WriteFiles(destFilePath, sourceContents)
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
	extensionNames := h.GetExtensionsFound(files)

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
		if !h.Contains(currentDirList, extension) {
			directoriesToMake = append(directoriesToMake, extension)
		}
	}

	for _, dirName := range directoriesToMake {
		// check for a file with the same name as the extension
		if h.Contains(currentNonDirList, dirName) {
			fmt.Println("File already exists with name of proposed directory " + dirName)
		} else {
			os.Mkdir(destPath+dirName, os.ModePerm)
		}
	}
}

func getSourceAndDestPathFromConfig(config c.Config) (string, string) {
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

func runGuiApplication(initialConfig c.Config) {
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
			config, _, _ := c.GetConfig()
			DoTheThing(folderNameInput.Text, config)
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
				c.WriteConfig(uri.Name(), configProperty)
				callback(uri.Name())
			}
		}, w))

	d.Resize(w.Content().Size())

	d.Show()
}
