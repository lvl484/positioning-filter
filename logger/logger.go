package logger

import (
	"errors"
	"os"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	graylog "gopkg.in/gemnasium/logrus-graylog-hook.v2"
)

const (
	srctype = "API"
	srcname = "Positioning filter"
)

// ErrBadLogDestination is the error returned when configuring failed becase of wrong destination
var ErrBadLogDestination = errors.New("logger: bad destination for logger ")

// NewLogger initialized logger according to configuration
func NewLogger(lc *Config) error {
	var err error
	switch lc.Output {
	case "Stdout":
		lc.setLoggerToStdout()
	case "File":
		err = lc.setLoggerToFile()
		if err != nil {
			return err
		}
	case "Graylog":
		lc.setLoggerToGraylog()
	default:
		err = ErrBadLogDestination
		return err
	}
	return nil
}

// setLoggerToFile initialize logger for writing to file
func (lc *Config) setLoggerToFile() error {
	f, err := os.OpenFile("positioning_filter.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(f)
	return err
}

// setLoggerToStdout initialize logger for writing to stdout
func (lc *Config) setLoggerToStdout() {
	formatter := &logrus.TextFormatter{
		FullTimestamp: true,
	}
	log.SetFormatter(formatter)
	log.SetOutput(os.Stdout)
}

// setLoggerToGraylog initialize logger for writing to Graylog
func (lc *Config) setLoggerToGraylog() {
	hook := graylog.NewGraylogHook(lc.Host+":"+lc.Port, map[string]interface{}{srctype: srcname})
	log.AddHook(hook)
}
