package diagnostics

import (
	"net"
	"time"
)

func CheckPing(host string, timeout time.Duration) (time.Duration, error) {
	start := time.Now()

	conn, err := net.DialTimeout("tcp", host, timeout)
	if err != nil {
		return 0, err
	}

	conn.Close()

	return time.Since(start), nil
}

func TestCheckPing(url string) {
	latency, err := CheckPing(url, 2*time.Second)

	if err != nil {
		println("ping failed:", err.Error())
		return
	}

	println("latency:", latency.String())
}
