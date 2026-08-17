package speedtest

import (
	"time"

	"wifi-diagnostic/diagnostics"
)

type BufferbloatResult struct {
	BaselineLatency time.Duration
	LoadedLatency   time.Duration
	PeakLatency     time.Duration
	Increase        time.Duration
}

func Bufferbloat(
	server string,
	pingHost string,
	duration time.Duration,
	streams int,
) (BufferbloatResult, error) {

	baseline, err := diagnostics.Ping(
		pingHost,
		5,
		2*time.Second,
	)

	if err != nil {
		return BufferbloatResult{}, err
	}

	if len(baseline) == 0 {
		return BufferbloatResult{}, errNoPingResponse
	}

	var baselineTotal time.Duration

	for _, response := range baseline {
		baselineTotal += response.Latency
	}

	baselineLatency := baselineTotal / time.Duration(len(baseline))

	type pingResult struct {
		responses []diagnostics.PingResponse
		err       error
	}

	resultChannel := make(chan pingResult, 1)

	go func() {
		responses, err := diagnostics.Ping(
			pingHost,
			int(duration.Seconds()),
			2*time.Second,
		)

		resultChannel <- pingResult{
			responses: responses,
			err:       err,
		}
	}()

	_, err = Download(server, duration, streams)
	if err != nil {
		return BufferbloatResult{}, err
	}

	pings := <-resultChannel

	if pings.err != nil {
		return BufferbloatResult{}, pings.err
	}

	if len(pings.responses) == 0 {
		return BufferbloatResult{}, errNoPingResponse
	}

	var total time.Duration
	peak := time.Duration(0)

	for _, response := range pings.responses {
		total += response.Latency

		if response.Latency > peak {
			peak = response.Latency
		}
	}

	loadedLatency := total / time.Duration(len(pings.responses))

	return BufferbloatResult{
		BaselineLatency: baselineLatency,
		LoadedLatency:   loadedLatency,
		PeakLatency:     peak,
		Increase:        loadedLatency - baselineLatency,
	}, nil
}