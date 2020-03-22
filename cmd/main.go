package main

import (
	"flag"
	"io"
	"log"

	"github.com/lvl484/positioning-filter/config"
	"github.com/lvl484/positioning-filter/storage"
)

var components []io.Closer

func main() {

	done := make(chan bool)

	configPath := flag.String("cp", "../config", "Path to config file")
	configName := flag.String("cn", "viper.config", "Name of config file")

	flag.Parse()

	viper, err := config.NewConfig(*configName, *configPath)
	if err != nil {
		log.Fatal(err)
	}

	consulConfig := viper.NewConsulConfig()
	consulClient, err := consulConfig.NewClient()

	if err != nil {
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

	components = append(components,
		//Put connection variables here
		db)

	gracefulShutdown(done)

	<-done

	log.Println("Service successfuly shutdown")
}
