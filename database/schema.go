package database

const schema = `
PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS health (
    id INTEGER PRIMARY KEY,
    timestamp INTEGER NOT NULL DEFAULT (unixepoch()),

    latency_ms REAL NOT NULL,
    packet_loss REAL NOT NULL,
    dns_ms REAL NOT NULL,
    http_ms REAL NOT NULL,
    http_success INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS wifi (
    id INTEGER PRIMARY KEY,
    timestamp INTEGER NOT NULL DEFAULT (unixepoch()),

    rssi_dbm INTEGER NOT NULL,
    snr_db REAL NOT NULL,
    channel INTEGER NOT NULL,
    frequency_mhz INTEGER NOT NULL,
    band TEXT NOT NULL,

    bssid_hash TEXT NOT NULL,
    ssid_hash TEXT NOT NULL,

    signal_percent INTEGER NOT NULL,
    link_mbps REAL NOT NULL
);

CREATE TABLE IF NOT EXISTS speed_tests (
    id INTEGER PRIMARY KEY,
    timestamp INTEGER NOT NULL DEFAULT (unixepoch()),

    download_mbps REAL NOT NULL,
    upload_mbps REAL NOT NULL,

    baseline_latency_ms REAL NOT NULL,
    loaded_latency_ms REAL NOT NULL,
    peak_latency_ms REAL NOT NULL,
    bufferbloat_ms REAL NOT NULL
);

CREATE TABLE IF NOT EXISTS traceroutes (
    id INTEGER PRIMARY KEY,
    timestamp INTEGER NOT NULL DEFAULT (unixepoch()),

    destination TEXT NOT NULL,
    hops_json TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS incidents (
    id INTEGER PRIMARY KEY,
    timestamp INTEGER NOT NULL DEFAULT (unixepoch()),

    reason TEXT NOT NULL,
    severity TEXT NOT NULL,
    status TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS diagnostics (
    id INTEGER PRIMARY KEY,
    timestamp INTEGER NOT NULL DEFAULT (unixepoch()),

    incident_id INTEGER,

    health_id INTEGER NOT NULL,
    wifi_id INTEGER,
    speed_test_id INTEGER,
    traceroute_id INTEGER,

    trigger TEXT NOT NULL,
    reason TEXT NOT NULL,
    status TEXT NOT NULL,
    duration_ms INTEGER NOT NULL,

    FOREIGN KEY (incident_id)
        REFERENCES incidents(id),

    FOREIGN KEY (health_id)
        REFERENCES health(id),

    FOREIGN KEY (wifi_id)
        REFERENCES wifi(id),

    FOREIGN KEY (speed_test_id)
        REFERENCES speed_tests(id),

    FOREIGN KEY (traceroute_id)
        REFERENCES traceroutes(id)
);

CREATE TABLE IF NOT EXISTS config (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL
);
`
