package main

import (
	"context"
	"fmt"

	"github.com/osquery/osquery-go"
)

type Application struct {
	OSQueryClient *osquery.ExtensionManagerClient
}

func NewApplication(osQueryClient *osquery.ExtensionManagerClient) *Application {
	return &Application{
		OSQueryClient: osQueryClient,
	}
}

func (app *Application) Run(ctx context.Context) error {
	query := "SELECT computer_name FROM system_info;"
	resp, err := app.OSQueryClient.QueryContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to query system info: %w", err)
	}

	fmt.Println("Query Results:")
	for _, row := range resp.Response {
		fmt.Println(row)
	}

	return nil
}
