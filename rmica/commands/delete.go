package commands

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/opencontainers/runtime-spec/specs-go"
	"github.com/urfave/cli"
	"golang.org/x/sys/unix"

	"rmica/logger"
	pseudo_container "rmica/pseudo-container"
	"rmica/utils"
)

// TODO: handle mica task killing properly
func killContainer(container *pseudo_container.Container) error {
	ct := utils.GetMicaTaskConfig()
	_ = container.Signal(unix.SIGKILL, *ct)

	for range 100 {
		time.Sleep(100 * time.Millisecond)
		if err := container.Signal(unix.Signal(0), *ct); err != nil {
			return container.Destroy()
		}
	}
	return errors.New("container init still running")
}

var DeleteCommand = cli.Command{
	Name:  "delete",
	Usage: "delete any resources held by the container",
	ArgsUsage: `<container-id>

Where "<container-id>" is the name for the instance of the container to be deleted.`,
	Description: `The delete command deletes any resources held by the container. The container
must be stopped before it can be deleted.`,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "force",
			Usage: "force delete the container if it is still running (uses SIGKILL)",
		},
	},
	Action: func(context *cli.Context) error {
		if err := utils.CheckArgs(context, 1, utils.ExactArgs); err != nil {
			return err
		}

		id := context.Args().First()
		force := context.Bool("force")
		cntr, err := pseudo_container.GetContainer(context)
		containerDir := cntr.StateDir()
		if err != nil {
			logger.Debugf("deleting containerDir<%s>: %s", id, containerDir)
			if errors.Is(err, utils.ErrNotExist) {
				if e := os.RemoveAll(containerDir); e != nil {
					logger.Fprintf("failed to remove container directory %s: %w",containerDir, e)
					logger.Errorf("failed to remove container directory %s",containerDir)
				}
				if force {
					return nil
				}
			}
			return err
		}
		logger.Fprintf("err when get container %s: %s", id, err)
		logger.Debugf("err when get container %s: %s", id, err)
		
		// For demo, we do not need to kill Container
		if force {
			return killContainer(cntr)
		}

		status := cntr.State().Status
		logger.Fprintf("err when get container %s: %s", id, err)
		logger.Debugf("err when get container %s: %s", id, err)
		switch status {
		case specs.StateCreated:
			return killContainer(cntr)
		case specs.StateStopped:
			return cntr.Destroy()
		default:
			return fmt.Errorf("cannot delete container %s that is not stopped: %s", id, status)
		}
		// Delete container directory

	},
} 