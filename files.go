package main

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
)

func readFile(filePath string) []byte {
	contents, err := ioutil.ReadFile(filePath)

	if err != nil {
		fmt.Println("Error Reading file...")
	}

	return contents
}

func readFiles(path string) []fs.FileInfo {
	files, err := ioutil.ReadDir(path)

	if err != nil {
		fmt.Println("Error reading from source directory...")
	}

	return files
}

func writeFiles(destFilePath string, contents []byte) {
	err := ioutil.WriteFile(destFilePath, contents, os.ModePerm)

	if err != nil {
		fmt.Println("Error Writing file...")
	}
}
