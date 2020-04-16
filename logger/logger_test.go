// Package logger provides implementetion of writing Log messages to Graylog, file, and stdout.
// Supports log levels and destination.
package logger

import (
	"os"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

const (
	loggerHost     = "logger.Host"
	loggerPort     = "logger.Port"
	loggerOutput   = "logger.Output"
	loggerFileName = "logger.FileName"
)

func TestNewLogger(t *testing.T) {
	v := viper.New()
	v.AddConfigPath("../config/")
	v.SetConfigName("viper.config")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		t.Error(err)
	}

	conf := &Config{
		Host:     v.GetString(loggerHost),
		Port:     v.GetString(loggerPort),
		Output:   v.GetString(loggerOutput),
		FileName: v.GetString(loggerFileName),
	}
	incorrectConf := &Config{
		Host:   "locallviv",
		Port:   "15000",
		Output: "Graynlog13",
	}
	confFile := &Config{
		Output:   "File",
		FileName: "positioning_filter_test.log",
	}
	confStdout := &Config{
		Output: "Stdout",
	}

	tests := []struct {
		name    string
		lc      *Config
		wantErr bool
	}{
		{
			name:    "CorrectConfig1",
			lc:      conf,
			wantErr: false,
		}, {
			name:    "CorrectConfig2",
			lc:      confFile,
			wantErr: false,
		}, {
			name:    "CorrectConfig3",
			lc:      confStdout,
			wantErr: false,
		}, {
			name:    "IncorrectConfig",
			lc:      incorrectConf,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewLogger(tt.lc)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewLogger() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestConfigsetLoggerToGraylog(t *testing.T) {
	v := viper.New()
	v.AddConfigPath("../config/")
	v.SetConfigName("viper.config")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		t.Error(err)
	}

	conf := &Config{
		Host:   v.GetString(loggerHost),
		Port:   v.GetString(loggerPort),
		Output: v.GetString(loggerOutput),
	}
	incorrectConf := &Config{
		Host:   "localhost123",
		Port:   "12345",
		Output: "Graylog123",
	}

	tests := []struct {
		name    string
		lc      *Config
		wantErr bool
	}{
		{
			name:    "CorrectConfig1",
			lc:      conf,
			wantErr: false,
		}, {
			name:    "IncorrectConfig",
			lc:      incorrectConf,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewLogger(tt.lc)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewLogger() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
func TestLogConfigsetLoggerToStdout(t *testing.T) {
	confStdout := &Config{
		Output: "Stdout",
	}

	logger := logrus.New()
	confStdout.setLoggerToStdout(logger)

	assert.Equal(t, os.Stdout, logger.Out)
}

func TestLogConfigsetLoggerToFile(t *testing.T) {
	confFile := &Config{
		Output:   "File",
		FileName: "positioning_filter_test.log",
	}

	logger := logrus.New()
	err := confFile.setLoggerToFile(logger)
	assert.NoError(t, err)
	assert.NotNil(t, logger.Out)
}
