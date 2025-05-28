package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/Egg12138/mica-OCI/rmica/constants"
	"github.com/opencontainers/runtime-spec/specs-go"
	"github.com/urfave/cli"
)

const (
	ExactArgs = iota
	MinArgs
	MaxArgs
)

// CheckArgs checks if the number of arguments matches the expected count
func CheckArgs(context *cli.Context, expected int, typ int) error {
	argc := context.NArg()
	if typ == ExactArgs && argc != expected {
		return fmt.Errorf("incorrect number of arguments, got %d, expected %d", argc, expected)
	}
	if typ == MinArgs && argc < expected {
		return fmt.Errorf("incorrect number of arguments, got %d, expected at least %d", argc, expected)
	}
	if typ == MaxArgs && argc > expected {
		return fmt.Errorf("incorrect number of arguments, got %d, expected at most %d", argc, expected)
	}
	return nil
}
// GetRootDir returns the root directory for containers
func GetRootDir(context *cli.Context) string {
	root := context.GlobalString("root")
	if root == "" {
		root = "/run/rmica"
	}
	return root
}

// WriteJSON writes the given data to a JSON file
func WriteJSON(path string, v interface{}) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewEncoder(f).Encode(v)
}

// CreatePidFile creates a PID file for the given process
func CreatePidFile(path string, pid int) error {
	var (
		tmpDir  = filepath.Dir(path)
		tmpName = filepath.Join(tmpDir, "."+filepath.Base(path))
	)
	f, err := os.OpenFile(tmpName, os.O_RDWR|os.O_CREATE|os.O_EXCL|os.O_SYNC, 0o666)
	if err != nil {
		return err
	}
	_, err = f.WriteString(strconv.Itoa(pid))
	f.Close()
	if err != nil {
		return err
	}
	return os.Rename(tmpName, path)
} 

// setupSpec performs initial setup based on the cli.Context for the container
func SetupSpec(context *cli.Context) (*specs.Spec, error) {
	bundle := context.String("bundle")
	if bundle != "" {
		if err := os.Chdir(bundle); err != nil {
			return nil, err
		}
	}
	spec, err := loadSpec(constants.SpecConfig)
	if err != nil {
		return nil, err
	}
	return spec, nil
}

// loadSpec loads the specification from the provided path.
func loadSpec(cPath string) (spec *specs.Spec, err error) {
	cf, err := os.Open(cPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("JSON specification file %s not found", cPath)
		}
		return nil, err
	}
	defer cf.Close()

	if err = json.NewDecoder(cf).Decode(&spec); err != nil {
		return nil, err
	}
	if spec == nil {
		return nil, errors.New("config cannot be null")
	}
	// return spec, validateProcessSpec(spec.Process)
	return spec, validateTaskSpec(spec.Process)
}

// TODO: migrate to utils_mcs
func validateTaskSpec(spec *specs.Process) error {
	return nil
}