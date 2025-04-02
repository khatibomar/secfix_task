package main

import (
	"context"
	"fmt"
	"log"

	"github.com/khatibomar/secfix_challenge/internal/database"
)

const appsQuery = "SELECT name FROM apps;"
const osVersionQuery = "SELECT version FROM os_version;"
const osQueryVersionQuery = "SELECT version FROM osquery_info;"

type OSQueryClient interface {
	QueryRowContext(ctx context.Context, query string) (map[string]string, error)
	QueryRowsContext(ctx context.Context, query string) ([]map[string]string, error)
}

type Application struct {
	OSQueryClient OSQueryClient
	Database      database.Querier
	Log           *log.Logger
	Verbose       bool
}

func NewApplication(
	logger *log.Logger,
	db database.Querier,
	osQueryClient OSQueryClient,
) *Application {
	return &Application{
		OSQueryClient: osQueryClient,
		Database:      db,
		Log:           logger,
	}
}

func (app *Application) Run(ctx context.Context) error {
	app.debug("Starting application run")

	var osVersion string
	var osqueryVersion string

	app.debug("Querying OS version")
	resp, err := app.OSQueryClient.QueryRowContext(ctx, osVersionQuery)
	if err != nil {
		return fmt.Errorf("failed to get OS version: %w", err)
	}
	osVersion = resp["version"]
	app.info("OS version: %s", osVersion)

	app.debug("Querying OSQuery version")
	resp, err = app.OSQueryClient.QueryRowContext(ctx, osQueryVersionQuery)
	if err != nil {
		return fmt.Errorf("failed to get OSQuery version: %w", err)
	}
	osqueryVersion = resp["version"]
	app.info("OSQuery version: %s", osqueryVersion)

	app.debug("Querying installed apps")
	appsQueryResp, err := app.OSQueryClient.QueryRowsContext(ctx, appsQuery)
	if err != nil {
		return fmt.Errorf("failed to get installed apps: %w", err)
	}

	appNames := make([]string, len(appsQueryResp))
	for i, row := range appsQueryResp {
		appNames[i] = row["name"]
	}

	app.info("Installed apps: %v", appNames)

	app.debug("Saving results to the database...")

	versionSnapshot, err := app.Database.InsertOSSnapshot(ctx, database.InsertOSSnapshotParams{
		OsVersion:      osVersion,
		OsQueryVersion: osqueryVersion,
	})
	if err != nil {
		return fmt.Errorf("failed to take versions snap shot: %w", err)
	}
	app.info("versions snapshot taken successfully with id = %d", versionSnapshot.ID)

	appsSnapshot, err := app.Database.InsertAppsSnapshot(ctx, appNames)
	if err != nil {
		return fmt.Errorf("failed to take versions snap shot: %w", err)
	}
	app.info("apps snapshot taken successfully with id = %d", appsSnapshot.ID)

	return nil
}

func (app *Application) debug(format string, v ...any) {
	if app.Verbose {
		app.Log.Printf("APP: [DEBUG] "+format, v...)
	}
}

func (app *Application) info(format string, v ...any) {
	app.Log.Printf("APP: [INFO] "+format, v...)
}
