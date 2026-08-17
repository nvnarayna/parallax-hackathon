package diagnostics

import (
	"testing"
	"time"
)

func TestCheckLatency(t *testing.T) {
	latency, err := CheckLatency("8.8.8.8:443", 2*time.Second)

	if err != nil {
		t.Fatal(err)
	}

	t.Log("latency:", latency)

	if latency <= 0 {
		t.Fatalf("invalid latency: %v", latency)
	}
}	