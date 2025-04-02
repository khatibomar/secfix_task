package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/osquery/osquery-go"
)

const appsQuery = "SELECT name FROM apps;"
const osVersionQuery = "SELECT version FROM os_version;"
const osQueryVersionQuery = "SELECT version FROM osquery_info;"

type Application struct {
	OSQueryClient *osquery.ExtensionManagerClient
}

func NewApplication(osQueryClient *osquery.ExtensionManagerClient) *Application {
	return &Application{
		OSQueryClient: osQueryClient,
	}
}

func (app *Application) Run(ctx context.Context) error {
	var osVersion string
	var osqueryVersion string

	resp, err := app.OSQueryClient.QueryRowContext(ctx, osVersionQuery)
	if err != nil {
		return fmt.Errorf("failed to get OS version: %w", err)
	}

	osVersion = resp["version"]

	resp, err = app.OSQueryClient.QueryRowContext(ctx, osQueryVersionQuery)
	if err != nil {
		return fmt.Errorf("failed to get OSQuery version: %w", err)
	}
	osqueryVersion = resp["version"]

	appsQueryResp, err := app.OSQueryClient.QueryRowsContext(ctx, appsQuery)
	if err != nil {
		return fmt.Errorf("failed to get installed apps: %w", err)
	}

	appNames := make([]string, len(appsQueryResp))
	fmt.Println("Query Results:")
	for i, row := range appsQueryResp {
		appNames[i] = row["name"]
	}

	var sb strings.Builder
	sb.Grow(1024)
	for _, n := range appNames {
		sb.WriteString(n)
		sb.WriteRune('\n')
	}

	fmt.Printf("OS version: %s\n", osVersion)
	fmt.Printf("OSQuery version: %s\n", osqueryVersion)
	fmt.Println("Installed apps:")
	fmt.Print(sb.String())

	return nil
}
