package diagnostics

import (
	"net/http"
	"time"
)

func CheckHTTP(url string) (time.Duration, int, error) {
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	start := time.Now()

	resp, err := client.Get(url)

	duration := time.Since(start)

	if err != nil {
		return duration, 0, err
	}

	defer resp.Body.Close()

	return duration, resp.StatusCode, nil
}
