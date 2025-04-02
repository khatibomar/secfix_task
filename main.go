package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/osquery/osquery-go"
)

type Application struct {
	OSQueryClient *osquery.ExtensionManagerClient
}

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

	parseFlags(&cfg)

	app := Application{}
	app.OSQueryClient, err = osquery.NewClient(cfg.SocketPath, 2*time.Second)
	if err != nil {
		log.Fatalf("Failed to create osquery client: %v", err)
	}
	defer app.OSQueryClient.Close()

	ctx := context.Background()
	query := "SELECT computer_name FROM system_info;"
	resp, err := app.OSQueryClient.QueryContext(ctx, query)
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}

	fmt.Println("Query Results:")
	for _, row := range resp.Response {
		fmt.Println(row)
	}
}
