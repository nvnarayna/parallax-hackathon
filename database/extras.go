package database

import (
	"encoding/json"
	"fmt"
)

type HealthRecord struct {
	ID          int64   `json:"id"`
	Timestamp   int64   `json:"timestamp"`
	LatencyMs   float64 `json:"latency_ms"`
	PacketLoss  float64 `json:"packet_loss"`
	DNSMs       float64 `json:"dns_ms"`
	HTTPMs      float64 `json:"http_ms"`
	HTTPSuccess bool    `json:"http_success"`
}

type WiFiRecord struct {
	ID            int64   `json:"id"`
	Timestamp     int64   `json:"timestamp"`
	RSSIDBm       int     `json:"rssi_dbm"`
	SNRDb         float64 `json:"snr_db"`
	Channel       int     `json:"channel"`
	FrequencyMHz  int     `json:"frequency_mhz"`
	Band          string  `json:"band"`
	BSSIDHash     string  `json:"bssid_hash"`
	SSIDHash      string  `json:"ssid_hash"`
	SignalPercent int     `json:"signal_percent"`
	LinkMbps      float64 `json:"link_mbps"`
}

type SpeedTestRecord struct {
	ID                int64   `json:"id"`
	Timestamp         int64   `json:"timestamp"`
	DownloadMbps      float64 `json:"download_mbps"`
	UploadMbps        float64 `json:"upload_mbps"`
	BaselineLatencyMs float64 `json:"baseline_latency_ms"`
	LoadedLatencyMs   float64 `json:"loaded_latency_ms"`
	PeakLatencyMs     float64 `json:"peak_latency_ms"`
	BufferbloatMs     float64 `json:"bufferbloat_ms"`
}

type TracerouteRecord struct {
	ID          int64  `json:"id"`
	Timestamp   int64  `json:"timestamp"`
	Destination string `json:"destination"`
	Hops        any    `json:"hops"`
}

type IncidentRecord struct {
	ID        int64  `json:"id"`
	Timestamp int64  `json:"timestamp"`
	Reason    string `json:"reason"`
	Severity  string `json:"severity"`
	Status    string `json:"status"`
}

type DiagnosticRecord struct {
	ID              int64  `json:"id"`
	Timestamp       int64  `json:"timestamp"`
	IncidentID      *int64 `json:"incident_id,omitempty"`
	HealthID        int64  `json:"health_id"`
	WiFiID          *int64 `json:"wifi_id,omitempty"`
	SpeedTestID     *int64 `json:"speed_test_id,omitempty"`
	TracerouteID    *int64 `json:"traceroute_id,omitempty"`
	Trigger         string `json:"trigger"`
	Reason          string `json:"reason"`
	Status          string `json:"status"`
	DurationMs      int64  `json:"duration_ms"`
}

func (d *Database) LatestHealth() (HealthRecord, error) {
	var r HealthRecord

	err := d.db.QueryRow(`
		SELECT
			id,
			timestamp,
			latency_ms,
			packet_loss,
			dns_ms,
			http_ms,
			http_success
		FROM health
		ORDER BY id DESC
		LIMIT 1
	`).Scan(
		&r.ID,
		&r.Timestamp,
		&r.LatencyMs,
		&r.PacketLoss,
		&r.DNSMs,
		&r.HTTPMs,
		&r.HTTPSuccess,
	)

	return r, err
}

