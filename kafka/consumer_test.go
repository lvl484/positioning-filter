// Package kafka provides producer and consumer to work with kafka topics
package kafka

import (
	"reflect"
	"testing"
)

func TestNewConsumer(t *testing.T) {
	type args struct {
		config *Config
	}

	tests := []struct {
		name    string
		args    args
		want    *Consumer
		wantErr bool
	}{
		{
			name: "TestIncorrectVersion",
			args: args{
				&Config{
					Host:            "localhost",
					Port:            "9092",
					Version:         "11111",
					ConsumerGroupID: "1",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "TestIncorrectHost",
			args: args{
				&Config{
					Host:            "localghost",
					Port:            "0",
					Version:         "2.4.1",
					ConsumerGroupID: "1",
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewConsumer(tt.args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewConsumer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewConsumer() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestIntegrationNewConsumer will be passed only if kafka broker is started on localhost:9092
func TestIntegrationNewConsumer(t *testing.T) {
	config := &Config{
		Host:            "localhost",
		Port:            "9092",
		Version:         "2.4.1",
		ConsumerTopic:   "testTopic",
		ConsumerGroupID: "1",
	}
	_, err := NewConsumer(config)

	if err != nil {
		t.Errorf("NewConsumer() got %v, want nil", err)
	}
}
