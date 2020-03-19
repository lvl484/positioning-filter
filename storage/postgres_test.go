package storage

import (
	"strings"
	"testing"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

const (
	postgresHost = "postgres.HOST"
	postgresPort = "postgres.PORT"
	postgresUser = "postgres.USER"
	postgresPass = "postgres.PASS"
	postgresDB   = "postgres.DB"
)

func TestConnect(t *testing.T) {
	v := viper.New()
	v.AddConfigPath("../config/")
	v.SetConfigName("viper.config")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		t.Error(err)
	}

	conf := &DBConfig{
		Host: v.GetString(postgresHost),
		Port: v.GetString(postgresPort),
		User: v.GetString(postgresUser),
		Pass: v.GetString(postgresPass),
		DB:   v.GetString(postgresDB),
	}

	incorrectConf := &DBConfig{
		Host: "localhouston",
		Port: "we",
		User: "have",
		Pass: "a",
		DB:   "problem",
	}

	tests := []struct {
		name    string
		config  *DBConfig
		wantErr bool
	}{
		{
			name:    "TestWithCorrectInput",
			config:  conf,
			wantErr: false,
		},
		{
			name:    "TestWithIncorrectInput",
			config:  incorrectConf,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			_, err := Connect(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("Connect() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
