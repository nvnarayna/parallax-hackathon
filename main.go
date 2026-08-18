package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"wifi-diagnostic/database"
	"wifi-diagnostic/diagnostics"
	"wifi-diagnostic/engine"
	"wifi-diagnostic/network"
	"wifi-diagnostic/speedtest"
	"wifi-diagnostic/telemetry"
)

func main() {
	logger := log.New(
		os.Stdout,
		"[wifi-diagnostic] ",
		log.LstdFlags,
	)

	db, err := database.Open("wifi-diagnostic.db")
	if err != nil {
		logger.Fatalf("database: %v", err)
	}
	defer db.Close()

	healthProbe := func() (database.Health, error) {
		const pingCount = 3

		pings, err := diagnostics.Ping(
			"8.8.8.8",
			pingCount,
			2*time.Second,
		)
		if err != nil {
			return database.Health{}, err
		}

		var totalLatency time.Duration

		for _, ping := range pings {
			totalLatency += ping.Latency
		}

		var latencyMs float64

		if len(pings) > 0 {
			averageLatency := totalLatency /
				time.Duration(len(pings))

			latencyMs =
				float64(averageLatency.Microseconds()) / 1000
		}

		packetLoss := float64(pingCount-len(pings)) /
			float64(pingCount) * 100

		dnsDuration, err := diagnostics.CheckDNS("google.com")
		if err != nil {
			return database.Health{}, err
		}

		httpDuration, status, err :=
			diagnostics.CheckHTTP("https://google.com")
		if err != nil {
			return database.Health{}, err
		}

		return database.Health{
			LatencyMs:   latencyMs,
			PacketLoss:  packetLoss,
			DNSMs:       float64(dnsDuration.Microseconds()) / 1000,
			HTTPMs:      float64(httpDuration.Microseconds()) / 1000,
			HTTPSuccess: status >= 200 && status < 400,
		}, nil
	}

	downloadProbe := func() (float64, error) {
		return speedtest.Download(
			"proof.ovh.net",
			10*time.Second,
			4,
		)
	}

	uploadProbe := func() (float64, error) {
		return speedtest.Upload(
			"proof.ovh.net",
			10*time.Second,
			4,
		)
	}

	tracerouteProbe := func() (database.Traceroute, error) {
		hops, err := speedtest.Traceroute(
			"8.8.8.8",
			30,
			2*time.Second,
		)
		if err != nil {
			return database.Traceroute{}, err
		}

		return database.Traceroute{
			Destination: "8.8.8.8",
			Hops:        hops,
		}, nil
	}

	config := engine.Config{
		HealthInterval:    30 * time.Second,
		TelemetryInterval: 60 * time.Second,
		SpeedTestInterval: 30 * time.Minute,
	}
	wifiProbe := func() (telemetry.WiFiInfo, error) {
		return telemetry.CollectWiFi()
	}

	bufferbloatProbe := func() (speedtest.BufferbloatResult, error) {
		return speedtest.Bufferbloat(
			"proof.ovh.net",
			"8.8.8.8",
			10*time.Second,
			4,
		)
	}

	eng := engine.New(
		db,
		config,
		healthProbe,
		wifiProbe,
		downloadProbe,
		uploadProbe,
		bufferbloatProbe,
		tracerouteProbe,
		nil,
		logger,
	)

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer cancel()

	eng.Start(ctx)

	logger.Println("engine started")
	logger.Println("http server: http://127.0.0.1:8080")

	server := network.New(
		"127.0.0.1:8080",
		db,
		eng,
		logger,
	)

	serverErr := make(chan error, 1)

	go func() {
		serverErr <- server.Start()
	}()

	select {
	case <-ctx.Done():
		logger.Println("shutdown requested")

	case err := <-serverErr:
		if err != nil {
			logger.Fatalf("http server: %v", err)
		}
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer shutdownCancel()

	if err := server.Stop(shutdownCtx); err != nil {
		logger.Printf("http shutdown: %v", err)
	}

	eng.Stop()

	logger.Println("shutdown complete")
}
