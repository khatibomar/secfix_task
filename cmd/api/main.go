package main

import (
	"context"
	"flag"
	"log"
	"log/slog"
	"os"
	"strings"

	"github.com/jackc/pgx/v5"

	"github.com/khatibomar/secfix_challenge/internal/database"
)

type config struct {
	port int
	env  string
	cors struct {
		trustedOrigins []string
	}
}

type application struct {
	ctx    context.Context
	config config
	logger *slog.Logger
	db     database.Querier
}

func parseFlags(cfg *config) {
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.Func("cors-trusted-origins", "Trusted CORS origins (space separated)", func(val string) error {
		cfg.cors.trustedOrigins = strings.Fields(val)
		return nil
	})
	flag.Parse()
}

func main() {
	var cfg config

	parseFlags(&cfg)

	ctx := context.Background()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

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
		ctx:    ctx,
		config: cfg,
		logger: logger,
		db:     db,
	}

	if err = app.serve(); err != nil {
		log.Fatalf("failed to start listening on server: %v", err)
	}
}
