//go:build darwin

package telemetry

/*
#cgo LDFLAGS: -framework CoreWLAN -framework Foundation
#include "wifi_corewlan_darwin.h"
*/
import "C"

import "fmt"

func collectWiFi() (WiFiInfo, error) {
	var cInfo C.WiFiInfo

	if C.getWiFiInfo(&cInfo) != 0 {
		return WiFiInfo{}, fmt.Errorf("Wi-Fi information unavailable")
	}

	defer C.freeWiFiInfo(&cInfo)

	info := WiFiInfo{}

	if cInfo.ssid != nil {
		info.SSID = C.GoString(cInfo.ssid)
	}

	if cInfo.bssid != nil {
		info.BSSID = C.GoString(cInfo.bssid)
	}

	if cInfo.rssi_valid != 0 {
		rssi := int(cInfo.rssi)
		info.RSSI = &rssi

		// Convert RSSI to an approximate 0-100% signal value.
		//
		// -30 dBm ≈ excellent
		// -90 dBm ≈ unusable
		percent := (rssi + 90) * 100 / 60

		if percent < 0 {
			percent = 0
		}
		if percent > 100 {
			percent = 100
		}

		info.SignalPercent = &percent
	}

	if cInfo.noise_valid != 0 {
		noise := float64(cInfo.noise)

		if info.RSSI != nil {
			snr := float64(*info.RSSI) - noise
			info.SNR = &snr
		}
	}

	if cInfo.channel_valid != 0 {
		channel := int(cInfo.channel)
		info.Channel = &channel
	}

	if cInfo.frequency_valid != 0 {
		frequency := int(cInfo.frequency)
		info.Frequency = &frequency
	}

	if cInfo.link_mbps_valid != 0 {
		linkMbps := float64(cInfo.link_mbps)
		info.LinkMbps = &linkMbps
	}

	if cInfo.band != nil {
		info.Band = C.GoString(cInfo.band)
	}

	return info, nil
}