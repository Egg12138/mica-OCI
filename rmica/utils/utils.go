package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"runtime"
	"time"

	"os"
	"path/filepath"
	"strconv"

	"rmica/defs"

	"github.com/opencontainers/runc/libcontainer/utils"
	"github.com/opencontainers/runtime-spec/specs-go"
	"github.com/urfave/cli"
)

const (
	ExactArgs = iota
	MinArgs
	MaxArgs
	debugFileName = "/var/log/rmica.log"
)



func CleanDebugFile() error {
	f, err := os.OpenFile(debugFileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	_, err = fmt.Fprintf(f, "%s\n", timestamp)
	return err
}

func getDebugInfoPrefix() string {
	var prefix = ""
	pc_parent, _, _, ok := runtime.Caller(3)
	if ok {
		fullFuncName := runtime.FuncForPC(pc_parent).Name()
		funcName := filepath.Base(fullFuncName)
		prefix += fmt.Sprintf("%s-> ", funcName)
	}
	pc, _, _, ok := runtime.Caller(2)
	if ok {
		var caller string
		fullFuncName := runtime.FuncForPC(pc).Name()
		file, line := runtime.FuncForPC(pc).FileLine(pc)
		file = filepath.Base(file)
		caller = fullFuncName
		prefix += fmt.Sprintf("[%s:%d] %s \n\t", file, line, caller)
	}
	return prefix
}

func DebugPrintf(format string, args ...interface{}) error {
	if _, err := os.Stat(debugFileName); os.IsNotExist(err) {
		if err := CleanDebugFile(); err != nil {
			return err
		}
	}

	f, err := os.OpenFile(debugFileName, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	prefix := getDebugInfoPrefix()	
	_, err = fmt.Fprintf(f, prefix + format+"\n", args...)
	return err
}

func CheckArgs(context *cli.Context, expected int, typ int) error {
	var err error
	argc := context.NArg()
	if typ == ExactArgs && argc != expected {
		err = fmt.Errorf("incorrect number of arguments, got %d, expected %d", argc, expected)
	}
	if typ == MinArgs && argc < expected {
		err = fmt.Errorf("incorrect number of arguments, got %d, expected at least %d", argc, expected)
	}
	if typ == MaxArgs && argc > expected {
		err = fmt.Errorf("incorrect number of arguments, got %d, expected at most %d", argc, expected)
	}
	if err != nil {
		cli.ShowCommandHelp(context, context.Command.Name)
		return err
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
	spec, err := LoadSpec(defs.SpecConfig)
	if err != nil {
		return nil, err
	}
	return spec, nil
}

// NOTICE: bundle内应该包含了 适配 clientRTOS运行 的 二进制 
// IDEA: 让mica监控RTOS上的task process
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
	return spec, ValidateTaskSpec(spec)
}

// TODO: Client task
func ValidateTaskSpec(spec *specs.Spec) error {
	return nil
}

// TODO: Client task
func ValidateSpec(spec *specs.Spec) error {
	return nil
}

// From runc::libcontainer
func ValidateID(id string) error {
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


// Revise the value of flag "pid-file" to the absolute path.
func RevisePidFile(context *cli.Context) error {
	pidFile := context.GlobalString("pid-file")
	if pidFile == "" {
		return nil
	}
	pidFile, err := filepath.Abs(pidFile)
	if err != nil {
		return err
	}
	return context.Set("pid-file", pidFile)
}

// Revise the value of flag "root" to the absolute path.
func ReviseRootDir(context *cli.Context) error {
	if !context.IsSet("root") {
		return nil
	}
	root, err := filepath.Abs(context.GlobalString("root"))
	if err != nil {
		return err
	}
	if root == "/" {
		// This can happen if --root argument is
		//  - "" (i.e. empty);
		//  - "." (and the CWD is /);
		//  - "../../.." (enough to get to /);
		//  - "/" (the actual /).
		return errors.New("ojption --root argument should not be set to /")
	}

	if err := os.MkdirAll(root, 0o700); err != nil {
		return fmt.Errorf("failed to create root directory: %v", err)
	}

	if err := os.Chmod(root, 0o700); err != nil {
		return fmt.Errorf("failed to set root directory permissions: %v", err)
	}

	return context.GlobalSet("root", root)
}
