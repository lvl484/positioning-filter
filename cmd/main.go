package main

import (
	"flag"
	"log"

	"github.com/lvl484/positioning-filter/config"
	"github.com/lvl484/positioning-filter/storage"
)

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
	consulClient, err := consulConfig.NewClient()

	if err != nil {
		log.Println(err)
		return
	}

	if err = consulConfig.ServiceRegister(consulClient, agentConfig); err != nil {
		log.Println(err)
		return
	}

	defer consulClient.Agent().ServiceDeregister(consulConfig.ServiceName)

	postgresConfig := viper.NewDBConfig()
	db, err := storage.Connect(postgresConfig)

	if err != nil {
		log.Println(err)
		return
	}

	defer db.Close()

	ConnectedComponents := &structForClose{

		//Put connection variables here
	}

	done := make(chan bool)

	ConnectedComponents.GracefulShutdown(done)
	if err != nil {
		log.Fatalf("Service graceful shutdown failed: %v", err)
	}

	<-done

	log.Println("Service successfuly shutdown")

}
