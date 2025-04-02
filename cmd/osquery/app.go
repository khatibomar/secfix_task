package main

import (
	"context"
	"fmt"
	"log"

	"github.com/osquery/osquery-go"
)

const appsQuery = "SELECT name FROM apps;"
const osVersionQuery = "SELECT version FROM os_version;"
const osQueryVersionQuery = "SELECT version FROM osquery_info;"

type Application struct {
	OSQueryClient *osquery.ExtensionManagerClient
	Log           *log.Logger
	Verbose       bool
}

func NewApplication(osQueryClient *osquery.ExtensionManagerClient) *Application {
	return &Application{
		OSQueryClient: osQueryClient,
		Log:           log.New(log.Writer(), "APP: ", log.LstdFlags),
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
	app.info("OS version: %s\n", osVersion)

	app.debug("Querying OSQuery version")
	resp, err = app.OSQueryClient.QueryRowContext(ctx, osQueryVersionQuery)
	if err != nil {
		return fmt.Errorf("failed to get OSQuery version: %w", err)
	}
	osqueryVersion = resp["version"]
	app.info("OSQuery version: %s\n", osqueryVersion)

	app.debug("Querying installed apps")
	appsQueryResp, err := app.OSQueryClient.QueryRowsContext(ctx, appsQuery)
	if err != nil {
		return fmt.Errorf("failed to get installed apps: %w", err)
	}

	appNames := make([]string, len(appsQueryResp))
	for i, row := range appsQueryResp {
		appNames[i] = row["name"]
	}

	app.info("Installed apps: ")
	app.info("%v", appNames)

	return nil
}

func (app *Application) debug(format string, v ...any) {
	if app.Verbose {
		app.Log.Printf("[DEBUG] "+format, v...)
	}
}

func (app *Application) info(format string, v ...any) {
	app.Log.Printf("[INFO] "+format, v...)
}
