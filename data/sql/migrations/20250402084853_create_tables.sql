-- migrate:up
CREATE TABLE info_history (
    id SERIAL PRIMARY KEY,
    os_version TEXT NOT NULL,
    os_query_version TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE system_snapshots (
    id SERIAL PRIMARY KEY,
    snapshot_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    installed_apps TEXT[] NOT NULL
);

-- migrate:down
DROP TABLE IF EXISTS system_snapshots;
DROP TABLE IF EXISTS info_history;
