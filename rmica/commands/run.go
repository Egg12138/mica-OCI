package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"rmica/defs"
	"rmica/utils"

	"github.com/opencontainers/runtime-spec/specs-go"
	"github.com/urfave/cli"
)

var RunCommand = cli.Command{
	Name:  "run",
	Usage: "create and run a container",
	ArgsUsage: `<container-id>

Where "<container-id>" is your name for the instance of the container that you
are starting. The name you provide for the container instance must be unique on
your host.`,
	Description: `The run command creates an instance of a container for a bundle and starts the
process inside the container. The bundle is a directory with a specification
file named "config.json" and a root filesystem.`,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "bundle, b",
			Value: "",
			Usage: `path to the root of the bundle directory, defaults to the current directory`,
		},
		cli.StringFlag{
			Name:  "pid-file",
			Value: "",
			Usage: "specify the file to write the process id to",
		},
		cli.BoolFlag{
			Name:  "detach, d",
			Usage: "detach from the container's process",
		},
	},
	Action: func(context *cli.Context) error {
		if err := utils.CheckArgs(context, 1, utils.ExactArgs); err != nil {
			return err
		}

		id := context.Args().First()

		spec, err := utils.SetupSpec(context)
		if err != nil {
			return fmt.Errorf("failed to load spec: %w", err)
		}

		// Create container directory
		root := utils.GetRootDir(context)
		containerDir := filepath.Join(root, id)
		if err := os.MkdirAll(containerDir, 0o700); err != nil {
			return fmt.Errorf("failed to create container directory: %w", err)
		}

		// Create state file
		state := &specs.State{
			Version:     spec.Version,
			ID:          id,
			Status:      specs.StateCreating,
			Bundle:      context.String("bundle"),
			Annotations: spec.Annotations,
		}

		statePath := filepath.Join(containerDir, defs.StateFilename)
		stateFile, err := os.OpenFile(statePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
		if err != nil {
			return err
		}
		defer stateFile.Close()
		if err := utils.WriteJSON(stateFile, state); err != nil {
			return fmt.Errorf("failed to write state file: %w", err)
		}

		if pidFile := context.String("pid-file"); pidFile != "" {
			if err := utils.CreatePidFile(pidFile, os.Getpid()); err != nil {
				return fmt.Errorf("failed to create pid file: %w", err)
			}
		}

		// TODO: Implement container process execution
		// This is where we would start the container process
		// For now, we just update the state to running
		state.Status = specs.StateRunning
		if err := utils.WriteJSON(stateFile, state); err != nil {
			return fmt.Errorf("failed to update state file: %w", err)
		}

		return nil
	},
} 