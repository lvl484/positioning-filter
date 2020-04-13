package logger

import (
	"errors"
	"os"

	"github.com/sirupsen/logrus"
	graylog "gopkg.in/gemnasium/logrus-graylog-hook.v2"
)

const (
	srctype = "API"
	srcname = "Positioning filter"
)

// ErrBadLogDestination is the error returned when configuring failed becase of wrong destination
var ErrBadLogDestination = errors.New("logger: bad destination for logger ")

// NewLogger initialized logger according to configuration
func NewLogger(lc *Config) (*logrus.Logger, error) {
	var err error
	logger := logrus.New()
	switch lc.Output {
	case "Stdout":
		lc.setLoggerToStdout(logger)
	case "File":
		err = lc.setLoggerToFile(logger)
		if err != nil {
			return nil, err
		}
	case "Graylog":
		lc.setLoggerToGraylog(logger)
	default:
		err = ErrBadLogDestination
		return nil, err
	}
	return logger, nil
}

// setLoggerToFile initialize logger for writing to file
func (lc *Config) setLoggerToFile(logger *logrus.Logger) error {
	f, err := os.OpenFile(lc.FileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(f)
	return err
}

// setLoggerToStdout initialize logger for writing to stdout
func (lc *Config) setLoggerToStdout(logger *logrus.Logger) {
	formatter := &logrus.TextFormatter{
		FullTimestamp: true,
	}
	logger.SetFormatter(formatter)
	logger.SetOutput(os.Stdout)
}

// setLoggerToGraylog initialize logger for writing to Graylog
func (lc *Config) setLoggerToGraylog(logger *logrus.Logger) {
	hook := graylog.NewGraylogHook(lc.Host+":"+lc.Port, map[string]interface{}{srctype: srcname})
	logger.AddHook(hook)
}
