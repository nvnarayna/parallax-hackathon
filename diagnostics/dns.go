package diagnostics

import (
	"context"
	"net"
	"time"
)

func CheckDNS(host string) (time.Duration, error) {
	start := time.Now()

	_, err := net.DefaultResolver.LookupHost(
		context.Background(),
		host,
	)

	return time.Since(start), err
}
