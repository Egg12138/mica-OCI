package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"rmica/constants"
	pseudo_container "rmica/pseudo-container"
	"rmica/utils"

	"github.com/opencontainers/runtime-spec/specs-go"
	"github.com/urfave/cli"
)

var CreateCommand = cli.Command{
	Name:  "create",
	Usage: "create a container",
	ArgsUsage: `<container-id>

Where "<container-id>" is your name for the instance of the container that you
are starting. The name you provide for the container instance must be unique on
your host.`,
	Description: `The create command creates an instance of a container for a bundle. The bundle
is a directory with a specification file named "` + constants.SpecConfig + `" and a root
filesystem.`,
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
	},
	Action: func(context *cli.Context) error {
		if err := utils.CheckArgs(context, 1, utils.ExactArgs); err != nil {
			return err
		}

		id := context.Args().First()

		// Load the spec
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

		stateFile := filepath.Join(containerDir, "state.json")
		if err := utils.WriteJSON(stateFile, state); err != nil {
			return fmt.Errorf("failed to write state file: %w", err)
		}

		if _, err := pseudo_container.Load(root, id); err != nil {
			return fmt.Errorf("failed to verify container instance: %w", err)
		}

		if pidFile := context.String("pid-file"); pidFile != "" {
			if err := utils.CreatePidFile(pidFile, os.Getpid()); err != nil {
				return fmt.Errorf("failed to create pid file: %w", err)
			}
		}

		return nil
	},
} 