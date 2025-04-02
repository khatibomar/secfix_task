package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/pgtype"

	"github.com/khatibomar/secfix_challenge/internal/database"
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
	logger := log.Default()

	logger.Println("Commit Hash:", getCommitHash())
	logger.Println("Git Tag:", getGitTag())

	parseFlags(&cfg)

	client, err := osquery.NewClient(cfg.Socket.Path, time.Duration(cfg.Socket.OpenTimeout)*time.Second)
	if err != nil {
		log.Fatalf("Failed to create osquery client: %v", err)
	}
	defer client.Close()

	connString := os.Getenv("SECFIX_CONNECTION_STRING")
	if connString == "" {
		log.Fatalf("Connection string is empty, please set env variable SECFIX_CONNECTION_STRING")
	}

	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer conn.Close(ctx)

	db := database.New(conn)
	app := NewApplication(logger, db, client)
	app.Verbose = cfg.Verbose
	if err := realMain(ctx, app); err != nil {
		log.Fatalf("Application failed to run: %v", err)
	}
}

func realMain(ctx context.Context, app *Application) error {
	return app.Run(ctx)
}
