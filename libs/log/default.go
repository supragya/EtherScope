package log

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
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

var _ Logger = (*defaultLogger)(nil)

type defaultLogger struct {
	logger *logrus.Logger
}

// NewDefaultLogger returns a default text logger that can be used
// and that fulfills the Logger interface. The underlying logging provider is a
// logrun logger that supports typical log levels along with JSON and plain/text
// log formats.
func NewDefaultLogger(level string) (*defaultLogger, error) {
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		return nil, err
	}
	logger := &logrus.Logger{
		Out: os.Stdout,
		Formatter: &RuntimeFormatter{ChildFormatter: &logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: time.RFC822Z,
		},
			File:         true,
			BaseNameOnly: true,
			Line:         true,
		},
		Hooks:        make(logrus.LevelHooks),
		Level:        lvl,
		ExitFunc:     os.Exit,
		ReportCaller: false,
	}
	return &defaultLogger{logger}, nil
}

func (l defaultLogger) Info(msg string, keyVals ...interface{}) {
	l.logger.WithFields(getLogFields(keyVals...)).Info(msg)
}

func (l defaultLogger) Error(msg string, keyVals ...interface{}) {
	l.logger.WithFields(getLogFields(keyVals...)).Error(msg)
}
func (l defaultLogger) Debug(msg string, keyVals ...interface{}) {
	l.logger.WithFields(getLogFields(keyVals...)).Debug(msg)
}

func (l defaultLogger) With(keyVals ...interface{}) Logger {
	return &defaultLogger{
		logger: l.logger.WithFields(getLogFields(keyVals...)).Logger,
	}
}

func getLogFields(keyVals ...interface{}) map[string]interface{} {
	if len(keyVals)%2 != 0 {
		return nil
	}

	fields := make(map[string]interface{}, len(keyVals))
	for i := 0; i < len(keyVals); i += 2 {
		fields[fmt.Sprint(keyVals[i])] = keyVals[i+1]
	}

	return fields
}
