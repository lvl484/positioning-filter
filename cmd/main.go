package main

import (
	"flag"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lvl484/positioning-filter/logger"
	"github.com/lvl484/positioning-filter/repository"
	"github.com/lvl484/positioning-filter/web"

	"github.com/lvl484/positioning-filter/config"
	"github.com/lvl484/positioning-filter/storage"
)

const (
	shutdownTimeout = 10 * time.Second
)

var components []io.Closer

func main() {

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
		return
	}

	filters := repository.NewFiltersRepository(db)
	srv := web.NewWebServer(filters, *serviceAddr)
	go srv.Run()

	components = append(components,
		srv,
		db)

	sigs := make(chan os.Signal)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigs
	logger.Info("Recieved", sig, "signal")

	if err := gracefulShutdown(shutdownTimeout, components); err != nil {
		logger.Error(err)
	}

	logger.Info("Service successfuly shutdown")
}
