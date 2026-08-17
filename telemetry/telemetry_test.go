package telemetry

import (
	"testing"
)

func TestCollectWiFi(t *testing.T) {
	info, err := CollectWiFi()

	if err != nil {
		t.Fatal(err)
	}

	t.Logf("SSID: %s", info.SSID)
	t.Logf("BSSID: %s", info.BSSID)

	if info.RSSI != nil {
		t.Logf("RSSI: %d dBm", *info.RSSI)
	} else {
		t.Log("RSSI: unavailable")
	}

	if info.SNR != nil {
		t.Logf("SNR: %.2f dB", *info.SNR)
	} else {
		t.Log("SNR: unavailable")
	}

	if info.Channel != nil {
		t.Logf("channel: %d", *info.Channel)
	} else {
		t.Log("channel: unavailable")
	}

	if info.Frequency != nil {
		t.Logf("frequency: %d MHz", *info.Frequency)
	} else {
		t.Log("frequency: unavailable")
	}

	t.Logf("band: %s", info.Band)

	if info.LinkMbps != nil {
		t.Logf("link speed: %.2f Mbps", *info.LinkMbps)
	} else {
		t.Log("link speed: unavailable")
	}
}