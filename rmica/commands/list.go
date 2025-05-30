package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/opencontainers/runtime-spec/specs-go"
	"github.com/urfave/cli"

	"rmica/logger"
	"rmica/utils"
)

var ListCommand = cli.Command{
	Name:  "list",
	Usage: "list containers",
	Description: `The list command lists all containers.`,
	Action: func(context *cli.Context) error {
		// Get root directory
		root := utils.GetRootDir(context)

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
				logger.Errorf("failed to read state file for container %s: %v", entry.Name(), err)
				continue
			}

			var state specs.State
			if err := json.NewDecoder(f).Decode(&state); err != nil {
				f.Close()
				logger.Errorf("failed to decode state file for container %s: %v", entry.Name(), err)
				continue
			}
			f.Close()

			// Print container information
			logger.Infof("ID: %s", state.ID)
			logger.Infof("Status: %s", state.Status)
			logger.Infof("Bundle: %s", state.Bundle)
			fmt.Println()
		}

		return nil
	},
} 