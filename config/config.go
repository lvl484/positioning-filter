// Package config contains application configuration.
package config

import (
	"strings"

	"github.com/spf13/viper"
)

const (
	configPath = "./"
	configName = "viper.config"
)

type PostgresConfig struct {
	Host string
	Port string
	User string
	Pass string
	DB   string
}

// NewPostgresConfig returns pointer to PointerConfig with data read from viper.config.json
func NewPostgresConfig(configName, configPath string) (*PostgresConfig, error) {
	v := viper.New()
	v.AddConfigPath(configPath)
	v.SetConfigName(configName)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	return &PostgresConfig{
		Host: v.GetString("postgres.HOST"),
		Port: v.GetString("postgres.PORT"),
		User: v.GetString("postgres.USER"),
		Pass: v.GetString("postgres.PASS"),
		DB:   v.GetString("postgres.DB"),
	}, nil
}
