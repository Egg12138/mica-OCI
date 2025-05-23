package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/urfave/cli"
)

var createCommand = cli.Command{
	Name:  "create",
	Usage: "create a container",
	ArgsUsage: `<container-id>

Where "<container-id>" is your name for the instance of the container that you
are starting. The name you provide for the container instance must be unique on
your host.`,
	Description: `The create command creates an instance of a container for a bundle. The bundle
is a directory with a specification file named "` + specConfig + `" and a root
filesystem.

The specification file includes an args parameter. The args parameter is used
to specify command(s) that get run when the container is started. To change the
command(s) that get executed on start, edit the args parameter of the spec. See
"runc spec --help" for more explanation.`,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "bundle, b",
			Value: ".",
			Usage: "path to the root of the bundle directory",
		},
		
		cli.StringFlag{
			Name:  "console-socket",
			Value: "",
			Usage: "path to an AF_UNIX socket which will receive a file descriptor referencing the master end of the console's pseudoterminal",
		},
		cli.StringFlag{
			Name:  "pidfd-socket",
			Usage: "path to an AF_UNIX socket which will receive a file descriptor referencing the init process",
		},
		cli.StringFlag{
			Name:  "pid-file",
			Usage: "specify the file to write the process id to",
		},
		cli.BoolFlag{
			Name:  "no-pivot",
			Usage: "do not use pivot root to jail process inside rootfs",
		},
		cli.BoolFlag{
			Name:  "no-new-keyring",
			Usage: "do not create a new session keyring for the container",
		},
		cli.StringFlag{
			Name:  "preserve-fds",
			Usage: "Pass N additional file descriptors to the container (stdio + $LISTEN_FDS + N in total)",
		},
	},
	Action: func(context *cli.Context) error {
		bundle := context.String("bundle")
		if bundle == "" {
			bundle = "."
		}
		containerID := context.Args().First()
		if containerID == "" {
			return fmt.Errorf("container ID cannot be empty")
		}

		// get pid file path
		// TODO: make log easier to locate error scene
		if pidFile := context.String("pid-file"); pidFile != "" {
			if err := os.MkdirAll(filepath.Dir(pidFile), 0755); err != nil {
				return fmt.Errorf("failed to create pid file directory: %v", err)
			}
			if err := os.WriteFile(pidFile, []byte(fmt.Sprintf("%d\n", os.Getpid())), 0644); err != nil {
				return fmt.Errorf("failed to write pid file: %v", err)
			}
		}

		if consoleSocket := context.String("console-socket"); consoleSocket != "" {
			fmt.Printf("Console socket: %s\n", consoleSocket)
		}

		fmt.Printf("Hello create\n")
		fmt.Printf("Container ID: %s\n", containerID)
		fmt.Printf("Bundle: %s\n", bundle)
		fmt.Printf("Args: %v\n", context.Args())
		fmt.Printf("Flags: %v\n", context.FlagNames())
		
		configPath := filepath.Join(bundle, "config.json")
		if _, err := os.Stat(configPath); err != nil {
			fmt.Printf("Warning: config.json not found at %s\n", configPath)
		}

		return nil
	},
}

var startCommand = cli.Command{
	Name:  "start",
	Usage: "start a container",
	ArgsUsage: `<container-id>

Where "<container-id>" is your name for the instance of the container that you
are starting.`,
	Description: `The start command executes the user defined process in a created container.`,
	Action: func(context *cli.Context) error {
		containerID := context.Args().First()
		if containerID == "" {
			return fmt.Errorf("container ID cannot be empty")
		}
		fmt.Printf("Hello start, container-id: %s\n", containerID)
		return nil
	},
}

var runCommand = cli.Command{
	Name:  "run",
	Usage: "create and run a container",
	ArgsUsage: `<container-id>

Where "<container-id>" is your name for the instance of the container that you
are starting. The name you provide for the container instance must be unique on
your host.`,
	Description: `The run command creates an instance of a container for a bundle and starts the
container. The bundle is a directory with a specification file named "` + specConfig + `"
and a root filesystem.

The specification file includes an args parameter. The args parameter is used
to specify command(s) that get run when the container is started. To change the
command(s) that get executed on start, edit the args parameter of the spec. See
"runc spec --help" for more explanation.`,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "bundle, b",
			Value: ".",
			Usage: "path to the root of the bundle directory",
		},
		cli.BoolFlag{
			Name:  "no-pivot",
			Usage: "do not use pivot root to jail process inside rootfs",
		},
		cli.BoolFlag{
			Name:  "no-new-keyring",
			Usage: "do not create a new session keyring for the container",
		},
		cli.StringFlag{
			Name:  "preserve-fds",
			Usage: "Pass N additional file descriptors to the container (stdio + $LISTEN_FDS + N in total)",
		},
	},
	Action: func(context *cli.Context) error {
		containerID := context.Args().First()
		if containerID == "" {
			return fmt.Errorf("container ID cannot be empty")
		}
		fmt.Printf("Hello run, container-id: %s\n", containerID)
		return nil
	},
}

