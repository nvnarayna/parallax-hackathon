package diagnostics

import "time"

func CheckPacketLoss(host string, count int, timeout time.Duration) (float64, error) {
	responses, err := Ping(host, count, timeout)

	if err != nil {
		return 0, err
	}

	lost := count - len(responses)

	return float64(lost) / float64(count) * 100, nil
}
