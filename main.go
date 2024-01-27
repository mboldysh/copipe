package main

import (
	"log"
	"os"

	"github.com/docker/docker/client"
)

const filePath = ".copipe.yaml"

func main() {
	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

	if err != nil {
		log.Fatalf("Not able to connect to Docker %v", err)
	}

	defer dockerClient.Close()

	launchConfig, err := LoadLaunchConfig(filePath)

	if err != nil {
		log.Fatalf("Can not unmarshal .copipe.yaml %v", err)
	}

	wdir, err := os.Getwd()

	if err != nil {
		log.Fatalf("Can not get current directory %v", err)
	}

	for _, step := range launchConfig.Steps {
		err = ExecuteStep(dockerClient, step, wdir)

		if err != nil {
			log.Fatalf("%v", err)
		}
	}

}
