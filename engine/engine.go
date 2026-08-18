package engine

import (
	"context"
	"log"
	"sync"
	"time"

	"wifi-diagnostic/database"
	"wifi-diagnostic/speedtest"
	"wifi-diagnostic/telemetry"
)

type HealthProbe func() (database.Health, error)
type WiFiProbe func() (telemetry.WiFiInfo, error)
type DownloadProbe func() (float64, error)
type UploadProbe func() (float64, error)
type BufferbloatProbe func() (speedtest.BufferbloatResult, error)
type TracerouteProbe func() (database.Traceroute, error)

type Store interface {
	StoreHealth(database.Health) error
	StoreWiFi(database.WiFi) error
	StoreSpeedTest(database.SpeedTest) error
	StoreDiagnostic(database.DiagnosticData) error
}

type Engine struct {
	store Store

	healthProbe      HealthProbe
	wifiProbe        WiFiProbe
	downloadProbe    DownloadProbe
	uploadProbe      UploadProbe
	bufferbloatProbe BufferbloatProbe
	tracerouteProbe  TracerouteProbe

	healthInterval    time.Duration
	telemetryInterval time.Duration
	speedTestInterval time.Duration

	ticketClient TicketClient
	logger       *log.Logger

	diagnosticMu sync.Mutex

	mu      sync.RWMutex
	running bool
	cancel  context.CancelFunc
	wg      sync.WaitGroup
}

type Config struct {
	HealthInterval    time.Duration
	TelemetryInterval time.Duration
	SpeedTestInterval time.Duration
}

func New(
	store Store,
	config Config,
	healthProbe HealthProbe,
	wifiProbe WiFiProbe,
	downloadProbe DownloadProbe,
	uploadProbe UploadProbe,
	bufferbloatProbe BufferbloatProbe,
	tracerouteProbe TracerouteProbe,
	ticketClient TicketClient,
	logger *log.Logger,
) *Engine {
	if logger == nil {
		logger = log.Default()
	}

	return &Engine{
		store:              store,
		healthProbe:        healthProbe,
		wifiProbe:          wifiProbe,
		downloadProbe:      downloadProbe,
		uploadProbe:        uploadProbe,
		bufferbloatProbe:   bufferbloatProbe,
		tracerouteProbe:    tracerouteProbe,
		healthInterval:     config.HealthInterval,
		telemetryInterval:  config.TelemetryInterval,
		speedTestInterval:  config.SpeedTestInterval,
		ticketClient:       ticketClient,
		logger:             logger,
	}
}

func (e *Engine) Start(ctx context.Context) {
	e.mu.Lock()

	if e.running {
		e.mu.Unlock()
		return
	}

	e.logger.Println("engine initial run")

	e.runHealth()
	e.runWiFi()
	e.runSpeedTest()

	ctx, cancel := context.WithCancel(ctx)

	e.cancel = cancel
	e.running = true

	e.mu.Unlock()

	e.wg.Add(3)

	go e.healthLoop(ctx)
	go e.telemetryLoop(ctx)
	go e.speedTestLoop(ctx)

	e.logger.Println("engine started")
}

func (e *Engine) Stop() {
	e.mu.Lock()

	if !e.running {
		e.mu.Unlock()
		return
	}

	e.logger.Println("engine stopping")

	e.cancel()
	e.running = false
	e.cancel = nil

	e.mu.Unlock()

	e.wg.Wait()

	e.logger.Println("engine stopped")
}

func (e *Engine) Running() bool {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return e.running
}

func (e *Engine) RunHealth() error {
	return e.runHealth()
}

func (e *Engine) runHealth() error {
	if e.healthProbe == nil {
		return nil
	}

	start := time.Now()

	e.logger.Println("[health] starting")

	health, err := e.healthProbe()
	if err != nil {
		e.logger.Printf("[health] failed: %v", err)
		return err
	}

	if err := e.store.StoreHealth(health); err != nil {
		e.logger.Printf("[health] storage failed: %v", err)
		return err
	}

	e.logger.Printf(
		"[health] completed in %s: latency=%.2fms loss=%.2f%% dns=%.2fms http=%.2fms success=%t",
		time.Since(start),
		health.LatencyMs,
		health.PacketLoss,
		health.DNSMs,
		health.HTTPMs,
		health.HTTPSuccess,
	)

	return nil
}

