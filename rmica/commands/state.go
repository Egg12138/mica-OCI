package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/opencontainers/runtime-spec/specs-go"
	"github.com/urfave/cli"
)

var StateCommand = cli.Command{
	Name:  "state",
	Usage: "output the state of a container",
	ArgsUsage: `<container-id>

Where "<container-id>" is the name for the instance of the container to query.`,
	Description: `The state command outputs the current state of a container.`,
	Action: func(context *cli.Context) error {
		if err := checkArgs(context, 1, exactArgs); err != nil {
			return err
		}

		// Get container ID from arguments
		id := context.Args().First()

		// Get container directory
		root := getRootDir(context)
		containerDir := filepath.Join(root, id)

		// Check if container exists
		if _, err := os.Stat(containerDir); os.IsNotExist(err) {
			return fmt.Errorf("container %s does not exist", id)
		}

		// Read state file
		stateFile := filepath.Join(containerDir, "state.json")
		f, err := os.Open(stateFile)
		if err != nil {
			return fmt.Errorf("failed to read state file: %w", err)
		}
		defer f.Close()

		var state specs.State
		if err := json.NewDecoder(f).Decode(&state); err != nil {
			return fmt.Errorf("failed to decode state file: %w", err)
		}

		// Print state information
		fmt.Printf("ID: %s\n", state.ID)
		fmt.Printf("Status: %s\n", state.Status)
		fmt.Printf("Bundle: %s\n", state.Bundle)
		if state.Pid != 0 {
			fmt.Printf("Pid: %d\n", state.Pid)
		}

		return nil
	},
} 