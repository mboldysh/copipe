package main

import (
	"os"

	"gopkg.in/yaml.v3"
)

func LoadLaunchConfig(filePath string) (LaunchConfig, error) {
	launchConfig := LaunchConfig{}

	configFile, err := os.ReadFile(filePath)

	if err != nil {
		return launchConfig, err
	}

	err = yaml.Unmarshal(configFile, &launchConfig)

	return launchConfig, err
}
