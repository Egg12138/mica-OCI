package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/urfave/cli"
)

var DeleteCommand = cli.Command{
	Name:  "delete",
	Usage: "delete any resources held by the container",
	ArgsUsage: `<container-id>

Where "<container-id>" is the name for the instance of the container to be deleted.`,
	Description: `The delete command deletes any resources held by the container. The container
must be stopped before it can be deleted.`,
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

		// TODO: Check if container is stopped
		// For now, we just delete the container directory

		// Delete container directory
		if err := os.RemoveAll(containerDir); err != nil {
			return fmt.Errorf("failed to delete container directory: %w", err)
		}

		return nil
	},
} 