package files

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
)

func ReadFile(filePath string) []byte {
	contents, err := ioutil.ReadFile(filePath)

	if err != nil {
		fmt.Println("Error Reading file...")
	}

	return contents
}

func ReadFiles(path string) []fs.FileInfo {
	files, err := ioutil.ReadDir(path)

	if err != nil {
		fmt.Println("Error reading from source directory...")
	}

	return files
}

func WriteFiles(destFilePath string, contents []byte) {
	err := ioutil.WriteFile(destFilePath, contents, os.ModePerm)

	if err != nil {
		fmt.Println("Error Writing file...")
	}
}
