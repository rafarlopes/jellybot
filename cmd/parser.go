package cmd

import (
	"errors"
	"strings"

	log "github.com/sirupsen/logrus"
)

type command struct {
	name string
	exec func(params []string) (string, error)
}

var (
	errInvalidCommand = errors.New("invalid command")
	tools             = map[string]map[string]command{
		"docker": map[string]command{
			"run": command{
				name: "run",
				exec: Run,
			},
		},
	}
)

// Parse is responsable to find the proper command and run it
func Parse(cmd string) (string, error) {
	log.WithField("cmd", cmd).Info("command received")
	pieces := strings.Split(cmd, " ")

	if len(pieces) < 4 {
		return "", errInvalidCommand
	}

	app, ok := tools[pieces[1]]

	if !ok {
		return "", errInvalidCommand
	}

	c, ok := app[pieces[2]]

	if !ok {
		return "", errInvalidCommand
	}

	return c.exec(pieces[3:])
}
