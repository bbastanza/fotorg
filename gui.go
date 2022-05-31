package main

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

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
			doTheThing(folderNameInput.Text, config)
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
