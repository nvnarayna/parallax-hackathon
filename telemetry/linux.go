//go:build linux

package telemetry

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func collectWiFi() (WiFiInfo, error) {
	output, err := exec.Command("sh", "-c", "iw dev | awk '$1==\"Interface\"{print $2; exit}'").Output()
	if err != nil {
		return WiFiInfo{}, err
	}

	iface := strings.TrimSpace(string(output))

	if iface == "" {
		return WiFiInfo{}, fmt.Errorf("no Wi-Fi interface found")
	}

	output, err = exec.Command("iw", "dev", iface, "link").Output()
	if err != nil {
		return WiFiInfo{}, err
	}

	lines := strings.Split(string(output), "\n")

	var info WiFiInfo

	for _, line := range lines {
		line = strings.TrimSpace(line)

		switch {
		case strings.HasPrefix(line, "Connected to "):
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				info.BSSID = parts[2]
			}

		case strings.HasPrefix(line, "SSID:"):
			info.SSID = strings.TrimSpace(strings.TrimPrefix(line, "SSID:"))

		case strings.HasPrefix(line, "freq:"):
			value := strings.TrimSpace(strings.TrimPrefix(line, "freq:"))

			if frequency, err := strconv.Atoi(value); err == nil {
				info.Frequency = &frequency

				channel := frequencyToChannel(frequency)
				if channel > 0 {
					info.Channel = &channel
				}

				info.Band = frequencyToBand(frequency)
			}

		case strings.HasPrefix(line, "signal:"):
			fields := strings.Fields(line)

			if len(fields) >= 2 {
				if rssi, err := strconv.Atoi(fields[1]); err == nil {
					info.RSSI = &rssi
				}
			}
		}
	}

	if info.BSSID == "" {
		return WiFiInfo{}, fmt.Errorf("not connected to Wi-Fi")
	}

	return info, nil
}

func frequencyToBand(frequency int) string {
	switch {
	case frequency >= 2400 && frequency < 2500:
		return "2.4GHz"
	case frequency >= 4900 && frequency < 5900:
		return "5GHz"
	case frequency >= 5900 && frequency < 7200:
		return "6GHz"
	default:
		return ""
	}
}

func frequencyToChannel(frequency int) int {
	if frequency >= 2412 && frequency <= 2472 {
		return (frequency - 2407) / 5
	}

	if frequency >= 5000 && frequency <= 5900 {
		return (frequency - 5000) / 5
	}

	if frequency >= 5955 && frequency <= 7115 {
		return (frequency - 5950) / 5
	}

	return 0
}