package logger

import (
	"os"

	runtime "github.com/banzaicloud/logrus-runtime-formatter"
	log "github.com/sirupsen/logrus"
)

// Infallible
func SetupLog() {
	formatter := runtime.Formatter{ChildFormatter: &log.TextFormatter{
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
