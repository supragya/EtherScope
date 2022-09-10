package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	// runtime "github.com/banzaicloud/logrus-runtime-formatter"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

const (
	logrusStackJump          = 4
	logrusFieldlessStackJump = 6
)

// Formatter decorates log entries with function name and package name (optional) and line number (optional)
type RuntimeFormatter struct {
	ChildFormatter logrus.Formatter
	// When true, line number will be tagged to fields as well
	Line bool
	// When true, package name will be tagged to fields as well
	Package bool
	// When true, file name will be tagged to fields as well
	File bool
	// When true, only base name of the file will be tagged to fields
	BaseNameOnly bool
}

// Format the current log entry by adding the function name and line number of the caller.
func (f *RuntimeFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	function, file, line := f.getCurrentPosition(entry)

	packageEnd := strings.LastIndex(function, ".")
	functionName := function[packageEnd+1:]

	if _file, ok := entry.Data["file"]; ok {
		// ENOK or ENOKS invoked
		line, _ := entry.Data["line"]
		entry.Data = logrus.Fields{"CODE": fmt.Sprintf("(%s) %s:%d", functionName, _file.(string), line.(int))}
		return f.ChildFormatter.Format(entry)
	}

	entry.Data = logrus.Fields{"CODE": fmt.Sprintf("(%s) %s:%s", functionName, filepath.Base(file), line)}
	return f.ChildFormatter.Format(entry)
}

func (f *RuntimeFormatter) getCurrentPosition(entry *logrus.Entry) (string, string, string) {
	skip := logrusStackJump
	if len(entry.Data) == 0 {
		skip = logrusFieldlessStackJump
	}
start:
	pc, file, line, _ := runtime.Caller(skip)
	lineNumber := ""
	if f.Line {
		lineNumber = fmt.Sprintf("%d", line)
	}
	function := runtime.FuncForPC(pc).Name()
	if strings.LastIndex(function, "sirupsen/logrus.") != -1 {
		skip++
		goto start
	}
	return function, file, lineNumber
}

// Infallible
func SetupLog() {
	formatter := RuntimeFormatter{ChildFormatter: &log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "06/01/02 15:04:05.00 Z0700",
	},
		File:         true,
		BaseNameOnly: true}
	formatter.Line = true
	log.SetFormatter(&formatter)
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func SetLogLevel(lvl string) error {
	parsedLevel, err := log.ParseLevel(lvl)
	if err != nil {
		return err
	}
	log.SetLevel(parsedLevel)
	return nil
}
