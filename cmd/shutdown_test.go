package main

import (
	"database/sql"
	"log"
	"testing"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/connect"
)

func Test_structForClose_GracefulShutdown(t *testing.T) {
	signal := make(chan bool)

	client, _ := api.NewClient(api.DefaultConfig())
	consul, err := connect.NewService("my-service", client)
	if err != nil {
		log.Fatalf("Server shutdown failed:%+v", err)
	}

	cons := structForClose{
		components: []closer{
			consul,
		},
	}

	postgresDB, err := sql.Open("driven", "DSN")
	if err != nil {
		log.Fatalf("Server shutdown failed:%+v", err)
	}

	db := structForClose{
		components: []closer{
			postgresDB,
		},
	}

	type args struct {
		done chan<- bool
	}
	tests := []struct {
		name    string
		closing *structForClose
		args    args
		wantErr bool
	}{
		{
			name:    "Successful with consul",
			closing: &cons,
			args:    args{done: signal},
			wantErr: false,
		},

		{
			name:    "Failed with database",
			closing: &db,
			args:    args{done: signal},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.closing.GracefulShutdown(tt.args.done); (err != nil) != tt.wantErr {
				t.Errorf("structForClose.GracefulShutdown() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
