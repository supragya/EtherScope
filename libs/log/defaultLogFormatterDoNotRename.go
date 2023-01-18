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
		line := entry.Data["line"]
		entry.Data["zz"] = fmt.Sprintf("(%s) %s:%d", functionName, _file.(string), line.(int))
		return f.ChildFormatter.Format(entry)
	}

	entry.Data["zz"] = fmt.Sprintf("(%s) %s:%s", functionName, filepath.Base(file), line)

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
	if strings.LastIndex(function, "sirupsen/logrus.") != -1 ||
		strings.LastIndex(file, "defaultLogFormatterDoNotRename") != -1 {
		skip++
		goto start
	}
	return function, file, lineNumber
}

var _ Logger = (*defaultLogger)(nil)

type defaultLogger struct {
	vals   []interface{}
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
	return &defaultLogger{vals: nil, logger: logger}, nil
}

func (l defaultLogger) Info(msg string, keyVals ...interface{}) {
	l.logger.WithFields(getLogFields(append(keyVals, l.vals...)...)).Info(msg)
}

func (l defaultLogger) Error(msg string, keyVals ...interface{}) {
	l.logger.WithFields(getLogFields(append(keyVals, l.vals...)...)).Error(msg)
}

func (l defaultLogger) Warn(msg string, keyVals ...interface{}) {
	l.logger.WithFields(getLogFields(append(keyVals, l.vals...)...)).Warn(msg)
}

func (l defaultLogger) Fatal(msg string, keyVals ...interface{}) {
	l.logger.WithFields(getLogFields(append(keyVals, l.vals...)...)).Fatal(msg)
}

func (l defaultLogger) Debug(msg string, keyVals ...interface{}) {
	l.logger.WithFields(getLogFields(append(keyVals, l.vals...)...)).Debug(msg)
}

func (l defaultLogger) Errorf(msg string, keyVals ...interface{}) {
	l.logger.WithFields(getLogFields(l.vals...)).Errorf(msg, keyVals...)
}

func (l defaultLogger) Warningf(msg string, keyVals ...interface{}) {
	l.logger.WithFields(getLogFields(l.vals...)).Warnf(msg, keyVals...)
}

func (l defaultLogger) Infof(msg string, keyVals ...interface{}) {
	l.logger.WithFields(getLogFields(l.vals...)).Infof(msg, keyVals...)
}

func (l defaultLogger) Debugf(msg string, keyVals ...interface{}) {
	l.logger.WithFields(getLogFields(l.vals...)).Debugf(msg, keyVals...)
}

func (l defaultLogger) With(keyVals ...interface{}) Logger {
	return &defaultLogger{
		vals:   append(keyVals, l.vals...),
		logger: l.logger,
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
