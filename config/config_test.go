// Package config contains application configuration.
package config

import (
	"testing"
)

func TestNewDBConfig(t *testing.T) {
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
			_, err := NewDBConfig(tt.args.configName, tt.args.configPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPostgresConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
