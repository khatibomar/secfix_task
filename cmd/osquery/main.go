package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/khatibomar/secfix_challenge/internal"
	"github.com/khatibomar/secfix_challenge/internal/database"
	"github.com/osquery/osquery-go"
)

type config struct {
	Socket struct {
		Path        string
		OpenTimeout int
	}
	Verbose bool
}

func parseFlags(cfg *config) {
	flag.StringVar(&cfg.Socket.Path, "socket-path", "/var/run/osquery.sock", "Path to osquery socket file")
	flag.IntVar(&cfg.Socket.OpenTimeout, "socket-open-timeout", 2, "Timeout when trying to open socket connection")
	flag.BoolVar(&cfg.Verbose, "verbose", false, "Enable verbose output, good for debugging")
	flag.Parse()
}

func main() {
	var (
		cfg config
		err error
	)

	ctx := context.Background()
	logger := log.Default()

	logger.Println("Commit Hash:", internal.GetCommitHash())
	logger.Println("Git Tag:", internal.GetGitTag())

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
	app := &application{
		osQueryClient: client,
		db:            db,
		log:           logger,
		verbose:       cfg.Verbose,
	}

	if err := app.Run(ctx); err != nil {
		log.Fatalf("Application failed to run: %v", err)
	}
}
