package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

var (
	commands = map[string]func(cli *DockerClient, params []string) (string, error){
		"run": Run,
	}
)

// DockerClient is used to wrap around official docker Client strunct
type DockerClient struct {
	*client.Client
}

// DockerCommandRunner interface to represent our docker cmd runnner
type DockerCommandRunner interface {
	RunCommand(cmd string) (string, error)
}

// NewDockerClient creates an instace of docker Client and returns our interface
func NewDockerClient() (DockerCommandRunner, error) {
	cli, err := client.NewEnvClient()

	if err != nil {
		return nil, err
	}

	return &DockerClient{
		cli,
	}, nil
}

// RunCommand is responsable to find the proper command and run it
func (cli *DockerClient) RunCommand(cmd string) (string, error) {
	log.WithField("cmd", cmd).Info("command received")
	args := strings.Split(cmd, " ")

	if len(args) < 4 {
		return "", fmt.Errorf("command should have a least 3 arguments, cmd: %s", cmd)
	}

	cmdFn, ok := commands[args[2]]
	if !ok {
		return "", fmt.Errorf("command [docker %s] is not implemented yet", args[2])
	}

	return cmdFn(cli, args[3:])
}

// Run starts an image in background and returns the containerID
func Run(cli *DockerClient, params []string) (string, error) {
	log.WithField("params", params).Info("docker run")

	if len(params) < 1 {
		return "", errors.New("invalid params for docker run")
	}

	ctx := context.Background()

	imageName := params[0]

	out, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		return "", err
	}

	l, err := ioutil.ReadAll(out)
	if err != nil {
		return "", err
	}
	log.Info(string(l))

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: imageName,
		Cmd:   params[1:],
	}, nil, nil, "")

	if err != nil {
		return "", err
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return "", err
	}

	return resp.ID, nil
}
