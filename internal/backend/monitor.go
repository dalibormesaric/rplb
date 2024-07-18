package backend

import (
	"fmt"
	"math"
	"net"
	"time"
)

const (
	monitorInterval    = 1 * time.Second
	healthCheckTimeout = 2 * time.Second
)

func (b Backends) Monitor() chan interface{} {
	messages := make(chan interface{})

	for _, v := range b {
		for _, backend := range v {
			go backend.monitor(messages)
		}
	}

	return messages
}

func (be *Backend) monitor(messages chan interface{}) {
	for {
		latency := healthCheck(be.URL.Host)
		be.SetLive(latency > 0)

		colorCode := getColorCode(latency)

		mf := MonitorFrame{Live: be.GetLive(), Latency: latency, ColorCode: colorCode}
		be.SetMonitorFrames(last20(append(be.GetMonitorFrames(), mf)))

		lmf := LiveMonitorFrame{Type: "monitor", Name: be.Name, Live: be.GetLive(), Latency: fmt.Sprintf("%v", latency), ColorCode: colorCode}
		messages <- lmf

		time.Sleep(monitorInterval)
	}
}

func healthCheck(host string) (latency time.Duration) {
	start := time.Now()
	conn, err := net.DialTimeout("tcp", host, healthCheckTimeout)
	if err != nil {
		return 0
	}
	latency = time.Since(start)
	conn.Close()
	return
}

func getColorCode(latency time.Duration) (colorCode int64) {
	switch l := latency.Microseconds(); {
	case l > 0 && l < 5_000:
		colorCode = 5
	case l >= 5_000 && l < 10_000:
		colorCode = 10
	case l >= 10_000 && l < 110_000:
		colorCode = (l / 10_000) * 10
	case l >= 110_000 && l < 1_100_000:
		colorCode = (l / 100_000) * 100
	case l >= 1_100_000:
		colorCode = 10000
	}
	return
}

func last20(m []MonitorFrame) []MonitorFrame {
	l := int(math.Max(0, float64(len(m)-20)))
	return m[l:]
}
