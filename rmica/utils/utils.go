package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"os"
	"path/filepath"
	"strconv"

	"rmica/constants"

	"github.com/opencontainers/runc/libcontainer/utils"
	"github.com/opencontainers/runtime-spec/specs-go"
	"github.com/urfave/cli"
)

const (
	ExactArgs = iota
	MinArgs
	MaxArgs
)

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

func GetRootDir(context *cli.Context) string {
	root := context.GlobalString("root")
	if root == "" {
		root = "/run/rmica"
	}
	return root
}

// WriteJSON writes the provided struct v to w using standard json marshaling
// without a trailing newline. This is used instead of json.Encoder because
// there might be a problem in json decoder in some cases, see:
// https://github.com/docker/docker/issues/14203#issuecomment-174177790
func WriteJSON(w io.Writer, v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

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

func SetupSpec(context *cli.Context) (*specs.Spec, error) {
	bundle := context.String("bundle")
	if bundle != "" {
		if err := os.Chdir(bundle); err != nil {
			return nil, err
		}
	}
	spec, err := LoadSpec(constants.SpecConfig)
	if err != nil {
		return nil, err
	}
	return spec, nil
}

func LoadSpec(cPath string) (spec *specs.Spec, err error) {
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

// TODO: Client task
func validateTaskSpec(spec *specs.Process) error {
	return nil
}

// From runc::libcontainer
func validateID(id string) error {
	if len(id) < 1 {
		return ErrInvalidID
	}

	// Allowed characters: 0-9 A-Z a-z _ + - .
	for i := range len(id) {
		c := id[i]
		switch {
		case c >= 'a' && c <= 'z':
		case c >= 'A' && c <= 'Z':
		case c >= '0' && c <= '9':
		case c == '_':
		case c == '+':
		case c == '-':
		case c == '.':
		default:
			return ErrInvalidID
		}

	}

	if string(os.PathSeparator)+id != utils.CleanPath(string(os.PathSeparator)+id) {
		return ErrInvalidID
	}

	return nil
}


func getDefaultImagePath() string {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return filepath.Join(cwd, "checkpoint")
}

// NOTICE: consider the tty handler of mica
func SetupIO() (*tty, error) {
	return nil, nil
}

