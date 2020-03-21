// Package config contains application configuration.
package config

import (
	"reflect"
	"testing"

	"github.com/lvl484/positioning-filter/consul"
	"github.com/lvl484/positioning-filter/logger"
	"github.com/lvl484/positioning-filter/storage"
	"github.com/spf13/viper"
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
				t.Errorf("NewPostgresConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestConfigNewConsulConfig(t *testing.T) {
	v, err := NewConfig("testConfigForViper", "./testData/")
	if err != nil {
		t.Errorf("Cant start test, err: %v", err)
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
		t.Errorf("Cant start test, err: %v", err)
	}

	type fields struct {
		v *viper.Viper
	}

	test := struct {
		name   string
		fields fields
		want   *storage.DBConfig
	}{
		name:   "test",
		fields: fields{v: v.v},
		want: &storage.DBConfig{
			Host: "HOST1",
			Port: "PORT1",
			User: "USER1",
			Pass: "PASSWORD1",
			DB:   "DB1",
		},
	}

	t.Run(test.name, func(t *testing.T) {
		vcfg := &Config{
			v: test.fields.v,
		}
		if got := vcfg.NewDBConfig(); !reflect.DeepEqual(got, test.want) {
			t.Errorf("ViperCfg.NewConsulConfig() = %v, want %v", got, test.want)
		}
	})
}

func TestViperCfgNewLoggerConfig(t *testing.T) {
	v, err := NewConfig("testConfigForViper", "./testData/")
	if err != nil {
		t.Errorf("Cant start test, err: %v", err)
	}

	type fields struct {
		v *viper.Viper
	}

	test := struct {
		name   string
		fields fields
		want   *logger.Config
	}{
		name:   "test",
		fields: fields{v: v.v},
		want: &logger.Config{
			Host: "HOST3",
			Port: "PORT3",
		},
	}

	t.Run(test.name, func(t *testing.T) {
		vcfg := &Config{
			v: test.fields.v,
		}
		if got := vcfg.NewLoggerConfig(); !reflect.DeepEqual(got, test.want) {
			t.Errorf("ViperCfg.NewLoggerConfig() = %v, want %v", got, test.want)
		}
	})
}
