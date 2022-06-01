package helpers

import (
	"errors"
	"io/fs"
	"path/filepath"
	"runtime"
)

func Contains(arr []string, value string) bool {
	for _, v := range arr {
		if v == value {
			return true
		}
	}
	return false
}

func RemoveDotSafely(ext string) (string, error) {
	if len(ext) < 2 {
		return ext, errors.New("Extension too short " + ext)
	}

	return ext[1:], nil
}

func GetSeparator() string {
	OS := runtime.GOOS

	if OS == "windows" {
		return "\\"
	}

	return "/"
}

func GetExtensionsFound(files []fs.FileInfo) []string {
	fileExtensionsFound := make([]string, 0)

	for _, item := range files {
		fullExtension := filepath.Ext(item.Name())

		if fullExtension == "" || len(fullExtension) < 2 {
			continue
		}

		ext, _ := RemoveDotSafely(fullExtension)

		if !Contains(fileExtensionsFound, ext) {
			fileExtensionsFound = append(fileExtensionsFound, ext)
		}
	}

	return fileExtensionsFound
}
