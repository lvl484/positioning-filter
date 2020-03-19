// Package config contains application configuration.
package config

import (
	"strings"

	"github.com/lvl484/positioning-filter/storage"
	"github.com/spf13/viper"
)

// NewPostgresConfig returns pointer to PointerConfig with data read from viper.config.json
func NewDBConfig(configName, configPath string) (*storage.DBConfig, error) {
	v := viper.New()
	v.AddConfigPath(configPath)
	v.SetConfigName(configName)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	return &storage.DBConfig{
		Host: v.GetString("postgres.HOST"),
		Port: v.GetString("postgres.PORT"),
		User: v.GetString("postgres.USER"),
		Pass: v.GetString("postgres.PASS"),
		DB:   v.GetString("postgres.DB"),
	}, nil
}
