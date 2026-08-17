//go:build windows

package telemetry

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func collectWiFi() (WiFiInfo, error) {
	output, err := exec.Command("netsh", "wlan", "show", "interfaces").Output()
	if err != nil {
		return WiFiInfo{}, err
	}

	var info WiFiInfo

	lines := strings.Split(string(output), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)

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

		case "Signal":
			value = strings.TrimSuffix(value, "%")
			if signal, err := strconv.Atoi(strings.TrimSpace(value)); err == nil {
				info.RSSI = &signal
			}

		case "Channel":
			if channel, err := strconv.Atoi(value); err == nil {
				info.Channel = &channel
			}

		case "Receive rate (Mbps)":
			if rate, err := strconv.ParseFloat(value, 64); err == nil {
				info.LinkMbps = &rate
			}
		}
	}

	if info.SSID == "" {
		return WiFiInfo{}, fmt.Errorf("no active Wi-Fi connection found")
	}

	return info, nil
}