// Package logger provides implementetion of writing Log messages to Graylog, file, and stdout.
// Supports log levels and destination.
package logger

import (
	"os"
	"strings"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

const (
	loggerHost   = "logger.Host"
	loggerPort   = "logger.Port"
	loggerOutput = "logger.Output"
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
		Host:   v.GetString(loggerHost),
		Port:   v.GetString(loggerPort),
		Output: v.GetString(loggerOutput),
	}
	incorrectConf := &Config{
		Host:   "locallviv",
		Port:   "15000",
		Output: "Graynlog13",
	}
	confFile := &Config{
		Output: "File",
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
			err := NewLogger(tt.lc)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewLogger() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestLogConfigsetLoggerToFile(t *testing.T) {

	_, err := os.OpenFile("positioning_filter.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)

	if err != nil {
		t.Errorf("Config.setLoggerToFile() error = %v", err)
	}
}

func TestLogConfigsetLoggerToStdout(t *testing.T) {

	conf_file := &Config{
		Output: "Filename",
	}
	conf_stdout := &Config{
		Output: "Stdout",
	}

	conf_stdout.setLoggerToStdout()
	assert.Equal(t, os.Stdout, log.StandardLogger().Out)
	conf_file.setLoggerToStdout()
	assert.Equal(t, os.Stdout, log.StandardLogger().Out)

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
			err := NewLogger(tt.lc)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewLogger() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
