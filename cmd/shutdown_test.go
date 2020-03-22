package main

import (
	"database/sql"
	"log"
	"testing"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/connect"
)

func Test_gracefulShutdown_Successful(t *testing.T) {
	signal := make(chan bool)

	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		log.Fatalf("New Client implementation failed:%+v", err)
	}

	consul, err := connect.NewService("my-service", client)
	if err != nil {
		log.Fatalf("New Service implementation failed:%+v", err)
	}

	components = append(components, consul)

	gracefulShutdown(signal)

	<-signal

	er := gracefulShutdown(signal)
	if er != nil {
		t.Error(err)
	}

}

func Test_gracefulShutdown_Failed(t *testing.T) {
	signal := make(chan bool)

	postgresDB, err := sql.Open("driven", "DSN")
	if err != nil {
		log.Fatalf("Database connection failed failed:%+v", err)
	}

	components = append(components, postgresDB)

	gracefulShutdown(signal)

	<-signal

	er := gracefulShutdown(signal)
	if er == nil {
		t.Error("want error got nil")
	}

}
