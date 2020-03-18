package storage

import (
	"testing"

	_ "github.com/lib/pq"

	"github.com/lvl484/positioning-filter/config"
)

func Test_Connect(t *testing.T) {
	conf, err := config.NewPostgresConfig()
	if err != nil {
		t.Error(err)
	}

	incorrectConf := &config.PostgresConfig{
		Host: "localhouston",
		Port: "we",
		User: "have",
		Pass: "a",
		DB:   "problem",
	}

	tests := []struct {
		name    string
		config  *config.PostgresConfig
		wantErr bool
	}{
		{
			name:    "TestWithCorrectInput",
			config:  conf,
			wantErr: false,
		},
		{
			name:    "TestWithIncorrectInput",
			config:  incorrectConf,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			_, err := Connect(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("Connect() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
