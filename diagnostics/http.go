package diagnostics

import (
	"fmt"
	"net/http"
	"time"
)

func TestHTTP(url string) {
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	start := time.Now()

	resp, err := client.Get(url)

	duration := time.Since(start)

	fmt.Println("\nhttp")
	fmt.Println("  url:", url)
	fmt.Println("  response time:", duration)

	if err != nil {
		fmt.Println("  success: false")
		fmt.Println("  error:", err)
		return
	}

	defer resp.Body.Close()

	fmt.Println("  status:", resp.StatusCode)
	fmt.Println("  success:", resp.StatusCode >= 200 && resp.StatusCode < 400)
}
