package main

import (
	"flag"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lvl484/positioning-filter/config"
	"github.com/lvl484/positioning-filter/kafka"
	"github.com/lvl484/positioning-filter/storage"
)

const (
	shutdownTimeout = 10 * time.Second
)

var components []io.Closer

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

	kafkaConfig := viper.NewKafkaConfig()
	consumer, err := kafka.NewConsumer(kafkaConfig)

	if err != nil {
		log.Println(err)
		return
	}

	filters := repository.NewFiltersRepository(db) //TODO: implement repository package
	matcher := matcher.NewMatcher(filters)         //TODO: implement matcher package
	producer, err := kafka.NewProducer(kafkaConfig)

	if err != nil {
		log.Println(err)
		return
	}

	go consumer.Consume(matcher, producer)

	components = append(components,
		//Put connection variables here
		db)

	sigs := make(chan os.Signal)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigs
	log.Println("Received", sig, "signal")

	if err := gracefulShutdown(shutdownTimeout, components); err != nil {
		log.Println(err)
	}

	log.Println("Service successfully shutdown")
}
