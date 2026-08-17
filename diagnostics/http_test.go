package diagnostics

import (
	"testing"
)

func TestHTTP(t *testing.T) {
	duration, status, err := CheckHTTP("https://google.com")

	if err != nil {
		t.Fatal(err)
	}

	t.Log("response time:", duration)
	t.Log("status:", status)

	if status < 200 || status >= 400 {
		t.Fatalf("unexpected HTTP status: %d", status)
	}
}
