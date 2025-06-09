package commands

import (
	"fmt"
	"os"
	"rmica/defs"
	"rmica/logger"
	"rmica/utils"

	pseudo_container "rmica/pseudo-container"

	"github.com/urfave/cli"
)

// FIXME: create network namespace which is required by moby!
func CreateAction(context *cli.Context) error {
	if err := utils.CheckArgs(context, 1, utils.ExactArgs); err != nil {
		return err
	}
	status, err := pseudo_container.StartContainer(context, defs.CT_ACT_CREATE, nil)
	logger.Debugf("status = %d", status)
	logger.Fprintf("status = %d", status)
	if err == nil {
		os.Exit(status)
	}
	return fmt.Errorf("`rmica create` failed: %w", err)

}


var CreateCommand = cli.Command{
	Name:  "create",
	Usage: "create a container",
	ArgsUsage: `<container-id>

Where "<container-id>" is your name for the instance of the container that you
are starting. The name you provide for the container instance must be unique on
your host.`,
	Description: `The create command creates an instance of a container for a bundle. The bundle
is a directory with a specification file named "` + defs.SpecConfig + `" and a root
filesystem.`,
	// bundle dir is the root of OCI bundle, including: config.json, rootfs, state.json, ...
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "bundle, b",
			Value: "",
			Usage: `path to the root of the bundle directory, defaults to the current directory`,
		},
		cli.StringFlag{
			Name:  "console-socket",
			Value: "",
			Usage: "path to an AF_UNIX socket which will receive a file descriptor referencing the master end of the console's pseudoterminal",
		},
		cli.StringFlag{
			Name:  "pid-file",
			Value: "",
			Usage: "specify the file to write the process id to",
		},
		cli.BoolFlag{
			Name:  "no-pivot",
			Usage: "do not use pivot root to jail process inside rootfs.  This should be used whenever the rootfs is on top of a ramdisk",
		},
	},

	Action: CreateAction,
} 