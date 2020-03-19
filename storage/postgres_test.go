package storage

import (
	"log"
	"strings"
	"testing"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
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
		Host: v.GetString("postgres.HOST"),
		Port: v.GetString("postgres.PORT"),
		User: v.GetString("postgres.USER"),
		Pass: v.GetString("postgres.PASS"),
		DB:   v.GetString("postgres.DB"),
	}
	log.Println(conf)

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
