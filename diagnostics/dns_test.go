package diagnostics

import (
	"testing"
)

func TestDNS(t *testing.T) {
	duration, err := CheckDNS("google.com")

	if err != nil {
		t.Fatal(err)
	}

	t.Log("dns response time:", duration)
}
