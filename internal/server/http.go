package server

import (
	"fmt"
	"math"
	"time"

	"github.com/dalibormesaric/rplb/internal/backend"
)

type server struct {
	Backends map[string][]*backend.Backend
	messages chan interface{}
}

func NewServer(backends map[string][]*backend.Backend, messages chan interface{}) *server {
	return &server{
		Backends: backends,
		messages: messages,
	}
}

func (s *server) Monitor() {
	for {
		for _, v := range s.Backends {
			for _, b := range v {
				go func(b *backend.Backend) {
					start := time.Now()
					res, err := b.MonitorClient.Get(b.Url)
					latency := time.Since(start)
					if err != nil {
						b.Alive = false
						latency = 0
						_ = fmt.Sprintf("%s\n\terror: %v\n", b.Url, err)
					} else {
						b.Alive = true
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
					case d >= 110:
						colorCode = 1000
					}
					mf := backend.MonitorFrame{Live: b.Alive, Latency: latency, ColorCode: colorCode}
					lmf := backend.LiveMonitorFrame{Name: b.Name, Alive: b.Alive, Latency: fmt.Sprintf("%v", latency), ColorCode: colorCode}
					s.messages <- lmf
					b.Monitor = last20(append(b.Monitor, mf))
				}(b)
			}
			time.Sleep(1 * time.Second)
		}
	}
}

func last20(m []backend.MonitorFrame) []backend.MonitorFrame {
	l := int(math.Max(0, float64(len(m)-20)))
	return m[l:]
}
