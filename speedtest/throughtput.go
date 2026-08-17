package speedtest

import (
	"time"

	iperf "github.com/BGrewell/go-iperf"
)

func Download(server string, duration time.Duration, streams int) (float64, error) {
	client := iperf.NewClient(server)

	client.SetJSON(true)
	client.SetStreams(streams)
	client.SetTimeSec(int(duration.Seconds()))
	client.SetIncludeServer(true)

	client.SetReverse(true)

	err := client.Start()
	if err != nil {
		return 0, err
	}

	<-client.Done

	report := client.Report()

	return report.End.SumReceived.BitsPerSecond / 1_000_000, nil
}

func Upload(server string, duration time.Duration, streams int) (float64, error) {
	client := iperf.NewClient(server)

	client.SetJSON(true)
	client.SetStreams(streams)
	client.SetTimeSec(int(duration.Seconds()))
	client.SetIncludeServer(true)

	err := client.Start()
	if err != nil {
		return 0, err
	}

	<-client.Done

	report := client.Report()

	return report.End.SumSent.BitsPerSecond / 1_000_000, nil
}
