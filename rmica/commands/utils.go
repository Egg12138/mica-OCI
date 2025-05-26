package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/opencontainers/runtime-spec/specs-go"
	"github.com/urfave/cli"
)

const (
	exactArgs = iota
	minArgs
	maxArgs
)

// checkArgs checks if the number of arguments matches the expected count
func checkArgs(context *cli.Context, expected int, typ int) error {
	argc := context.NArg()
	if typ == exactArgs && argc != expected {
		return fmt.Errorf("incorrect number of arguments, got %d, expected %d", argc, expected)
	}
	if typ == minArgs && argc < expected {
		return fmt.Errorf("incorrect number of arguments, got %d, expected at least %d", argc, expected)
	}
	if typ == maxArgs && argc > expected {
		return fmt.Errorf("incorrect number of arguments, got %d, expected at most %d", argc, expected)
	}
	return nil
}

// setupSpec loads the OCI spec from the bundle directory
func setupSpec(context *cli.Context) (*specs.Spec, error) {
	bundle := context.String("bundle")
	if bundle == "" {
		bundle = "."
	}

	specFile := filepath.Join(bundle, "config.json")
	f, err := os.Open(specFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open spec file: %w", err)
	}
	defer f.Close()

	var spec specs.Spec
	if err := json.NewDecoder(f).Decode(&spec); err != nil {
		return nil, fmt.Errorf("failed to decode spec file: %w", err)
	}

	return &spec, nil
}

// getRootDir returns the root directory for containers
func getRootDir(context *cli.Context) string {
	root := context.GlobalString("root")
	if root == "" {
		root = "/run/rmica"
	}
	return root
}

// writeJSON writes the given data to a JSON file
func writeJSON(path string, v interface{}) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewEncoder(f).Encode(v)
}

func createPidFile(path string, pid int) error {
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