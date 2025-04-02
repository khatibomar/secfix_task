package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/osquery/osquery-go"
)

type Config struct {
	SocketPath string
}

func parseFlags(cfg *Config) {
	flag.StringVar(&cfg.SocketPath, "socket-path", "/var/run/osquery.sock", "Path to osquery socket file")
	flag.Parse()

	if cfg.SocketPath == "" {
		log.Fatal("socket-path is required")
	}
}

func main() {
	var cfg Config
	var err error

	ctx := context.Background()

	parseFlags(&cfg)

	client, err := osquery.NewClient(cfg.SocketPath, 2*time.Second)
	if err != nil {
		log.Fatalf("Failed to create osquery client: %v", err)
	}
	defer client.Close()

	app := NewApplication(client)

	if err = app.Run(ctx); err != nil {
		log.Fatalf("Application failed to run: %v", err)
	}
}
