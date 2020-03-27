// Package config contains application configuration.
package config

import (
	"reflect"
	"testing"

	"github.com/lvl484/positioning-filter/consul"
	"github.com/lvl484/positioning-filter/kafka"
	"github.com/lvl484/positioning-filter/logger"
	"github.com/lvl484/positioning-filter/storage"
)

func TestNewConfig(t *testing.T) {
	type args struct {
		configName string
		configPath string
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "TestWithCorrectInput",
			args:    args{configName: "viper.config", configPath: "./"},
			wantErr: false,
		}, {
			name:    "TestWithIncorrectInput",
			args:    args{configName: "BestConfigInThisUniverse", configPath: "LongLonelyPath"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewConfig(tt.args.configName, tt.args.configPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestConfigNewConsulConfig(t *testing.T) {
	v, err := NewConfig("testConfigForViper", "./testData/")
	if err != nil {
		t.Fatalf("Cant start test, err: %v", err)
	}

	want := &consul.Config{
		Address:                "HOST2",
		ServicePort:            111,
		ServiceName:            "NAME2",
		ServiceHealthCheckPath: "HEALTH2",
	}

	if got := v.NewConsulConfig(); !reflect.DeepEqual(got, want) {
		t.Errorf("Config.NewConsulConfig() = %v, want %v", got, want)
	}
}

func TestConfigNewDBConfig(t *testing.T) {
	v, err := NewConfig("testConfigForViper", "./testData/")
	if err != nil {
		t.Fatalf("Cant start test, err: %v", err)
	}

	want := &storage.DBConfig{
		Host: "HOST1",
		Port: "PORT1",
		User: "USER1",
		Pass: "PASSWORD1",
		DB:   "DB1",
	}

	if got := v.NewDBConfig(); !reflect.DeepEqual(got, want) {
		t.Errorf("Config.NewConsulConfig() = %v, want %v", got, want)
	}
}

func TestConfigNewLoggerConfig(t *testing.T) {
	v, err := NewConfig("testConfigForViper", "./testData/")
	if err != nil {
		t.Fatalf("Cant start test, err: %v", err)
	}

	want := &logger.Config{
		Host: "HOST3",
		Port: "PORT3",
	}

	if got := v.NewLoggerConfig(); !reflect.DeepEqual(got, want) {
		t.Errorf("Config.NewLoggerConfig() = %v, want %v", got, want)
	}
}

func TestConfigNewKafkaConfig(t *testing.T) {
	v, err := NewConfig("testConfigForViper", "./testData/")
	if err != nil {
		t.Fatalf("Cant start test, err: %v", err)
	}

	want := &kafka.Config{
		Host:            "HOST4",
		Port:            "PORT4",
		ConsumerTopic:   "ConsumerTopic",
		ConsumerGroupID: "ConsumerGroupID",
		ProducerTopic:   "ProducerTopic",
	}
	if got := v.NewKafkaConfig(); !reflect.DeepEqual(got, want) {
		t.Errorf("Config.NewKafkaConfig() = %v, want %v", got, want)
	}
}
