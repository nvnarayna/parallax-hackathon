//go:build darwin

package telemetry

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func collectWiFi() (WiFiInfo, error) {
	output, err := exec.Command("wdutil", "info").CombinedOutput()
	if err != nil {
		return WiFiInfo{}, fmt.Errorf("wdutil failed: %w: %s", err, string(output))
	}

	var info WiFiInfo

	lines := strings.Split(string(output), "\n")

	inWiFi := false

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if line == "WIFI" {
			inWiFi = true
			continue
		}

		if !inWiFi || line == "" {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "SSID":
			info.SSID = value

		case "BSSID":
			info.BSSID = value

		case "RSSI":
			value = strings.TrimSpace(strings.TrimSuffix(value, "dBm"))

			if rssi, err := strconv.Atoi(value); err == nil {
				info.RSSI = &rssi
			}

		case "Noise":
			value = strings.TrimSpace(strings.TrimSuffix(value, "dBm"))

			if noise, err := strconv.ParseFloat(value, 64); err == nil && info.RSSI != nil {
				snr := float64(*info.RSSI) - noise
				info.SNR = &snr
			}

		case "Channel":
			info.Band = value

			channel := parseMacChannel(value)
			if channel > 0 {
				info.Channel = &channel
			}
		}
	}

	if info.RSSI == nil {
		return WiFiInfo{}, fmt.Errorf("Wi-Fi information unavailable")
	}

	return info, nil
}

func parseMacChannel(value string) int {
	value = strings.TrimSpace(value)

	var channel int

	_, err := fmt.Sscanf(value, "%*[^0-9]%d", &channel)

	if err != nil {
		return 0
	}

	return channel
}