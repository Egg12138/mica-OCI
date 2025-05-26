package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"

	"github.com/opencontainers/runtime-spec/specs-go"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)
 
const (
	exactArgs = iota
	minArgs
	maxArgs
)

func checkArgs(context *cli.Context, expected int, checkType int) error {
	var err error
	cmdName := context.Command.Name
	switch checkType {
	case exactArgs:
		if context.NArg() != expected {
			err = fmt.Errorf("%s: %q requires exactly %d argument(s)", os.Args[0], cmdName, expected)
		}
	case minArgs:
		if context.NArg() < expected {
			err = fmt.Errorf("%s: %q requires a minimum of %d argument(s)", os.Args[0], cmdName, expected)
		}
	case maxArgs:
		if context.NArg() > expected {
			err = fmt.Errorf("%s: %q requires a maximum of %d argument(s)", os.Args[0], cmdName, expected)
		}
	}

	if err != nil {
		fmt.Printf("Incorrect Usage.\n\n")
		_ = cli.ShowCommandHelp(context, cmdName)
		return err
	}
	return nil
}

func logrusToStderr() bool {
	l, ok := logrus.StandardLogger().Out.(*os.File)
	return ok && l.Fd() == os.Stderr.Fd()
}

// fatal prints the error's details if it is a libcontainer specific error type
// then exits the program with an exit status of 1.
func fatal(err error) {
	fatalWithCode(err, 1)
}

func fatalWithCode(err error, ret int) {
	// Make sure the error is written to the logger.
	logrus.Error(err)
	if !logrusToStderr() {
		fmt.Fprintln(os.Stderr, err)
	}

	os.Exit(ret)
}

// setupSpec performs initial setup based on the cli.Context for the container
// rmica just load the spec and take into account some items about RTOS tasks 
// and RTOS target OS
func setupSpec(context *cli.Context) (*specs.Spec, error) {
	bundle := context.String("bundle")
	if bundle != "" {
		if err := os.Chdir(bundle); err != nil {
			return nil, err
		}
	}
	spec, err := loadSpec(specConfig)
	if err != nil {
		return nil, err
	}
	return spec, nil
}

// 
func revisePidFile(context *cli.Context) error {
	pidFile := context.String("pid-file")
	if pidFile == "" {
		return nil
	}
	// convert pid-file to an absolute path so we can write to the right
	// file after chdir to bundle
	// NOTICE: for mica, we do not need chroot!
	pidFile, err := filepath.Abs(pidFile)
	if err != nil {
		return err
	}
	logrus.Debugln("pidFile is revised")
	return context.Set("pid-file", pidFile)

}

// reviseRootDir ensures that the --root option argument,
// if specified, is converted to an absolute and cleaned path,
// and that this path is sane.
func reviseRootDir(context *cli.Context) error {
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

func configLogrus(context *cli.Context) error {
	if context.GlobalBool("debug") {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}

	switch context.GlobalString("log-format") {
	case "json":
		logrus.SetFormatter(&logrus.JSONFormatter{})
	case "text":
		logrus.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	default:
		return fmt.Errorf("invalid log format: %s", context.GlobalString("log-format"))
	}

	if logPath := context.GlobalString("log"); logPath != "" {
		f, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)
		if err != nil {
			return fmt.Errorf("failed to open log file: %v", err)
		}
		logrus.SetOutput(f)
	}

	return nil
}

func getRootDir(context *cli.Context) string {
	// root := context.GlobalString("root")
	root := "/run/rmica"
	
	return root
}

// ensureRootDir 确保根目录存在并具有正确的权限
func ensureRootDir(root string) error {
	if err := os.MkdirAll(root, 0o700); err != nil {
		return fmt.Errorf("failed to create root directory: %v", err)
	}

	if err := os.Chmod(root, 0o700); err != nil {
		return fmt.Errorf("failed to set root directory permissions: %v", err)
	}

	return nil
}

// TODO: migrate to spec.go
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

type FatalWriter struct {
	cliErrWriter io.Writer
}

func (f *FatalWriter) Write(p []byte) (n int, err error) {
	logrus.Error(string(p))
	if !logrusToStderr() {
		return f.cliErrWriter.Write(p)
	}
	return len(p), nil
}

func createPidFile(path string, pid int) error {
	tmpDir := filepath.Dir(path)
	tmpName := filepath.Join(tmpDir, "."+filepath.Base(path))
	
	f, err := os.OpenFile(tmpName, os.O_RDWR|os.O_CREATE|os.O_EXCL|os.O_SYNC, 0o666)
	if err != nil {
		return err
	}
	
	_, err = f.WriteString(strconv.Itoa(pid))
	f.Close()
	if err != nil {
		os.Remove(tmpName)
		return err
	}
	
	return os.Rename(tmpName, path)
}
