package speedtest

import (
	"testing"
	"time"
)

func TestBufferbloat(t *testing.T) {
	result, err := Bufferbloat(
		"proof.ovh.net",
		"8.8.8.8",
		10*time.Second,
		4,
	)

	if err != nil {
		t.Fatal(err)
	}

	t.Logf("baseline latency: %v", result.BaselineLatency)
	t.Logf("loaded latency: %v", result.LoadedLatency)
	t.Logf("peak latency: %v", result.PeakLatency)
	t.Logf("bufferbloat: %v", result.Increase)

	if result.BaselineLatency <= 0 {
		t.Fatalf("invalid baseline latency: %v", result.BaselineLatency)
	}

	if result.LoadedLatency <= 0 {
		t.Fatalf("invalid loaded latency: %v", result.LoadedLatency)
	}

	if result.PeakLatency <= 0 {
		t.Fatalf("invalid peak latency: %v", result.PeakLatency)
	}
}