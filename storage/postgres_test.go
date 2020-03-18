package storage

import (
	"log"
	"strings"
	"testing"

	_ "github.com/lib/pq"
	viper "github.com/spf13/viper"
)

const (
	configPath = "../config/"
	configName = "postgres.config"
)

func TestConnect(t *testing.T) {
	v := viper.New()
	v.AutomaticEnv()
	v.AddConfigPath(configPath)
	v.SetConfigName(configName)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	log.Println(v.ReadInConfig())

	tests := []struct {
		name    string
		host    string
		port    string
		user    string
		pass    string
		db      string
		wantErr bool
	}{
		{
			name:    "TestWithCorrectInput",
			host:    v.GetString("postgres.HOST"),
			port:    v.GetString("postgres.PORT"),
			user:    v.GetString("postgres.USER"),
			pass:    v.GetString("postgres.PASS"),
			db:      v.GetString("postgres.DB"),
			wantErr: false,
		},
		{
			name:    "TestWithIncorrectInput",
			host:    "localhouston",
			port:    "we have a problem",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			_, err := Connect(tt.host, tt.port, tt.user, tt.pass, tt.db)
			if (err != nil) != tt.wantErr {
				t.Errorf("Connect() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
