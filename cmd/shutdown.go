package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/connect"
	"github.com/segmentio/kafka-go"
)

type closer interface {
	Close() error
}

type structForClose struct {
	components []closer
}

//GracefulShutdown implements releasing all resouces it got from system, finish all request handling and return responses when service stopping.
func GracefulShutdown(sigs <-chan os.Signal, done chan<- bool) error {

	sig := <-sigs
	fmt.Println("Recieved", sig, "signal")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	select {
	case <-time.After(10 * time.Second):
		fmt.Println("Overslept")
	case <-ctx.Done():
		fmt.Println(ctx.Err())
	}

	// Please, replace connection variables with yours components values , or add new.
	server := &http.Server{
		Addr: ":8080",
	}

	db, err := sql.Open("driven", "DSN")
	if err != nil {
		log.Fatalf("Server shutdown failed:%+v", err)
	}

	nc, err := net.Dial("tcp", "golang.org:80")
	if err != nil {
		log.Fatalf("Server shutdown failed:%+v", err)
	}
	kafka := kafka.NewConn(nc, "topic", 0)

	client, _ := api.NewClient(api.DefaultConfig())
	consul, err := connect.NewService("my-service", client)
	if err != nil {
		log.Fatalf("Server shutdown failed:%+v", err)
	}

	otherComponents := structForClose{
		components: []closer{
			//Put connection variables here
			db,
			kafka,
			consul,
		},
	}

	defer func() {

		for _, component := range otherComponents.components {
			component.Close()
		}

		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("Server shutdown failed:%+v", err)
		}

		cancel()
	}()

	close(done)

	return nil
}
