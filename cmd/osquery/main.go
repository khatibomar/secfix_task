package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/osquery/osquery-go"
)

type Config struct {
	Socket struct {
		Path        string
		OpenTimeout int
	}
	Verbose bool
}

func parseFlags(cfg *Config) {
	flag.StringVar(&cfg.Socket.Path, "socket-path", "/var/run/osquery.sock", "Path to osquery socket file")
	flag.IntVar(&cfg.Socket.OpenTimeout, "socket-open-timeout", 2, "Timeout when trying to open socket connection")
	flag.BoolVar(&cfg.Verbose, "verbose", false, "Enable verbose output, good for debugging")
	flag.Parse()
}

func main() {
	var cfg Config
	var err error

	ctx := context.Background()

	parseFlags(&cfg)

	client, err := osquery.NewClient(cfg.Socket.Path, time.Duration(cfg.Socket.OpenTimeout)*time.Second)
	if err != nil {
		log.Fatalf("Failed to create osquery client: %v", err)
	}
	defer client.Close()

	app := NewApplication(client)
	if cfg.Verbose {
		app.Verbose = true
	}
	if err := realMain(ctx, app); err != nil {
		log.Fatalf("Application failed to run: %v", err)
	}
}

func realMain(ctx context.Context, app *Application) error {
	return app.Run(ctx)
}
