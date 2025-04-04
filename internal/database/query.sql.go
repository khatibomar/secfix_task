// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: query.sql

package database

import (
	"context"
)

const getLatestAppsSnapshot = `-- name: GetLatestAppsSnapshot :one
SELECT id, snapshot_time, installed_apps 
FROM system_snapshots 
ORDER BY snapshot_time DESC 
LIMIT 1
`

func (q *Queries) GetLatestAppsSnapshot(ctx context.Context) (SystemSnapshot, error) {
	row := q.db.QueryRow(ctx, getLatestAppsSnapshot)
	var i SystemSnapshot
	err := row.Scan(&i.ID, &i.SnapshotTime, &i.InstalledApps)
	return i, err
}

const getLatestOSInfo = `-- name: GetLatestOSInfo :one
SELECT id, os_version, os_query_version, created_at 
FROM info_history 
ORDER BY created_at DESC 
LIMIT 1
`

func (q *Queries) GetLatestOSInfo(ctx context.Context) (InfoHistory, error) {
	row := q.db.QueryRow(ctx, getLatestOSInfo)
	var i InfoHistory
	err := row.Scan(
		&i.ID,
		&i.OsVersion,
		&i.OsQueryVersion,
		&i.CreatedAt,
	)
	return i, err
}

const insertAppsSnapshot = `-- name: InsertAppsSnapshot :one
INSERT INTO system_snapshots (installed_apps)
VALUES ($1)
RETURNING id, snapshot_time, installed_apps
`

func (q *Queries) InsertAppsSnapshot(ctx context.Context, installedApps []string) (SystemSnapshot, error) {
	row := q.db.QueryRow(ctx, insertAppsSnapshot, installedApps)
	var i SystemSnapshot
	err := row.Scan(&i.ID, &i.SnapshotTime, &i.InstalledApps)
	return i, err
}

const insertOSSnapshot = `-- name: InsertOSSnapshot :one
INSERT INTO info_history (os_version, os_query_version)
VALUES ($1, $2)
RETURNING id, os_version, os_query_version, created_at
`

type InsertOSSnapshotParams struct {
	OsVersion      string `db:"os_version"`
	OsQueryVersion string `db:"os_query_version"`
}

func (q *Queries) InsertOSSnapshot(ctx context.Context, arg InsertOSSnapshotParams) (InfoHistory, error) {
	row := q.db.QueryRow(ctx, insertOSSnapshot, arg.OsVersion, arg.OsQueryVersion)
	var i InfoHistory
	err := row.Scan(
		&i.ID,
		&i.OsVersion,
		&i.OsQueryVersion,
		&i.CreatedAt,
	)
	return i, err
}
