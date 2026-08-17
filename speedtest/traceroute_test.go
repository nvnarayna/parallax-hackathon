package speedtest

import (
	"testing"
	"time"
)

func TestTraceroute(t *testing.T) {
	hops, err := Traceroute(
		"8.8.8.8",
		30,
		2*time.Second,
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(hops) == 0 {
		t.Fatal("traceroute returned no hops")
	}

	for _, hop := range hops {
		if hop.Address == "" {
			t.Logf(
				"hop %d: timeout",
				hop.TTL,
			)
			continue
		}

		t.Logf(
			"hop %d: %s - %v",
			hop.TTL,
			hop.Address,
			hop.Latency,
		)
	}
}