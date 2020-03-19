// Package config contains application configuration.
package config

import (
	"strings"

	"github.com/lvl484/positioning-filter/consul"
	"github.com/lvl484/positioning-filter/logger"
	"github.com/lvl484/positioning-filter/storage"

	"github.com/spf13/viper"
)

const (
	postgresHost = "postgres.Host"
	postgresPort = "postgres.Port"
	postgresUser = "postgres.User"
	postgresPass = "postgres.Pass"
	postgresDB   = "postgres.DB"

	loggerHost = "logger.Host"
	loggerPort = "logger.Port"

	consulAddr               = "consul.Addr"
	consulServiceName        = "consul.ServiceName"
	consulServicePort        = "consul.ServicePort"
	consulServiceHealthCheck = "consul.ServiceHealthCheck"
)

type ViperCfg struct {
	v *viper.Viper
}

func NewViperCfg(configName, configPath string) (*ViperCfg, error) {
	v := viper.New()
	v.AddConfigPath(configPath)
	v.SetConfigName(configName)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	return &ViperCfg{v: v}, nil
}

// NewDBConfig returns pointer to storage.DBConfig with data read from viper.config.json
func (vcfg *ViperCfg) NewDBConfig() *storage.DBConfig {
	return &storage.DBConfig{
		Host: vcfg.v.GetString(postgresHost),
		Port: vcfg.v.GetString(postgresPort),
		User: vcfg.v.GetString(postgresUser),
		Pass: vcfg.v.GetString(postgresPass),
		DB:   vcfg.v.GetString(postgresDB),
	}
}

// NewLoggerConfig returns pointer to logger.Config with data read from viper.config.json
func (vcfg *ViperCfg) NewLoggerConfig() *logger.Config {
	return &logger.Config{
		Host: vcfg.v.GetString(loggerHost),
		Port: vcfg.v.GetString(loggerPort),
	}
}

// NewConsulConfig returns pointer to consul.Config with data read from viper.config.json
func (vcfg *ViperCfg) NewConsulConfig() *consul.Config {
	return &consul.Config{
		Address:            vcfg.v.GetString(consulAddr),
		ServiceName:        vcfg.v.GetString(consulServiceName),
		ServicePort:        vcfg.v.GetInt(consulServicePort),
		ServiceHealthCheck: vcfg.v.GetString(consulServiceHealthCheck),
	}
}
