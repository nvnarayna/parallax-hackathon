package diagnostics

import (
	"testing"
	"time"
)

func TestCheckPacketLoss(t *testing.T) {
	packetLoss, err := CheckPacketLoss(
		"8.8.8.8",
		10,
		2*time.Second,
	)

	if err != nil {
		t.Fatal(err)
	}

	t.Logf("packet loss: %.2f%%", packetLoss)

	if packetLoss < 0 || packetLoss > 100 {
		t.Fatalf("invalid packet loss: %.2f%%", packetLoss)
	}
}