var killCommand = cli.Command{
	Name:  "kill",
	Usage: "kill a container",
	ArgsUsage: `<container-id> [signal]

Where "<container-id>" is the name for the instance of the container and
[signal] is the signal to be sent to the init process.`,
	Description: `The kill command sends the specified signal (default: SIGTERM) to the init
process of the container.`,
	Action: func(context *cli.Context) error {
		containerID := context.Args().First()
		if containerID == "" {
			return fmt.Errorf("container ID cannot be empty")
		}
		signal := context.Args().Get(1)
		if signal == "" {
			signal = "SIGTERM"
		}
		fmt.Printf("Hello kill, container-id: %s, signal: %s\n", containerID, signal)
		return nil
	},
}

var deleteCommand = cli.Command{
	Name:  "delete",
	Usage: "delete a container",
	ArgsUsage: `<container-id>

Where "<container-id>" is the name for the instance of the container.`,
	Description: `The delete command deletes the instance of a container.`,
	Action: func(context *cli.Context) error {
		containerID := context.Args().First()
		if containerID == "" {
			return fmt.Errorf("container ID cannot be empty")
		}
		fmt.Printf("Hello delete, container-id: %s\n", containerID)
		return nil
	},
}

var psCommand = cli.Command{
	Name:  "ps",
	Usage: "list processes running inside the container",
	ArgsUsage: `<container-id> [ps options]

Where "<container-id>" is the name for the instance of the container.`,
	Description: `The ps command lists the processes running inside the container.`,
	Action: func(context *cli.Context) error {
		containerID := context.Args().First()
		if containerID == "" {
			return fmt.Errorf("container ID cannot be empty")
		}
		fmt.Printf("Hello ps, container-id: %s\n", containerID)
		return nil
	},
}

var execCommand = cli.Command{
	Name:  "exec",
	Usage: "execute new process inside the container",
	ArgsUsage: `<container-id> <command> [command options]  || -p process.json <container-id>

Where "<container-id>" is the name for the instance of the container and
"<command>" is the command to be executed in the container.`,
	Description: `The exec command executes a new process inside the container.`,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "process, p",
			Usage: "path to the process.json",
		},
		cli.BoolFlag{
			Name:  "detach, d",
			Usage: "detach from the container's process",
		},
		cli.StringFlag{
			Name:  "pid-file",
			Usage: "specify the file to write the process id to",
		},
		cli.StringFlag{
			Name:  "process-label",
			Usage: "set the asm process label for the process commonly used with selinux",
		},
		cli.StringFlag{
			Name:  "apparmor",
			Usage: "set the apparmor profile for the process",
		},
		cli.BoolFlag{
			Name:  "no-new-privs",
			Usage: "set the no new privileges value for the process",
		},
		cli.StringSliceFlag{
			Name:  "cap, c",
			Usage: "add a capability to the bounding set for the process",
		},
		cli.BoolFlag{
			Name:  "no-subreaper",
			Usage: "disable the use of the subreaper used to reap reparented processes",
		},
	},
	Action: func(context *cli.Context) error {
		containerID := context.Args().First()
		if containerID == "" {
			return fmt.Errorf("container ID cannot be empty")
		}
		command := context.Args().Get(1)
		fmt.Printf("Hello exec, container-id: %s, command: %s\n", containerID, command)
		return nil
	},
}

var listCommand = cli.Command{
	Name:  "list",
	Usage: "list containers",
	Description: `The list command lists containers started by runc with the given root.`,
	Action: func(context *cli.Context) error {
		fmt.Println("Hello list")
		return nil
	},
}

var stateCommand = cli.Command{
	Name:  "state",
	Usage: "output the state of a container",
	ArgsUsage: `<container-id>

Where "<container-id>" is the name for the instance of the container.`,
	Description: `The state command outputs the state of a container.`,
	Action: func(context *cli.Context) error {
		containerID := context.Args().First()
		if containerID == "" {
			return fmt.Errorf("container ID cannot be empty")
		}
		fmt.Printf("Hello state, container-id: %s\n", containerID)
		return nil
	},
} 