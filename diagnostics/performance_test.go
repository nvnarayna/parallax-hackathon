package diagnostics

import (
	"runtime"
	"testing"
)

func TestGetPerformance(t *testing.T) {
	performance := GetPerformance()

	if performance.HeapAlloc == 0 {
		t.Error("heap allocation is zero")
	}

	if performance.HeapInUse == 0 {
		t.Error("heap in use is zero")
	}

	if performance.SysMemory == 0 {
		t.Error("runtime memory is zero")
	}

	if performance.Goroutines < 1 {
		t.Error("invalid goroutine count")
	}

	t.Logf("heap allocated: %.2f KB", float64(performance.HeapAlloc)/1024)
	t.Logf("heap in use: %.2f KB", float64(performance.HeapInUse)/1024)
	t.Logf("runtime memory: %.2f KB", float64(performance.SysMemory)/1024)
	t.Logf("gc cycles: %d", performance.NumGC)
	t.Logf("goroutines: %d", performance.Goroutines)

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	if performance.HeapAlloc > m.HeapAlloc {
		t.Error("performance snapshot is inconsistent")
	}
}