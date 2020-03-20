package main

import (
	"flag"
	"log"
	"time"

	"github.com/lvl484/positioning-filter/config"
	"github.com/lvl484/positioning-filter/storage"
)

const ()

func main() {
	configPath := flag.String("cp", "../config", "Path to config file")
	configName := flag.String("cn", "viper.config", "Name of config file")

	flag.Parse()

	viper, err := config.NewConfig(*configName, *configPath)
	if err != nil {
		log.Fatal(err)
	}

	consulConfig := viper.NewConsulConfig()
	agentConfig := consulConfig.AgentConfig()
	clientConfig, err := consulConfig.NewClient()

	if err != nil {
		log.Fatal(err)
	}

	if err = consulConfig.ServiceRegister(clientConfig, agentConfig); err != nil {
		log.Fatal(err)
	}

	postgresConfig := viper.NewDBConfig()
	db, err := storage.Connect(postgresConfig)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	for {
		log.Println(" [INFO] App is running.")
		time.Sleep(5 * time.Second)
	}
}
