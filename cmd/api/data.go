package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/khatibomar/secfix_challenge/internal/database"
)

var (
	ErrNegativeOrZeroPageNumber = errors.New("page_number should be greater than 0")
	ErrNegativeOrZeroPageSize   = errors.New("page_size should be greater than 0")
)

func (app *application) latestDataHandler(w http.ResponseWriter, r *http.Request) {
	info, err := app.db.GetLatestOSInfo(app.ctx)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	page_number := r.URL.Query().Get("page_number")
	pageNumberCasted, err := strconv.Atoi(page_number)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	pagesize := r.URL.Query().Get("page_size")
	pageSizeCasted, err := strconv.Atoi(pagesize)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if pageNumberCasted <= 0 {
		app.badRequestErrorResponse(w, r, ErrNegativeOrZeroPageNumber, ErrNegativeOrZeroPageNumber.Error())
		return
	}
	if pageSizeCasted <= 0 {
		app.badRequestErrorResponse(w, r, ErrNegativeOrZeroPageSize, ErrNegativeOrZeroPageSize.Error())
		return
	}
	chunkSize := (pageNumberCasted - 1) * pageSizeCasted

	apps, err := app.db.GetLatestAppsSnapshot(app.ctx, database.GetLatestAppsSnapshotParams{
		Column1: int32(chunkSize + 1),
		Column2: int32(chunkSize + pageSizeCasted),
	})
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
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
