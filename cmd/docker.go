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
	//TODO maybe we can remove this kind of mapping and use a convention instead
	commands = map[string]func(cli client.APIClient, params []string) (string, error){
		"run": Run,
	}
)

// DockerCommandRunner interface to represent our docker cmd runnner
type DockerCommandRunner interface {
	RunCommand(cmd string) (string, error)
}

// NewDockerClient creates an instace of docker Client and returns APIClient interface
func NewDockerClient() (client.APIClient, error) {
	cli, err := client.NewEnvClient()

	if err != nil {
		return nil, err
	}

	return cli, nil
}

// RunCommand is responsable to find the proper command and run it
func RunCommand(cli client.APIClient, cmd string) (string, error) {
	//TODO can we use Cobra cmd here to make our life easier?
	log.WithField("cmd", cmd).Info("command received")
	args := strings.Split(cmd, " ")

	if len(args) < 3 {
		return "", fmt.Errorf("command should have a least 3 arguments, cmd: %s", cmd)
	}

	cmdFn, ok := commands[args[1]]
	if !ok {
		return "", fmt.Errorf("command [docker %s] is not implemented yet", args[1])
	}

	return cmdFn(cli, args[2:])
}

// Run starts an image in background and returns the containerID
func Run(cli client.APIClient, params []string) (string, error) {
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