func (d *Database) HealthHistory(limit int) ([]HealthRecord, error) {
	rows, err := d.db.Query(`
		SELECT
			id,
			timestamp,
			latency_ms,
			packet_loss,
			dns_ms,
			http_ms,
			http_success
		FROM health
		ORDER BY id DESC
		LIMIT ?
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []HealthRecord

	for rows.Next() {
		var r HealthRecord

		if err := rows.Scan(
			&r.ID,
			&r.Timestamp,
			&r.LatencyMs,
			&r.PacketLoss,
			&r.DNSMs,
			&r.HTTPMs,
			&r.HTTPSuccess,
		); err != nil {
			return nil, err
		}

		results = append(results, r)
	}

	return results, rows.Err()
}

func (d *Database) LatestWiFi() (WiFiRecord, error) {
	var r WiFiRecord

	err := d.db.QueryRow(`
		SELECT
			id,
			timestamp,
			rssi_dbm,
			snr_db,
			channel,
			frequency_mhz,
			band,
			bssid_hash,
			ssid_hash,
			signal_percent,
			link_mbps
		FROM wifi
		ORDER BY id DESC
		LIMIT 1
	`).Scan(
		&r.ID,
		&r.Timestamp,
		&r.RSSIDBm,
		&r.SNRDb,
		&r.Channel,
		&r.FrequencyMHz,
		&r.Band,
		&r.BSSIDHash,
		&r.SSIDHash,
		&r.SignalPercent,
		&r.LinkMbps,
	)

	return r, err
}

func (d *Database) WiFiHistory(limit int) ([]WiFiRecord, error) {
	rows, err := d.db.Query(`
		SELECT
			id,
			timestamp,
			rssi_dbm,
			snr_db,
			channel,
			frequency_mhz,
			band,
			bssid_hash,
			ssid_hash,
			signal_percent,
			link_mbps
		FROM wifi
		ORDER BY id DESC
		LIMIT ?
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []WiFiRecord

	for rows.Next() {
		var r WiFiRecord

		if err := rows.Scan(
			&r.ID,
			&r.Timestamp,
			&r.RSSIDBm,
			&r.SNRDb,
			&r.Channel,
			&r.FrequencyMHz,
			&r.Band,
			&r.BSSIDHash,
			&r.SSIDHash,
			&r.SignalPercent,
			&r.LinkMbps,
		); err != nil {
			return nil, err
		}

		results = append(results, r)
	}

	return results, rows.Err()
}

func (d *Database) LatestSpeedTest() (SpeedTestRecord, error) {
	var r SpeedTestRecord

	err := d.db.QueryRow(`
		SELECT
			id,
			timestamp,
			download_mbps,
			upload_mbps,
			baseline_latency_ms,
			loaded_latency_ms,
			peak_latency_ms,
			bufferbloat_ms
		FROM speed_tests
		ORDER BY id DESC
		LIMIT 1
	`).Scan(
		&r.ID,
		&r.Timestamp,
		&r.DownloadMbps,
		&r.UploadMbps,
		&r.BaselineLatencyMs,
		&r.LoadedLatencyMs,
		&r.PeakLatencyMs,
		&r.BufferbloatMs,
	)

	return r, err
}

func (d *Database) SpeedTestHistory(limit int) ([]SpeedTestRecord, error) {
	rows, err := d.db.Query(`
		SELECT
			id,
			timestamp,
			download_mbps,
			upload_mbps,
			baseline_latency_ms,
			loaded_latency_ms,
			peak_latency_ms,
			bufferbloat_ms
		FROM speed_tests
		ORDER BY id DESC
		LIMIT ?
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []SpeedTestRecord

	for rows.Next() {
		var r SpeedTestRecord

		if err := rows.Scan(
			&r.ID,
			&r.Timestamp,
			&r.DownloadMbps,
			&r.UploadMbps,
			&r.BaselineLatencyMs,
			&r.LoadedLatencyMs,
			&r.PeakLatencyMs,
			&r.BufferbloatMs,
		); err != nil {
			return nil, err
		}

		results = append(results, r)
	}

	return results, rows.Err()
}

func (d *Database) LatestDiagnostic() (DiagnosticRecord, error) {
	var r DiagnosticRecord

	err := d.db.QueryRow(`
		SELECT
			id,
			timestamp,
			incident_id,
			health_id,
			wifi_id,
			speed_test_id,
			traceroute_id,
			trigger,
			reason,
			status,
			duration_ms
		FROM diagnostics
		ORDER BY id DESC
		LIMIT 1
	`).Scan(
		&r.ID,
		&r.Timestamp,
		&r.IncidentID,
		&r.HealthID,
		&r.WiFiID,
		&r.SpeedTestID,
		&r.TracerouteID,
		&r.Trigger,
		&r.Reason,
		&r.Status,
		&r.DurationMs,
	)

	return r, err
}

func (d *Database) DiagnosticHistory(limit int) ([]DiagnosticRecord, error) {
	rows, err := d.db.Query(`
		SELECT
			id,
			timestamp,
			incident_id,
			health_id,
			wifi_id,
			speed_test_id,
			traceroute_id,
			trigger,
			reason,
			status,
			duration_ms
		FROM diagnostics
		ORDER BY id DESC
		LIMIT ?
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []DiagnosticRecord

	for rows.Next() {
		var r DiagnosticRecord

		if err := rows.Scan(
			&r.ID,
			&r.Timestamp,
			&r.IncidentID,
			&r.HealthID,
			&r.WiFiID,
			&r.SpeedTestID,
			&r.TracerouteID,
			&r.Trigger,
			&r.Reason,
			&r.Status,
			&r.DurationMs,
		); err != nil {
			return nil, err
		}

		results = append(results, r)
	}

	return results, rows.Err()
}

func (d *Database) GetDiagnostic(id int64) (DiagnosticRecord, error) {
	var r DiagnosticRecord

	err := d.db.QueryRow(`
		SELECT
			id,
			timestamp,
			incident_id,
			health_id,
			wifi_id,
			speed_test_id,
			traceroute_id,
			trigger,
			reason,
			status,
			duration_ms
		FROM diagnostics
		WHERE id = ?
	`, id).Scan(
		&r.ID,
		&r.Timestamp,
		&r.IncidentID,
		&r.HealthID,
		&r.WiFiID,
		&r.SpeedTestID,
		&r.TracerouteID,
		&r.Trigger,
		&r.Reason,
		&r.Status,
		&r.DurationMs,
	)

	return r, err
}

func (d *Database) GetDiagnosticHealth(id int64) (HealthRecord, error) {
	var r HealthRecord

	err := d.db.QueryRow(`
		SELECT
			h.id,
			h.timestamp,
			h.latency_ms,
			h.packet_loss,
			h.dns_ms,
			h.http_ms,
			h.http_success
		FROM diagnostics d
		JOIN health h ON h.id = d.health_id
		WHERE d.id = ?
	`, id).Scan(
		&r.ID,
		&r.Timestamp,
		&r.LatencyMs,
		&r.PacketLoss,
		&r.DNSMs,
		&r.HTTPMs,
		&r.HTTPSuccess,
	)

	return r, err
}

func (d *Database) GetDiagnosticWiFi(id int64) (WiFiRecord, error) {
	var r WiFiRecord

	err := d.db.QueryRow(`
		SELECT
			w.id,
			w.timestamp,
			w.rssi_dbm,
			w.snr_db,
			w.channel,
			w.frequency_mhz,
			w.band,
			w.bssid_hash,
			w.ssid_hash,
			w.signal_percent,
			w.link_mbps
		FROM diagnostics d
		JOIN wifi w ON w.id = d.wifi_id
		WHERE d.id = ?
	`, id).Scan(
		&r.ID,
		&r.Timestamp,
		&r.RSSIDBm,
		&r.SNRDb,
		&r.Channel,
		&r.FrequencyMHz,
		&r.Band,
		&r.BSSIDHash,
		&r.SSIDHash,
		&r.SignalPercent,
		&r.LinkMbps,
	)

	return r, err
}

func (d *Database) GetDiagnosticSpeedTest(id int64) (SpeedTestRecord, error) {
	var r SpeedTestRecord

	err := d.db.QueryRow(`
		SELECT
			s.id,
			s.timestamp,
			s.download_mbps,
			s.upload_mbps,
			s.baseline_latency_ms,
			s.loaded_latency_ms,
			s.peak_latency_ms,
			s.bufferbloat_ms
		FROM diagnostics d
		JOIN speed_tests s ON s.id = d.speed_test_id
		WHERE d.id = ?
	`, id).Scan(
		&r.ID,
		&r.Timestamp,
		&r.DownloadMbps,
		&r.UploadMbps,
		&r.BaselineLatencyMs,
		&r.LoadedLatencyMs,
		&r.PeakLatencyMs,
		&r.BufferbloatMs,
	)

	return r, err
}

func (d *Database) GetDiagnosticTraceroute(id int64) (TracerouteRecord, error) {
	var r TracerouteRecord
	var hopsJSON string

	err := d.db.QueryRow(`
		SELECT
			t.id,
			t.timestamp,
			t.destination,
			t.hops_json
		FROM diagnostics d
		JOIN traceroutes t ON t.id = d.traceroute_id
		WHERE d.id = ?
	`, id).Scan(
		&r.ID,
		&r.Timestamp,
		&r.Destination,
		&hopsJSON,
	)
	if err != nil {
		return r, err
	}

	if err := json.Unmarshal([]byte(hopsJSON), &r.Hops); err != nil {
		return r, err
	}

	return r, nil
}

func (d *Database) GetDiagnosticIncident(id int64) (IncidentRecord, error) {
	var r IncidentRecord

	err := d.db.QueryRow(`
		SELECT
			i.id,
			i.timestamp,
			i.reason,
			i.severity,
			i.status
		FROM diagnostics d
		JOIN incidents i ON i.id = d.incident_id
		WHERE d.id = ?
	`, id).Scan(
		&r.ID,
		&r.Timestamp,
		&r.Reason,
		&r.Severity,
		&r.Status,
	)

	return r, err
}

func (d *Database) DiagnosticBundle(id int64) (DiagnosticBundle, error) {
	diagnostic, err := d.GetDiagnostic(id)
	if err != nil {
		return DiagnosticBundle{}, err
	}

	bundle := DiagnosticBundle{
		Diagnostic: diagnostic,
	}

	bundle.Health, err = d.GetDiagnosticHealth(id)
	if err != nil {
		return DiagnosticBundle{}, err
	}

	if diagnostic.WiFiID != nil {
		wifi, err := d.GetDiagnosticWiFi(id)
		if err != nil {
			return DiagnosticBundle{}, err
		}

		bundle.WiFi = &wifi
	}

	if diagnostic.SpeedTestID != nil {
		speed, err := d.GetDiagnosticSpeedTest(id)
		if err != nil {
			return DiagnosticBundle{}, err
		}

		bundle.SpeedTest = &speed
	}

	if diagnostic.TracerouteID != nil {
		trace, err := d.GetDiagnosticTraceroute(id)
		if err != nil {
			return DiagnosticBundle{}, err
		}

		bundle.Traceroute = &trace
	}

	if diagnostic.IncidentID != nil {
		incident, err := d.GetDiagnosticIncident(id)
		if err != nil {
			return DiagnosticBundle{}, err
		}

		bundle.Incident = &incident
	}

	return bundle, nil
}

type DiagnosticBundle struct {
	Diagnostic  DiagnosticRecord `json:"diagnostic"`
	Incident    *IncidentRecord  `json:"incident,omitempty"`
	Health      HealthRecord     `json:"health"`
	WiFi        *WiFiRecord      `json:"wifi,omitempty"`
	SpeedTest   *SpeedTestRecord `json:"speed_test,omitempty"`
	Traceroute  *TracerouteRecord `json:"traceroute,omitempty"`
}

func (d *Database) Counts() (map[string]int, error) {
	tables := []string{
		"health",
		"wifi",
		"speed_tests",
		"traceroutes",
		"incidents",
		"diagnostics",
	}

	result := make(map[string]int)

	for _, table := range tables {
		count, err := d.count(table)
		if err != nil {
			return nil, err
		}

		result[table] = count
	}

	return result, nil
}

func (d *Database) ValidateLimit(limit int) int {
	if limit < 1 {
		return 1
	}

	if limit > 1000 {
		return 1000
	}

	return limit
}

func (d *Database) Validate() error {
	var result int

	err := d.db.QueryRow(`SELECT 1`).Scan(&result)
	if err != nil {
		return fmt.Errorf("database validation failed: %w", err)
	}

	return nil
}