func (e *Engine) RunWiFi() error {
	return e.runWiFi()
}

func (e *Engine) runWiFi() error {
	if e.wifiProbe == nil {
		return nil
	}

	start := time.Now()

	e.logger.Println("[telemetry] starting")

	info, err := e.wifiProbe()
	if err != nil {
		e.logger.Printf("[telemetry] failed: %v", err)
		return err
	}

	info = telemetry.AnonymizeWiFi(info)
	info.SSID = telemetry.Anonymize(info.SSID)

	wifi := database.WiFi{
		Band:      info.Band,
		BSSIDHash: info.BSSID,
		SSIDHash:  info.SSID,
	}

	if info.RSSI != nil {
		wifi.RSSIDBm = *info.RSSI
	}

	if info.SNR != nil {
		wifi.SNRDb = *info.SNR
	}

	if info.Channel != nil {
		wifi.Channel = *info.Channel
	}

	if info.Frequency != nil {
		wifi.FrequencyMHz = *info.Frequency
	}

	if info.LinkMbps != nil {
		wifi.LinkMbps = *info.LinkMbps
	}

	if info.SignalPercent != nil {
		wifi.SignalPercent = *info.SignalPercent
	}

	if err := e.store.StoreWiFi(wifi); err != nil {
		e.logger.Printf("[telemetry] storage failed: %v", err)
		return err
	}

	e.logger.Printf(
		"[telemetry] completed in %s",
		time.Since(start),
	)

	return nil
}

func (e *Engine) RunSpeedTest() error {
	return e.runSpeedTest()
}

func (e *Engine) runSpeedTest() error {
	start := time.Now()

	e.logger.Println("[speedtest] starting")

	result := database.SpeedTest{}

	if e.downloadProbe != nil {
		download, err := e.downloadProbe()
		if err != nil {
			e.logger.Printf("[speedtest] download failed: %v", err)
			return err
		}

		result.DownloadMbps = download
	}

	if e.uploadProbe != nil {
		upload, err := e.uploadProbe()
		if err != nil {
			e.logger.Printf("[speedtest] upload failed: %v", err)
			return err
		}

		result.UploadMbps = upload
	}

	if e.bufferbloatProbe != nil {
		bufferbloat, err := e.bufferbloatProbe()
		if err != nil {
			e.logger.Printf("[speedtest] bufferbloat failed: %v", err)
			return err
		}

		result.BaselineLatencyMs =
			float64(bufferbloat.BaselineLatency.Microseconds()) / 1000

		result.LoadedLatencyMs =
			float64(bufferbloat.LoadedLatency.Microseconds()) / 1000

		result.PeakLatencyMs =
			float64(bufferbloat.PeakLatency.Microseconds()) / 1000

		result.BufferbloatMs =
			float64(bufferbloat.Increase.Microseconds()) / 1000
	}

	if err := e.store.StoreSpeedTest(result); err != nil {
		e.logger.Printf("[speedtest] storage failed: %v", err)
		return err
	}

	e.logger.Printf(
		"[speedtest] completed in %s: download=%.2fMbps upload=%.2fMbps bufferbloat=%.2fms",
		time.Since(start),
		result.DownloadMbps,
		result.UploadMbps,
		result.BufferbloatMs,
	)

	return nil
}

