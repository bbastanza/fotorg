package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
)

type Config struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
}

func GetConfig() (Config, string, error) {
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

func WriteConfig(path string, propertyName string) {
	config, configPath, err := GetConfig()

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

func GetPaths(config Config) (string, string) {
	sourcePath := config.Source
	destinationPath := config.Destination

	OS := runtime.GOOS

	if OS == "windows" {
		sourcePath = filepath.FromSlash(sourcePath)
		destinationPath = filepath.FromSlash(destinationPath)
	}

	return sourcePath, destinationPath
}
