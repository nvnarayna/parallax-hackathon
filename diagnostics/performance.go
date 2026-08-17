package diagnostics

import (
	"runtime"
)

type Performance struct {
	HeapAlloc  uint64
	HeapInUse  uint64
	SysMemory  uint64
	NumGC      uint32
	Goroutines int
}

func GetPerformance() Performance {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return Performance{
		HeapAlloc:  m.HeapAlloc,
		HeapInUse:  m.HeapInuse,
		SysMemory:  m.Sys,
		NumGC:      m.NumGC,
		Goroutines: runtime.NumGoroutine(),
	}
}