func (e *Engine) RunDiagnostic(
	trigger string,
	reason string,
	severity string,
) error {
	e.diagnosticMu.Lock()
	defer e.diagnosticMu.Unlock()

	start := time.Now()

	e.logger.Printf(
		"[diagnostic] starting: trigger=%s reason=%s",
		trigger,
		reason,
	)

	data := database.DiagnosticData{
		Diagnostic: database.Diagnostic{
			Trigger: trigger,
			Reason:  reason,
			Status:  "running",
		},
		Incident: &database.Incident{
			Reason:   reason,
			Severity: severity,
			Status:   "open",
		},
	}

	if e.healthProbe != nil {
		health, err := e.healthProbe()
		if err != nil {
			return e.failDiagnostic(start, err)
		}

		data.Health = health
	}

	if e.wifiProbe != nil {
		info, err := e.wifiProbe()
		if err != nil {
			return e.failDiagnostic(start, err)
		}

		info = telemetry.AnonymizeWiFi(info)
		info.SSID = telemetry.Anonymize(info.SSID)

		data.WiFi = database.WiFi{
			Band:      info.Band,
			BSSIDHash: info.BSSID,
			SSIDHash:  info.SSID,
		}

		if info.RSSI != nil {
			data.WiFi.RSSIDBm = *info.RSSI
		}

		if info.SNR != nil {
			data.WiFi.SNRDb = *info.SNR
		}

		if info.Channel != nil {
			data.WiFi.Channel = *info.Channel
		}

		if info.Frequency != nil {
			data.WiFi.FrequencyMHz = *info.Frequency
		}

		if info.LinkMbps != nil {
			data.WiFi.LinkMbps = *info.LinkMbps
		}

		if info.SignalPercent != nil {
			data.WiFi.SignalPercent = *info.SignalPercent
		}

		data.HasWiFi = true
	}

	if e.downloadProbe != nil {
		download, err := e.downloadProbe()
		if err != nil {
			return e.failDiagnostic(start, err)
		}

		data.SpeedTest.DownloadMbps = download
		data.HasSpeedTest = true
	}

	if e.uploadProbe != nil {
		upload, err := e.uploadProbe()
		if err != nil {
			return e.failDiagnostic(start, err)
		}

		data.SpeedTest.UploadMbps = upload
		data.HasSpeedTest = true
	}

	if e.bufferbloatProbe != nil {
		bufferbloat, err := e.bufferbloatProbe()
		if err != nil {
			return e.failDiagnostic(start, err)
		}

		data.SpeedTest.BaselineLatencyMs =
			float64(bufferbloat.BaselineLatency.Microseconds()) / 1000

		data.SpeedTest.LoadedLatencyMs =
			float64(bufferbloat.LoadedLatency.Microseconds()) / 1000

		data.SpeedTest.PeakLatencyMs =
			float64(bufferbloat.PeakLatency.Microseconds()) / 1000

		data.SpeedTest.BufferbloatMs =
			float64(bufferbloat.Increase.Microseconds()) / 1000

		data.HasSpeedTest = true
	}

	if e.tracerouteProbe != nil {
		traceroute, err := e.tracerouteProbe()
		if err != nil {
			return e.failDiagnostic(start, err)
		}

		data.Traceroute = traceroute
		data.HasTraceroute = true
	}

	data.Diagnostic.Status = "completed"
	data.Diagnostic.DurationMs = time.Since(start).Milliseconds()

	if err := e.store.StoreDiagnostic(data); err != nil {
		e.logger.Printf("[diagnostic] storage failed: %v", err)
		return err
	}

	e.logger.Printf(
		"[diagnostic] completed in %dms",
		data.Diagnostic.DurationMs,
	)

	return nil
}

func (e *Engine) failDiagnostic(
	start time.Time,
	err error,
) error {
	e.logger.Printf(
		"[diagnostic] failed after %dms: %v",
		time.Since(start).Milliseconds(),
		err,
	)

	return err
}

func (e *Engine) SubmitTicket(
	ctx context.Context,
	ticket Ticket,
) error {
	if e.ticketClient == nil {
		return ErrTicketClientUnavailable
	}

	e.logger.Println("[ticket] submitting")

	err := e.ticketClient.Create(ctx, ticket)
	if err != nil {
		e.logger.Printf("[ticket] failed: %v", err)
		return err
	}

	e.logger.Println("[ticket] submitted")

	return nil
}

func (e *Engine) healthLoop(ctx context.Context) {
	defer e.wg.Done()

	ticker := time.NewTicker(e.healthInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			_ = e.runHealth()

		case <-ctx.Done():
			return
		}
	}
}

func (e *Engine) telemetryLoop(ctx context.Context) {
	defer e.wg.Done()

	ticker := time.NewTicker(e.telemetryInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			_ = e.runWiFi()

		case <-ctx.Done():
			return
		}
	}
}

func (e *Engine) speedTestLoop(ctx context.Context) {
	defer e.wg.Done()

	ticker := time.NewTicker(e.speedTestInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			_ = e.runSpeedTest()

		case <-ctx.Done():
			return
		}
	}
}