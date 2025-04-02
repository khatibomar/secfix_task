-- name: InsertOSSnapshot :one
INSERT INTO info_history (os_version, os_query_version)
VALUES ($1, $2)
RETURNING id, os_version, os_query_version, created_at;

-- name: GetLatestOSInfo :one
SELECT id, os_version, os_query_version, created_at 
FROM info_history 
ORDER BY created_at DESC 
LIMIT 1;

-- name: InsertAppsSnapshot :one
INSERT INTO system_snapshots (installed_apps)
VALUES ($1)
RETURNING id, snapshot_time, installed_apps;

-- name: GetLatestAppsSnapshot :one
SELECT id, snapshot_time, installed_apps 
FROM system_snapshots 
ORDER BY snapshot_time DESC 
LIMIT 1;
