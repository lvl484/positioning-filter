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
	"github.com/lvl484/positioning-filter/logger"
	"github.com/lvl484/positioning-filter/matcher"
	"github.com/lvl484/positioning-filter/repository"
	"github.com/lvl484/positioning-filter/storage"
	"github.com/lvl484/positioning-filter/web"
)

const (
	shutdownTimeout = 10 * time.Second
)

func main() {
	var components []io.Closer

	configPath := flag.String("cp", "../config", "Path to config file")
	configName := flag.String("cn", "viper.config", "Name of config file")
	serviceAddr := flag.String("p", ":8000", "Service addr")

	flag.Parse()

	viper, err := config.NewConfig(*configName, *configPath)
	if err != nil {
		log.Fatal(err)
	}

	loggerConfig := viper.NewLoggerConfig()
	logger, err := logger.NewLogger(loggerConfig)

	if err != nil {
		log.Println(err)
		return
	}

	consulConfig := viper.NewConsulConfig()
	agentConfig := consulConfig.AgentConfig()
	consulClient, err := consulConfig.NewClient()

	if err != nil {
		logger.Error(err)
		return
	}

	if err = consulConfig.ServiceRegister(consulClient, agentConfig); err != nil {
		logger.Error(err)
		return
	}

	defer consulClient.Agent().ServiceDeregister(consulConfig.ServiceName)

	postgresConfig := viper.NewDBConfig()
	db, err := storage.Connect(postgresConfig)

	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}

	components = append(components, db)

	filters := repository.NewFiltersRepository(db)
	matcher := matcher.NewMatcher(filters)

	kafkaConfig := viper.NewKafkaConfig()

	producer, err := kafka.NewProducer(kafkaConfig)
	if err != nil {
		logger.Error(err)

		if err := gracefulShutdown(shutdownTimeout, components); err != nil {
			logger.Error(err)
		}

		os.Exit(1)
	}

	components = append(components, producer)

	consumer, err := kafka.NewConsumer(kafkaConfig, logger)
	if err != nil {
		logger.Error(err)

		if err := gracefulShutdown(shutdownTimeout, components); err != nil {
			logger.Error(err)
		}

		os.Exit(1)
	}

	components = append(components, consumer)

	go consumer.Consume(matcher, producer)

	srv := web.NewWebServer(filters, *serviceAddr, logger)

	go func() {
		if err := srv.Run(); err != nil {
			logger.Error(err)

			if err := gracefulShutdown(shutdownTimeout, components); err != nil {
				logger.Error(err)
			}

			os.Exit(1)
		}
	}()

	components = append(components, srv)

	sigs := make(chan os.Signal)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigs
	logger.Info("Received", sig, "signal")

	if err := gracefulShutdown(shutdownTimeout, components); err != nil {
		logger.Error(err)
	}

	logger.Info("Service successfully shutdown")
}
