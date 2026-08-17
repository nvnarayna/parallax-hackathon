package diagnostics

import (
	"context"
	"fmt"
	"net"
	"time"
)

func TestDNS(host string) {
	start := time.Now()

	_, err := net.DefaultResolver.LookupHost(
		context.Background(),
		host,
	)

	duration := time.Since(start)

	fmt.Println("\ndns")
	fmt.Println("  host:", host)
	fmt.Println("  response time:", duration)
	fmt.Println("  success:", err == nil)

	if err != nil {
		fmt.Println("  error:", err)
	}
}
