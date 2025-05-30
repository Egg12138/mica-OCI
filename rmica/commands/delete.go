package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/urfave/cli"

	"rmica/utils"
)

var DeleteCommand = cli.Command{
	Name:  "delete",
	Usage: "delete any resources held by the container",
	ArgsUsage: `<container-id>

Where "<container-id>" is the name for the instance of the container to be deleted.`,
	Description: `The delete command deletes any resources held by the container. The container
must be stopped before it can be deleted.`,
	Action: func(context *cli.Context) error {
		if err := utils.CheckArgs(context, 1, utils.ExactArgs); err != nil {
			return err
		}

		// Get container ID from arguments
		id := context.Args().First()

		// Get container directory
		root := utils.GetRootDir(context)
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