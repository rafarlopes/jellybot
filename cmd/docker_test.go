package cmd

import (
	"errors"
	"testing"

	"github.com/docker/docker/client"
	"github.com/rafarlopes/jellybot/mocks"
)

const succeed = "\u2713"
const failed = "\u2717"

var fakeRun = func(cli client.APIClient, params []string) (string, error) {
	return "ID", nil
}

func TestRunCommand(t *testing.T) {
	tests := []struct {
		name     string
		cmd      string
		expected string
		err      error
	}{
		{"a not implemented", "docker ps -l", "", errors.New("command [docker ps] is not implemented yet")},
		{"an invalid", "docker", "", errors.New("command should have a least 3 arguments, cmd: docker")},
		{"a valid", "docker run alpine", "ID", nil},
	}

	commands = map[string]func(cli client.APIClient, params []string) (string, error){
		"run": fakeRun,
	}

	cli := &mocks.APIClient{}

	t.Log("Given the need to run a docker command based on a message")
	{
		for i, tc := range tests {
			tf := func(t *testing.T) {

				t.Logf("\tTest: %d\tWhen sending %s command [%s]", i, tc.name, tc.cmd)
				{
					received, err := RunCommand(cli, tc.cmd)
					if tc.err != nil {
						if err != nil && tc.err.Error() == err.Error() {
							t.Logf("\t%s\tShould receive an error: %s ", succeed, tc.err.Error())
						} else {
							t.Errorf("\t%s\tShould receive an error: %s but got: %s", failed, tc.err.Error(), err.Error())
						}
					} else {
						if tc.expected == received {
							t.Logf("\t%s\tShould receive an output: %s ", succeed, tc.expected)
						} else {
							t.Errorf("\t%s\tShould receive an output: %s but got: %s", failed, tc.expected, received)
						}
					}
				}
			}
			t.Run(tc.name, tf)
		}
	}
}
