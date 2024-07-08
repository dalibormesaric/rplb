package backend

import (
	"fmt"
	"math"
	"time"
)

type Monitor struct {
	Backends Backends
	Messages chan interface{}
}

func (b Backends) NewMonitor() *Monitor {
	return &Monitor{
		Backends: b,
		Messages: make(chan interface{}),
	}
}

func (m *Monitor) Run() {
	for {
		for _, v := range m.Backends {
			for _, b := range v {
				go func(b *Backend) {
					start := time.Now()
					res, err := b.MonitorClient.Get(b.Url)
					latency := time.Since(start)
					if err != nil {
						b.Live = false
						latency = 0
						_ = fmt.Sprintf("%s\n\terror: %v\n", b.Url, err)
					} else {
						b.Live = true
						_ = fmt.Sprintf("%s\n\tstatus code: %v\n\tlatency: %v\n", b.Url, res.StatusCode, latency)
					}
					duration, _ := time.ParseDuration(fmt.Sprintf("%v", latency))
					colorCode := int64(0)
					switch d := duration.Milliseconds(); {
					case d > 0 && d < 5:
						colorCode = 5
					case d >= 5 && d < 10:
						colorCode = 10
					case d >= 10 && d < 110:
						colorCode = (d / 10) * 10
					case d >= 110 && d < 1100:
						colorCode = (d / 100) * 100
					case d >= 1100:
						colorCode = 10000
					}
					mf := MonitorFrame{Live: b.Live, Latency: latency, ColorCode: colorCode}
					lmf := LiveMonitorFrame{Type: "monitor", Name: b.Name, Live: b.Live, Latency: fmt.Sprintf("%v", latency), ColorCode: colorCode}
					m.Messages <- lmf
					b.Monitor = last20(append(b.Monitor, mf))
				}(b)
			}
			time.Sleep(1 * time.Second)
		}
	}
}

func last20(m []MonitorFrame) []MonitorFrame {
	l := int(math.Max(0, float64(len(m)-20)))
	return m[l:]
}