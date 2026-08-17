package telemetry

import "time"


type WiFiInfo struct {
	RSSI          *int     `json:"rssi_dbm,omitempty"`
	SNR           *float64 `json:"snr_db,omitempty"`
	Channel       *int     `json:"channel,omitempty"`
	Frequency     *int     `json:"frequency_mhz,omitempty"`
	Band          string   `json:"band,omitempty"`
	BSSID         string   `json:"bssid,omitempty"`
	SSID          string   `json:"ssid,omitempty"`
	LinkMbps      *float64 `json:"link_mbps,omitempty"`
	SignalPercent *int     `json:"signal_percent,omitempty"`
}

func CollectWiFi() (WiFiInfo, error) {
	return collectWiFi()
}
type DeviceInfo struct {
	OS           string `json:"os"`
	Architecture string `json:"architecture"`
	Hostname     string `json:"hostname"`
}

type NetworkEvents struct {
	Connected      bool      `json:"connected"`
	Interface      string    `json:"interface"`
	LastChange     time.Time `json:"last_change"`
	NetworkChanged bool      `json:"network_changed"`
}

type Telemetry struct {
	Timestamp time.Time     `json:"timestamp"`
	WiFi      WiFiInfo      `json:"wifi"`
	Device    DeviceInfo    `json:"device"`
	Network   NetworkEvents `json:"network"`
}
