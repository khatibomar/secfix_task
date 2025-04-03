package main

import (
	"net/http"
)

func (app *application) latestDataHandler(w http.ResponseWriter, r *http.Request) {
	info, err := app.db.GetLatestOSInfo(app.ctx)
	if err != nil {
    		app.serverErrorResponse(w, r, err)
	}

	apps, err := app.db.GetLatestAppsSnapshot(app.ctx)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	env := envelope{
		"os_versions": map[string]any{
			"os_query_version": info.OsQueryVersion,
			"os_version":       info.OsVersion,
			"snapshot_at":      apps.SnapshotTime,
		},
		"apps": map[string]any{
			"installed":   apps.InstalledApps,
			"snapshot_at": apps.SnapshotTime,
		},
	}
	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
