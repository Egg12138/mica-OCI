package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/opencontainers/runtime-spec/specs-go"
	"github.com/urfave/cli"
)

var ListCommand = cli.Command{
	Name:  "list",
	Usage: "list containers",
	Description: `The list command lists all containers.`,
	Action: func(context *cli.Context) error {
		// Get root directory
		root := getRootDir(context)

		// Read container directories
		entries, err := os.ReadDir(root)
		if err != nil {
			if os.IsNotExist(err) {
				return nil // No containers exist yet
			}
			return fmt.Errorf("failed to read container directory: %w", err)
		}

		// Print container information
		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}

			// Read state file
			stateFile := filepath.Join(root, entry.Name(), "state.json")
			f, err := os.Open(stateFile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to read state file for container %s: %v\n", entry.Name(), err)
				continue
			}

			var state specs.State
			if err := json.NewDecoder(f).Decode(&state); err != nil {
				f.Close()
				fmt.Fprintf(os.Stderr, "failed to decode state file for container %s: %v\n", entry.Name(), err)
				continue
			}
			f.Close()

			// Print container information
			fmt.Printf("ID: %s\n", state.ID)
			fmt.Printf("Status: %s\n", state.Status)
			fmt.Printf("Bundle: %s\n", state.Bundle)
			fmt.Println()
		}

		return nil
	},
} 