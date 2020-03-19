// Package config contains application configuration.
package config

import (
	"strings"

	"github.com/lvl484/positioning-filter/storage"
	"github.com/spf13/viper"
)

const (
	postgresHost = "postgres.HOST"
	postgresPort = "postgres.PORT"
	postgresUser = "postgres.USER"
	postgresPass = "postgres.PASS"
	postgresDB   = "postgres.DB"
)

// NewDBConfig returns pointer to storage.DBConfig with data read from viper.config.json
func NewDBConfig(configName, configPath string) (*storage.DBConfig, error) {
	v := viper.New()
	v.AddConfigPath(configPath)
	v.SetConfigName(configName)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	return &storage.DBConfig{
		Host: v.GetString(postgresHost),
		Port: v.GetString(postgresPort),
		User: v.GetString(postgresUser),
		Pass: v.GetString(postgresPass),
		DB:   v.GetString(postgresDB),
	}, nil
}
