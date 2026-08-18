package database

type Health struct {
	LatencyMs   float64
	PacketLoss  float64
	DNSMs       float64
	HTTPMs      float64
	HTTPSuccess bool
}

type WiFi struct {
	RSSIDBm       int
	SNRDb         float64
	Channel       int
	FrequencyMHz  int
	Band          string
	BSSIDHash     string
	SSIDHash      string
	SignalPercent int
	LinkMbps      float64
}

type SpeedTest struct {
	DownloadMbps      float64
	UploadMbps        float64
	BaselineLatencyMs float64
	LoadedLatencyMs   float64
	PeakLatencyMs     float64
	BufferbloatMs     float64
}

type Traceroute struct {
	Destination string
	Hops        any
}

type Incident struct {
	Reason   string
	Severity string
	Status   string
}

type Diagnostic struct {
	Trigger    string
	Reason     string
	Status     string
	DurationMs int64
}

type DiagnosticData struct {
	Diagnostic Diagnostic

	Incident *Incident

	Health     Health
	WiFi       WiFi
	SpeedTest  SpeedTest
	Traceroute Traceroute

	HasWiFi       bool
	HasSpeedTest  bool
	HasTraceroute bool
}
