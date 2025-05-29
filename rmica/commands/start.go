package commands

import (
	"errors"
	"fmt"

	"github.com/Egg12138/mica-OCI/rmica/utils"
	"github.com/opencontainers/runtime-spec/specs-go"
	"github.com/urfave/cli"
)

var StartCommand = cli.Command{
	Name:  "start",
	Usage: "executes the user defined process in a created container",
	ArgsUsage: `<container-id>

Where "<container-id>" is your name for the instance of the container that you
are starting. The name you provide for the container instance must be unique on
your host.`,
	Description: `The start command executes the user defined process in a created container.`,
	Action: func(context *cli.Context) error {
		if err := utils.CheckArgs(context, 1, utils.ExactArgs); err != nil {
			return err
		}

		container, err := utils.GetContainer(context)
		if err != nil {
			return err
		}

		status := container.Status()
		switch status {
		case specs.StateCreated:
			return container.Start()
		case specs.StateStopped:
			return errors.New("cannot start a container that has stopped")
		case specs.StateRunning:
			return errors.New("cannot start an already running container")
		default:
			return fmt.Errorf("cannot start a container in the %s state", status)
		}
	},
}