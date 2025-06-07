package main

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"rmica/logger"
	"runtime"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)


func logrusToStderr() bool {
	l, ok := logrus.StandardLogger().Out.(*os.File)
	return ok && l.Fd() == os.Stderr.Fd()
}

type FatalWriter struct {
	cliErrWriter io.Writer
}

func (f *FatalWriter) Write(p []byte) (n int, err error) {
	logger.Error(string(p))
	if !logrusToStderr() {
		return f.cliErrWriter.Write(p)
	}
	return len(p), nil
}



func configLogrus(context *cli.Context) error {
	if context.GlobalBool("debug") {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.SetReportCaller(true)
		// Shorten function and file names reported by the logger, by
		// trimming common "github.com/opencontainers/runc" prefix.
		// This is only done for text formatter.
		_, file, _, _ := runtime.Caller(0)
		prefix := filepath.Dir(file) + "/"
		logrus.SetFormatter(&logrus.TextFormatter{
			CallerPrettyfier: func(f *runtime.Frame) (string, string) {
				function := strings.TrimPrefix(f.Function, prefix) + "()"
				fileLine := strings.TrimPrefix(f.File, prefix) + ":" + strconv.Itoa(f.Line)
				return function, fileLine
			},
		})
	}

	switch f := context.GlobalString("log-format"); f {
	case "":
		// do nothing
	case "text":
		// do nothing
	case "json":
		logrus.SetFormatter(new(logrus.JSONFormatter))
	default:
		return errors.New("invalid log-format: " + f)
	}

	if file := context.GlobalString("log"); file != "" {
		f, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_APPEND|os.O_SYNC, 0o644)
		if err != nil {
			return err
		}
		logrus.SetOutput(f)
	}

	return nil
}