package database

import (
	"database/sql"
	"encoding/json"
	"fmt"

	_ "modernc.org/sqlite"
)

type Database struct {
	db *sql.DB
}

func Open(path string) (*Database, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	if _, err := db.Exec(schema); err != nil {
		db.Close()
		return nil, err
	}

	return &Database{db: db}, nil
}

func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) StoreHealth(h Health) error {
	_, err := d.db.Exec(`
		INSERT INTO health (
			latency_ms,
			packet_loss,
			dns_ms,
			http_ms,
			http_success
		)
		VALUES (?, ?, ?, ?, ?)
	`,
		h.LatencyMs,
		h.PacketLoss,
		h.DNSMs,
		h.HTTPMs,
		h.HTTPSuccess,
	)

	return err
}

func (d *Database) StoreDiagnostic(data DiagnosticData) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	var incidentID any

	if data.Incident != nil {
		result, err := tx.Exec(`
			INSERT INTO incidents (
				reason,
				severity,
				status
			)
			VALUES (?, ?, ?)
		`,
			data.Incident.Reason,
			data.Incident.Severity,
			data.Incident.Status,
		)

		if err != nil {
			return err
		}

		incidentID, err = result.LastInsertId()
		if err != nil {
			return err
		}
	}

	healthID, err := insertHealth(tx, data.Health)
	if err != nil {
		return err
	}

	var wifiID any

	if data.HasWiFi {
		id, err := insertWiFi(tx, data.WiFi)
		if err != nil {
			return err
		}

		wifiID = id
	}

	var speedTestID any

	if data.HasSpeedTest {
		id, err := insertSpeedTest(tx, data.SpeedTest)
		if err != nil {
			return err
		}

		speedTestID = id
	}

	var tracerouteID any

	if data.HasTraceroute {
		id, err := insertTraceroute(tx, data.Traceroute)
		if err != nil {
			return err
		}

		tracerouteID = id
	}

	_, err = tx.Exec(`
		INSERT INTO diagnostics (
			incident_id,
			health_id,
			wifi_id,
			speed_test_id,
			traceroute_id,
			trigger,
			reason,
			status,
			duration_ms
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		incidentID,
		healthID,
		wifiID,
		speedTestID,
		tracerouteID,
		data.Diagnostic.Trigger,
		data.Diagnostic.Reason,
		data.Diagnostic.Status,
		data.Diagnostic.DurationMs,
	)

	if err != nil {
		return err
	}

	return tx.Commit()
}

func insertHealth(tx *sql.Tx, h Health) (int64, error) {
	result, err := tx.Exec(`
		INSERT INTO health (
			latency_ms,
			packet_loss,
			dns_ms,
			http_ms,
			http_success
		)
		VALUES (?, ?, ?, ?, ?)
	`,
		h.LatencyMs,
		h.PacketLoss,
		h.DNSMs,
		h.HTTPMs,
		h.HTTPSuccess,
	)

	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func insertWiFi(tx *sql.Tx, w WiFi) (int64, error) {
	result, err := tx.Exec(`
		INSERT INTO wifi (
			rssi_dbm,
			snr_db,
			channel,
			frequency_mhz,
			band,
			bssid_hash,
			ssid_hash,
			signal_percent,
			link_mbps
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		w.RSSIDBm,
		w.SNRDb,
		w.Channel,
		w.FrequencyMHz,
		w.Band,
		w.BSSIDHash,
		w.SSIDHash,
		w.SignalPercent,
		w.LinkMbps,
	)

	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func insertSpeedTest(tx *sql.Tx, s SpeedTest) (int64, error) {
	result, err := tx.Exec(`
		INSERT INTO speed_tests (
			download_mbps,
			upload_mbps,
			baseline_latency_ms,
			loaded_latency_ms,
			peak_latency_ms,
			bufferbloat_ms
		)
		VALUES (?, ?, ?, ?, ?, ?)
	`,
		s.DownloadMbps,
		s.UploadMbps,
		s.BaselineLatencyMs,
		s.LoadedLatencyMs,
		s.PeakLatencyMs,
		s.BufferbloatMs,
	)

	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func insertTraceroute(tx *sql.Tx, t Traceroute) (int64, error) {
	hops, err := json.Marshal(t.Hops)
	if err != nil {
		return 0, err
	}

	result, err := tx.Exec(`
		INSERT INTO traceroutes (
			destination,
			hops_json
		)
		VALUES (?, ?)
	`,
		t.Destination,
		string(hops),
	)

	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func (d *Database) SetConfig(key, value string) error {
	_, err := d.db.Exec(`
		INSERT INTO config (key, value)
		VALUES (?, ?)
		ON CONFLICT(key)
		DO UPDATE SET value = excluded.value
	`, key, value)

	return err
}

func (d *Database) GetConfig(key string) (string, error) {
	var value string

	err := d.db.QueryRow(`
		SELECT value
		FROM config
		WHERE key = ?
	`, key).Scan(&value)

	return value, err
}

func (d *Database) HealthCount() (int, error) {
	return d.count("health")
}

func (d *Database) DiagnosticCount() (int, error) {
	return d.count("diagnostics")
}

func (d *Database) IncidentCount() (int, error) {
	return d.count("incidents")
}

func (d *Database) WiFiCount() (int, error) {
	return d.count("wifi")
}

func (d *Database) SpeedTestCount() (int, error) {
	return d.count("speed_tests")
}

func (d *Database) TracerouteCount() (int, error) {
	return d.count("traceroutes")
}

func (d *Database) count(table string) (int, error) {
	var count int

	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", table)

	err := d.db.QueryRow(query).Scan(&count)

	return count, err
}


func (d *Database) StoreWiFi(w WiFi) error {
	_, err := d.db.Exec(`
		INSERT INTO wifi (
			rssi_dbm,
			snr_db,
			channel,
			frequency_mhz,
			band,
			bssid_hash,
			ssid_hash,
			signal_percent,
			link_mbps
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		w.RSSIDBm,
		w.SNRDb,
		w.Channel,
		w.FrequencyMHz,
		w.Band,
		w.BSSIDHash,
		w.SSIDHash,
		w.SignalPercent,
		w.LinkMbps,
	)

	return err
}

func (d *Database) StoreSpeedTest(s SpeedTest) error {
	_, err := d.db.Exec(`
		INSERT INTO speed_tests (
			download_mbps,
			upload_mbps,
			baseline_latency_ms,
			loaded_latency_ms,
			peak_latency_ms,
			bufferbloat_ms
		)
		VALUES (?, ?, ?, ?, ?, ?)
	`,
		s.DownloadMbps,
		s.UploadMbps,
		s.BaselineLatencyMs,
		s.LoadedLatencyMs,
		s.PeakLatencyMs,
		s.BufferbloatMs,
	)

	return err
}
