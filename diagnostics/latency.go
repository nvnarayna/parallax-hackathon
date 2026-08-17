package diagnostics

import (
	"net"
	"time"
)

func CheckLatency(host string, timeout time.Duration) (time.Duration, error) {
	start := time.Now()

	conn, err := net.DialTimeout("tcp", host, timeout)
	if err != nil {
		return 0, err
	}

	conn.Close()

	return time.Since(start), nil
}
