package main

import (
	"flag"
	"math"
	"time"

	"bitbucket.org/bertimus9/systemstat"
)

const sleep = time.Duration(10) * time.Millisecond

func doSomething() {
	for i := 1; i < 10000000; i++ {
		x := float64(0)
		x += math.Sqrt(0)
	}
}

var (
	millicores  = flag.Int("millicores", 0, "millicores number")
	durationSec = flag.Int("duration-sec", 0, "duration time in seconds")
)

func main() {
	flag.Parse()
	// convert millicores to percentage
	millicoresPct := float64(*millicores) / float64(10)
	duration := time.Duration(*durationSec) * time.Second
	start := time.Now()
	first := systemstat.GetProcCPUSample()
	for time.Since(start) < duration {
		cpu := systemstat.GetProcCPUAverage(first, systemstat.GetProcCPUSample(), systemstat.GetUptime().Uptime)
		if cpu.TotalPct < millicoresPct {
			doSomething()
		} else {
			time.Sleep(sleep)
		}
	}
}
