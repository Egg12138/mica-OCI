package commands

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/opencontainers/runc/libcontainer/specconv"
	"github.com/urfave/cli"

	"github.com/Egg12138/mica-OCI/rmica/utils"
)

// Required fields:
// {
//     "ociVersion": "0.2.0",
//     "id": "oci-container1",
//     "status": "running",
//     "pid": 4422,
//     "bundle": "/containers/redis",
//     "annotations": {
//         "myKey": "myValue"
//     }
// }

const specConfig = "config.json"

func SpecAction(context *cli.Context) error {
	if err := utils.CheckArgs(context, 0, utils.ExactArgs); err != nil {
		return err
	}
	spec := specconv.Example()

	rootless := context.Bool("rootless")
	if rootless {
		specconv.ToRootless(spec)
	}

	checkNoFile := func(name string) error {
		_, err := os.Stat(name)
		if err == nil {
			return fmt.Errorf("file %s exists. Remove it first", name)
		}
		if !os.IsNotExist(err) {
			return err
		}
		return nil
	}
	bundle := context.String("bundle")
	if bundle != "" {
		if err := os.Chdir(bundle); err != nil {
			return err
		}
	}
	if err := checkNoFile(specConfig); err != nil {
		return err
	}
	data, err := json.MarshalIndent(spec, "", "\t")
	if err != nil {
		return err
	}
	return os.WriteFile(specConfig, data, 0o666)
}

var SpecCommand = cli.Command{
	Name:      "spec",
	Usage:     "create a new specification file",
	ArgsUsage: "",
	Description: `The spec command creates the new specification file named "` + specConfig + `" for
the bundle.
`,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "bundle, b",
			Value: "",
			Usage: "path to the root of the bundle directory",
		},
	},
	Action: SpecAction,
}