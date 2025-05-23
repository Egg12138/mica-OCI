package main

import (
	_ "embed"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/urfave/cli"
)

// version is set from the contents of VERSION file.
//
//go:embed VERSION
var version string

var extraVersion = ""
var gitCommit = ""

const (
	runtimeName = "rmica"
	specConfig = "config.json"
	usage      = `Simple Pseudo-Container Runtime

A simple drop-in replacement for runc that implements basic container lifecycle management APIs
but does not actually handling any containers following the OCI specification.
`
)

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
	app.Name = "rmica"
	app.Usage = usage
	app.Version = strings.TrimSpace(version) + extraVersion

	cli.VersionPrinter = printVersion

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "root",
			Value: "/run/rmica",
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

	// TODO: lots of commands were not provided
	app.Commands = []cli.Command{
		createCommand,
		startCommand,
		runCommand,
		killCommand,
		deleteCommand,
		psCommand,
		execCommand,
		listCommand,
		stateCommand,
	}

	app.Before = func(context *cli.Context) error {
		if err := reviseRootDir(context); err != nil {
			return err
		}

		if context.Bool("debug") {
			fmt.Println("Debug mode enabled")
		}

		if err := configLogrus(context); err != nil {
			return err
		}

		return nil
	}

	cli.ErrWriter = &FatalWriter{cli.ErrWriter}
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "rmica: %v\n", err)
		os.Exit(1)
	}
} 
