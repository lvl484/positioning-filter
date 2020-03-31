package logger

import (
	"errors"
	"os"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	graylog "gopkg.in/gemnasium/logrus-graylog-hook.v2"
)

// ErrFailedToConfigureLog is the error returned when configuring failed for some reasons
var ErrFailedToConfigureLog = errors.New("Failed to init log: failed to configure ")

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
		err = ErrFailedToConfigureLog
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
		FullTimestamp: false,
	}
	log.SetFormatter(formatter)
	log.SetOutput(os.Stdout)
}

// setLoggerToGraylog initialize logger for writing to Graylog
func (lc *Config) setLoggerToGraylog() {
	hook := graylog.NewGraylogHook(lc.Host+":"+lc.Port, map[string]interface{}{"API": "Positioning filter"})
	log.AddHook(hook)
}
