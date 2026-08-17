package speedtest

import (
	"testing"
	"time"
)

const testServer = "proof.ovh.net"

func TestDownload(t *testing.T) {
	speed, err := Download(
		testServer,
		10*time.Second,
		4,
	)

	if err != nil {
		t.Fatal(err)
	}

	t.Logf("download: %.2f Mbps", speed)

	if speed <= 0 {
		t.Fatalf("invalid download speed: %.2f Mbps", speed)
	}
}

func TestUpload(t *testing.T) {
	speed, err := Upload(
		testServer,
		10*time.Second,
		4,
	)

	if err != nil {
		t.Fatal(err)
	}

	t.Logf("upload: %.2f Mbps", speed)

	if speed <= 0 {
		t.Fatalf("invalid upload speed: %.2f Mbps", speed)
	}
}