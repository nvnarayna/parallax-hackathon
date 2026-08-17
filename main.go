package main

import (
	"fmt"
	"sync"
	"time"

	"wifi-diagnostic/diagnostics"
)

func main() {
	start := time.Now()
	before := diagnostics.GetPerformance()

	fmt.Println("network health test")
	fmt.Println("--------------------")

	var wg sync.WaitGroup
	wg.Add(4)

	go func() {
		defer wg.Done()
		fmt.Println("dns test")
		diagnostics.TestDNS("google.com")
	}()

	go func() {
		defer wg.Done()
		fmt.Println("http test")
		diagnostics.TestHTTP("https://google.com")
	}()

	go func() {
		defer wg.Done()
		fmt.Println("ping latency")
		diagnostics.TestCheckPing("8.8.8.8:443")
	}()

	go func() {
		defer wg.Done()
		fmt.Println("packet loss")
		diagnostics.TestPacketLoss("8.8.8.8")
	}()

	wg.Wait()

	after := diagnostics.GetPerformance()

	fmt.Println("\nperformance")
	fmt.Println("-----------")
	fmt.Println("execution time:", time.Since(start))
	fmt.Printf("heap allocated: %.2f MB\n", float64(after.HeapAlloc-before.HeapAlloc)/1024/1024)
	fmt.Printf("heap in use: %.2f MB\n", float64(after.HeapInUse-before.HeapInUse)/1024/1024)
	fmt.Printf("runtime memory/sys: %.2f MB\n", float64(after.SysMemory)/1024/1024)
	fmt.Println("gc cycles:", after.NumGC-before.NumGC)
	fmt.Println("goroutines:", after.Goroutines)
}
