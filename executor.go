package main

import (
	"context"
	"io"
	"log"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/google/uuid"
)

const (
	sandboxFolder = "/sandbox"

	imageNameSeparator = "-"

	tailAll = "all"

	commandsSeparator = "; "
	bash              = "sh"
	bashFlags         = "-c"

	ps4  = "PS4='$ ';"
	setX = "set -x; "
)

func ExecuteStep(dockerClient *client.Client, step Step, workingDir string) error {
	ctx := context.Background()

	pullImage(dockerClient, step.Image, ctx)

	cmd := prepareCmd(step.Script)
	containerName := prepareContainerName(step.Name)

	createdContainer, err := dockerClient.ContainerCreate(ctx,
		&container.Config{
			Image:      step.Image,
			Cmd:        cmd,
			WorkingDir: sandboxFolder,
		},
		&container.HostConfig{
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeBind,
					Source: workingDir,
					Target: sandboxFolder,
				},
			},
		}, nil, nil, containerName)

	if err != nil {
		return err
	}

	log.Printf("Starting container with name %s", containerName)

	err = dockerClient.ContainerStart(ctx, createdContainer.ID, container.StartOptions{})

	if err != nil {
		return err
	}

	logs, err := dockerClient.ContainerLogs(ctx, createdContainer.ID, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
		Tail:       tailAll,
	})

	if err != nil {
		return err
	}

	stdcopy.StdCopy(os.Stdout, os.Stderr, logs)

	statusCh, errCh := dockerClient.ContainerWait(ctx, createdContainer.ID, container.WaitConditionNotRunning)

	select {
	case err := <-errCh:
		if err != nil {
			return err
		}
	case <-statusCh:
	}

	stdcopy.StdCopy(os.Stdout, os.Stderr, logs)

	err = dockerClient.ContainerRemove(ctx, createdContainer.ID, container.RemoveOptions{})

	if err != nil {
		return err
	}

	return nil
}

func pullImage(dockerClient *client.Client, image string, ctx context.Context) error {
	reader, err := dockerClient.ImagePull(ctx, image, types.ImagePullOptions{})

	if err != nil {
		return err
	}

	defer reader.Close()

	io.Copy(os.Stdout, reader)

	return nil
}

func prepareContainerName(stepName string) string {
	return stepName + imageNameSeparator + uuid.NewString()
}

func prepareCmd(commands []string) []string {
	joinedCommands := strings.Join(commands[:], commandsSeparator)
	commandsWithConfig := ps4 + setX + joinedCommands
	cmd := []string{bash, bashFlags, commandsWithConfig}

	return cmd
}
