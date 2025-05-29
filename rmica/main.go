package main

import (
	_ "embed"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/Egg12138/mica-OCI/rmica/commands"
	"github.com/Egg12138/mica-OCI/rmica/constants"
	"github.com/Egg12138/mica-OCI/rmica/logger"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

// version is set from the contents of VERSION file.
//
//go:embed VERSION
var version string
var extraVersion = ""
var gitCommit = ""



func printVersion(c *cli.Context) {
	w := c.App.Writer
	fmt.Fprintln(w, "rmica version", c.App.Version)
	if gitCommit != "" {
		fmt.Fprintln(w, "commit:", gitCommit)
	}
	fmt.Fprintln(w, "go:", runtime.Version())
}

func main() {
	app := cli.NewApp()
	app.Name = constants.RuntimeName
	app.Usage = constants.Usage
	app.Version = strings.TrimSpace(version) + extraVersion

	cli.VersionPrinter = printVersion

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "root",
			Value: constants.Root,
			Usage: "root directory for storage of container state",
		},
		cli.BoolFlag{
			Name:  "debug",
			Usage: "enable debug output",
		},
		cli.StringFlag{
			Name:  "log",
			Value: "",
			Usage: "set the log file to write logs to",
		},
		cli.StringFlag{
			Name:  "log-format",
			Value: "text",
			Usage: "set the log format ('text' or 'json')",
		},
		cli.BoolFlag{
			Name:  "systemd-mica(TODO)",
			Usage: "enable systemd mica support(TODO)",
		},
		cli.StringFlag{
			Name:  "rootless(REMOVE?)",
			Value: "auto",
			Usage: "ignore resource permission errors ('true', 'false', or 'auto')",
		},
	}

	app.Commands = []cli.Command{
		// Required by OCI specifications
		commands.CreateCommand,
		commands.StartCommand,
		commands.DeleteCommand,
		commands.ListCommand,
		commands.StateCommand,
		// Common commands
		commands.RunCommand,
		commands.SpecCommand,
		// Extenstions
	}

	app.Before = func(context *cli.Context) error {
		if err := reviseRootDir(context); err != nil {
			return err
		}


		if err := configLogrus(context); err != nil {
			return err
		}

		if context.Bool("debug") {
			logrus.Debug("Debug mode enabled")
		}

		return nil
	}

	cli.ErrWriter = &FatalWriter{cli.ErrWriter}
	if err := app.Run(os.Args); err != nil {
		// fmt.Fprintf(os.Stderr, "error: %v\n", err)
		logger.Errorf("error: %v", err)
		os.Exit(1)
	}
} 
