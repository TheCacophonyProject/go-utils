package logging

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

// LogArgs holds the common logging arguments
type LogArgs struct {
	LogLevel string `arg:"-l, --log-level" default:"info" help:"Set the logging level (debug, info, warn, error)"`
}

type Logger struct {
	*logrus.Logger
}

// NewLogger returns a new logger with the given log level (info, debug, warn, error)
func NewLogger(logLevel string) *Logger {
	log := logrus.New()
	log.Formatter = new(customFormatter)
	switch logLevel {
	case "debug":
		log.SetLevel(logrus.DebugLevel)
	case "info":
		log.SetLevel(logrus.InfoLevel)
	case "warn":
		log.SetLevel(logrus.WarnLevel)
	case "error":
		log.SetLevel(logrus.ErrorLevel)
	default:
		log.SetLevel(logrus.InfoLevel)
		log.Warnf("Unknown log level '%s', defaulting to info", logLevel)
	}
	return &Logger{log}
}

type customFormatter struct{}

func (f *customFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	return []byte(fmt.Sprintf("[%s] %s\n", strings.ToUpper(entry.Level.String()), entry.Message)), nil
}
