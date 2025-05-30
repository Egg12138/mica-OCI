package main

import (
	_ "embed"
	"fmt"
	"os"
	"time"

	"mica"
)

//go:embed qemu-zephyr-rproc.conf
var defaultConfig string

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please specify a command (create, start, stop, rm, status)")
		return
	}

	command := os.Args[1]

	switch command {
	case "create":
		var configFile string
		if len(os.Args) == 3 {
			configFile = os.Args[2]
		} else {
			fmt.Println("Using the embedded default configuration file: qemu-zephyr-rproc.conf")
			fmt.Println(defaultConfig)
			// Write the embedded config to a temporary file
			tmpFile, err := os.CreateTemp("", "*.conf")
			if err != nil {
				fmt.Printf("Failed to create temporary config file: %v\n", err)
				return
			}
			defer os.Remove(tmpFile.Name())
			defer tmpFile.Close()

			if _, err := tmpFile.WriteString(defaultConfig); err != nil {
				fmt.Printf("Failed to write temporary config file: %v\n", err)
				return
			}
			configFile = tmpFile.Name()
		}

		// msg, err := mica.ParseConfig(configFile)
		// if err != nil {
		// 	fmt.Printf("Failed to parse config: %v\n", err)
		// 	return
		// }

		// msg.Debug = true

		if err := mica.SendCreateMsg(configFile); err != nil {
			fmt.Printf("Failed to create MICA instance: %v\n", err)
			return
		}

	case "start", "stop", "rm":
		if len(os.Args) != 3 {
			fmt.Printf("Usage: example %s <client-name>\n", command)
			return
		}
		clientName := os.Args[2]

		if err := mica.SendCtrlMsg(command, clientName); err != nil {
			fmt.Printf("Failed to %s MICA instance: %v\n", command, err)
			return
		}
		fmt.Printf("Successfully sent %s command to %s\n", command, clientName)

	case "status":
		time.Sleep(1 * time.Second)
		if err := mica.QueryStatus(); err != nil {
			fmt.Printf("Failed to query status: %v\n", err)
		}

	default:
		fmt.Printf("Unknown command: %s\n", command)
		fmt.Println("Available commands: create, start, stop, rm, status")
	}
}