package engine

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"wifi-diagnostic/database"
	"wifi-diagnostic/diagnostics"
	"wifi-diagnostic/speedtest"
	"wifi-diagnostic/telemetry"
)

func TestEngine(t *testing.T) {
	const databasePath = "engine-test1.db"


	db, err := database.Open(databasePath)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	healthProbe := func() (database.Health, error) {
		const pingCount = 3

		pings, err := diagnostics.Ping(
			"8.8.8.8",
			pingCount,
			2*time.Second,
		)
		if err != nil {
			return database.Health{}, err
		}

		var totalLatency time.Duration

		for _, ping := range pings {
			totalLatency += ping.Latency
		}

		var latencyMs float64

		if len(pings) > 0 {
			averageLatency := totalLatency / time.Duration(len(pings))
			latencyMs = float64(averageLatency.Microseconds()) / 1000
		}

		packetLoss := float64(pingCount-len(pings)) /
			float64(pingCount) * 100

		dnsDuration, err := diagnostics.CheckDNS("google.com")
		if err != nil {
			return database.Health{}, err
		}

		httpDuration, status, err := diagnostics.CheckHTTP(
			"https://google.com",
		)
		if err != nil {
			return database.Health{}, err
		}

		return database.Health{
			LatencyMs:   latencyMs,
			PacketLoss:  packetLoss,
			DNSMs:       float64(dnsDuration.Microseconds()) / 1000,
			HTTPMs:      float64(httpDuration.Microseconds()) / 1000,
			HTTPSuccess: status >= 200 && status < 400,
		}, nil
	}

	wifiProbe := func() (telemetry.WiFiInfo, error) {
		return telemetry.CollectWiFi()
	}

	downloadProbe := func() (float64, error) {
		return speedtest.Download(
			"proof.ovh.net",
			10*time.Second,
			4,
		)
	}

	uploadProbe := func() (float64, error) {
		return speedtest.Upload(
			"proof.ovh.net",
			10*time.Second,
			4,
		)
	}

	bufferbloatProbe := func() (speedtest.BufferbloatResult, error) {
		return speedtest.Bufferbloat(
			"proof.ovh.net",
			"8.8.8.8",
			10*time.Second,
			4,
		)
	}

	tracerouteProbe := func() (database.Traceroute, error) {
		hops, err := speedtest.Traceroute(
			"8.8.8.8",
			30,
			2*time.Second,
		)
		if err != nil {
			return database.Traceroute{}, err
		}

		return database.Traceroute{
			Destination: "8.8.8.8",
			Hops:        hops,
		}, nil
	}

	logger := log.New(
		os.Stdout,
		"[engine-test] ",
		log.LstdFlags,
	)

	engine := New(
		db,
		Config{
			HealthInterval:    30 * time.Second,
			TelemetryInterval: 60 * time.Second,
			SpeedTestInterval: 30 * time.Minute,
		},
		healthProbe,
		wifiProbe,
		downloadProbe,
		uploadProbe,
		bufferbloatProbe,
		tracerouteProbe,
		nil,
		logger,
	)

	t.Log("starting engine")

	engine.Start(context.Background())

	if !engine.Running() {
		t.Fatal("engine is not running")
	}

	t.Log("initial engine run completed")

	t.Log("running full diagnostic")

	err = engine.RunDiagnostic(
		"manual",
		"engine integration test",
		"low",
	)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("full diagnostic completed")

	t.Log("waiting for periodic health run")

	time.Sleep(35 * time.Second)

	engine.Stop()

	if engine.Running() {
		t.Fatal("engine did not stop")
	}

	healthCount, err := db.HealthCount()
	if err != nil {
		t.Fatal(err)
	}

	wifiCount, err := db.WiFiCount()
	if err != nil {
		t.Fatal(err)
	}

	speedCount, err := db.SpeedTestCount()
	if err != nil {
		t.Fatal(err)
	}

	tracerouteCount, err := db.TracerouteCount()
	if err != nil {
		t.Fatal(err)
	}

	diagnosticCount, err := db.DiagnosticCount()
	if err != nil {
		t.Fatal(err)
	}

	incidentCount, err := db.IncidentCount()
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("health records: %d", healthCount)
	t.Logf("wifi records: %d", wifiCount)
	t.Logf("speed test records: %d", speedCount)
	t.Logf("traceroute records: %d", tracerouteCount)
	t.Logf("diagnostic records: %d", diagnosticCount)
	t.Logf("incident records: %d", incidentCount)

	if healthCount < 2 {
		t.Fatalf(
			"expected initial health + periodic health, got %d",
			healthCount,
		)
	}

	if wifiCount < 1 {
		t.Fatalf(
			"expected initial telemetry, got %d",
			wifiCount,
		)
	}

	if speedCount < 1 {
		t.Fatalf(
			"expected initial speed test, got %d",
			speedCount,
		)
	}

	if tracerouteCount < 1 {
		t.Fatalf(
			"expected diagnostic traceroute, got %d",
			tracerouteCount,
		)
	}

	if diagnosticCount < 1 {
		t.Fatalf(
			"expected diagnostic, got %d",
			diagnosticCount,
		)
	}

	if incidentCount < 1 {
		t.Fatalf(
			"expected diagnostic incident, got %d",
			incidentCount,
		)
	}
}
