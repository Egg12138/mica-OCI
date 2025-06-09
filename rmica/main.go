package main

import (
	_ "embed"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/urfave/cli"

	"rmica/commands"
	"rmica/defs"
	"rmica/logger"
	"rmica/utils"
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
	if err := logger.Fprintf("---------debug rmica file enabled"); err != nil {
		logger.Errorf("error to start debug printf utils: %s", err)
		os.Exit(114)
	}
	app := cli.NewApp()
	app.Name = defs.RuntimeName
	app.Usage = defs.Usage
	app.Version = strings.TrimSpace(version) + extraVersion
	cli.VersionPrinter = printVersion

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "root",
			Value: defs.Root,
			Usage: "root directory for storage of container state",
		},
		cli.BoolFlag{
			Name:  "debug",
			Usage: "enable debug output and flesh the debug file: " + defs.DefaultLogFile,
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

	logger.CleanDebugFile()
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


	// Only executed for subcommands
	app.Before = func(context *cli.Context) error {
		if err := utils.ReviseRootDir(context); err != nil {
			return err
		}

		// Initialize logger with CLI flags
		err := logger.Init(&logger.Config{
			Level:  "info",
			Format: context.GlobalString("log-format"),
			Output: context.GlobalString("log"),
			Debug:  context.GlobalBool("debug"),
		})
		if err != nil {
			return fmt.Errorf("failed to configure logger: %v", err)
		}

		logger.CleanDebugFile()
		logger.Debug("Debug mode enabled")

		return nil
	}

	cli.ErrWriter = &FatalWriter{cli.ErrWriter}
	if err := app.Run(os.Args); err != nil {
		logger.Errorf("error: %v", err)
		os.Exit(1)
	}
} 
