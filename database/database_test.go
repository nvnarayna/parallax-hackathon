package database

import (
	"testing"
)

func newTestDatabase(t *testing.T) *Database {
	t.Helper()

	db, err := Open(":memory:") //change to db name
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	return db
}

func testHealth() Health {
	return Health{
		LatencyMs:   25,
		PacketLoss:  0,
		DNSMs:       18,
		HTTPMs:      42,
		HTTPSuccess: true,
	}
}

func testWiFi() WiFi {
	return WiFi{
		RSSIDBm:       -55,
		SNRDb:         35,
		Channel:       36,
		FrequencyMHz:  5180,
		Band:          "5GHz",
		BSSIDHash:     "test-bssid",
		SSIDHash:      "test-ssid",
		SignalPercent: 90,
		LinkMbps:      866,
	}
}

func testSpeedTest() SpeedTest {
	return SpeedTest{
		DownloadMbps:      92,
		UploadMbps:        18,
		BaselineLatencyMs: 22,
		LoadedLatencyMs:   104,
		PeakLatencyMs:     143,
		BufferbloatMs:     82,
	}
}

func testTraceroute() Traceroute {
	return Traceroute{
		Destination: "8.8.8.8",
		Hops: []map[string]any{
			{
				"ttl":        1,
				"address":    "192.168.1.1",
				"latency_ms": 2.1,
			},
			{
				"ttl":        2,
				"address":    "10.0.0.1",
				"latency_ms": 8.4,
			},
		},
	}
}

func TestDatabaseOpen(t *testing.T) {
	db := newTestDatabase(t)

	if err := db.Validate(); err != nil {
		t.Fatal(err)
	}
}

func TestStoreHealth(t *testing.T) {
	db := newTestDatabase(t)

	err := db.StoreHealth(testHealth())
	if err != nil {
		t.Fatal(err)
	}

	err = db.StoreHealth(testHealth())
	if err != nil {
		t.Fatal(err)
	}

	count, err := db.HealthCount()
	if err != nil {
		t.Fatal(err)
	}

	if count != 2 {
		t.Fatalf("expected 2 health records, got %d", count)
	}
}

func TestManualDiagnostic(t *testing.T) {
	db := newTestDatabase(t)

	err := db.StoreDiagnostic(DiagnosticData{
		Diagnostic: Diagnostic{
			Trigger:    "manual",
			Reason:     "user requested diagnostic",
			Status:     "completed",
			DurationMs: 15000,
		},

		Health:     testHealth(),
		WiFi:       testWiFi(),
		SpeedTest:  testSpeedTest(),
		Traceroute: testTraceroute(),

		HasWiFi:       true,
		HasSpeedTest:  true,
		HasTraceroute: true,
	})

	if err != nil {
		t.Fatal(err)
	}

	assertCount(t, db.DiagnosticCount, 1, "diagnostics")
	assertCount(t, db.IncidentCount, 0, "incidents")
	assertCount(t, db.HealthCount, 1, "health")
	assertCount(t, db.WiFiCount, 1, "wifi")
	assertCount(t, db.SpeedTestCount, 1, "speed tests")
	assertCount(t, db.TracerouteCount, 1, "traceroutes")
}

func TestAutomaticDiagnostic(t *testing.T) {
	db := newTestDatabase(t)

	err := db.StoreDiagnostic(DiagnosticData{
		Incident: &Incident{
			Reason:   "packet loss exceeded threshold",
			Severity: "warning",
			Status:   "open",
		},

		Diagnostic: Diagnostic{
			Trigger:    "automatic",
			Reason:     "packet loss exceeded threshold",
			Status:     "completed",
			DurationMs: 12000,
		},

		Health:     testHealth(),
		WiFi:       testWiFi(),
		SpeedTest:  testSpeedTest(),
		Traceroute: testTraceroute(),

		HasWiFi:       true,
		HasSpeedTest:  true,
		HasTraceroute: true,
	})

	if err != nil {
		t.Fatal(err)
	}

	assertCount(t, db.DiagnosticCount, 1, "diagnostics")
	assertCount(t, db.IncidentCount, 1, "incidents")
	assertCount(t, db.HealthCount, 1, "health")
	assertCount(t, db.WiFiCount, 1, "wifi")
	assertCount(t, db.SpeedTestCount, 1, "speed tests")
	assertCount(t, db.TracerouteCount, 1, "traceroutes")
}

func TestRepeatedHealthMonitoring(t *testing.T) {
	db := newTestDatabase(t)

	for i := 0; i < 10; i++ {
		err := db.StoreHealth(testHealth())
		if err != nil {
			t.Fatal(err)
		}
	}

	count, err := db.HealthCount()
	if err != nil {
		t.Fatal(err)
	}

	if count != 10 {
		t.Fatalf("expected 10 health records, got %d", count)
	}

	diagnostics, err := db.DiagnosticCount()
	if err != nil {
		t.Fatal(err)
	}

	if diagnostics != 0 {
		t.Fatalf("expected no diagnostics, got %d", diagnostics)
	}
}

func TestConfig(t *testing.T) {
	db := newTestDatabase(t)

	err := db.SetConfig("latency_threshold_ms", "100")
	if err != nil {
		t.Fatal(err)
	}

	value, err := db.GetConfig("latency_threshold_ms")
	if err != nil {
		t.Fatal(err)
	}

	if value != "100" {
		t.Fatalf("expected 100, got %s", value)
	}

	err = db.SetConfig("latency_threshold_ms", "200")
	if err != nil {
		t.Fatal(err)
	}

	value, err = db.GetConfig("latency_threshold_ms")
	if err != nil {
		t.Fatal(err)
	}

	if value != "200" {
		t.Fatalf("expected 200, got %s", value)
	}
}

func funcCount(fn func() (int, error)) func() (int, error) {
	return fn
}

func assertCount(
	t *testing.T,
	fn func() (int, error),
	expected int,
	name string,
) {
	t.Helper()

	count, err := fn()
	if err != nil {
		t.Fatalf("%s count failed: %v", name, err)
	}

	if count != expected {
		t.Fatalf("expected %d %s, got %d", expected, name, count)
	}
}
