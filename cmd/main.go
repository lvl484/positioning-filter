package main

import (
	"log"
	"time"

	"github.com/lvl484/positioning-filter/config"
	"github.com/lvl484/positioning-filter/storage"
)

const (
	configPath = "../config"
	configName = "viper.config"
)

func main() {
	viper, err := config.NewViperCfg(configName, configPath)
	if err != nil {
		log.Fatal(err)
	}

	// CONSUL ------------------------
	consulCfg := viper.NewConsulConfig()
	agentCfg := consulCfg.AgentConfig()
	clientCfg, err := consulCfg.NewClient()

	if err != nil {
		log.Fatal(err)
	}

	if err = consulCfg.ServiceRegister(clientCfg, agentCfg); err != nil {
		log.Fatal(err)
	}

	// -------------------------------

	// POSTGRES ----------------------

	postgresCfg := viper.NewDBConfig()
	db, err := storage.Connect(postgresCfg)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	// -------------------------------

	for {
		log.Println(" [INFO] App is running.")
		time.Sleep(5 * time.Second)
	}
}
