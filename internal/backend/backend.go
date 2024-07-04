package backend

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

type Backends map[string][]*Backend

type Backend struct {
	Name          string
	Url           string
	Proxy         *httputil.ReverseProxy
	Live          bool
	MonitorClient http.Client
	Monitor       []MonitorFrame
	Hits          int64
}

type MonitorFrame struct {
	Live      bool
	Latency   time.Duration
	ColorCode int64
}

type LiveMonitorFrame struct {
	Type      string
	Name      string
	Live      bool
	Latency   string
	ColorCode int64
}

const (
	monitorTimeout = 2 * time.Second
)

func CreateBackends(nameUrlPairs string) (Backends, error) {
	split := strings.Split(nameUrlPairs, ",")
	if len(split)%2 != 0 {
		return nil, fmt.Errorf("backends must be a comma-separated list containing even number of items")
	}

	be := make(Backends)

	for i, v := range split {
		if v == "" {
			return nil, fmt.Errorf("nameUrlPair at index %d must have a value", i)
		}

		if (i+1)%2 == 0 {
			k := split[i-1]
			_, ok := be[k]
			b, _ := createBackend(k, split[i])
			if !ok {
				be[k] = []*Backend{b}
			} else {
				be[k] = append(be[k], b)
			}
		}
	}

	return be, nil
}

func createBackend(key, urlString string) (*Backend, error) {
	urlParsed, err := url.Parse(urlString)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(urlParsed)

	return &Backend{
		string(strip([]byte(fmt.Sprintf("%s%s", key, urlString)))),
		urlString,
		proxy,
		false,
		http.Client{Timeout: monitorTimeout},
		[]MonitorFrame{},
		0,
	}, nil
}

func strip(s []byte) []byte {
	n := 0
	for _, b := range s {
		if ('a' <= b && b <= 'z') ||
			('A' <= b && b <= 'Z') ||
			('0' <= b && b <= '9') ||
			b == ' ' {
			s[n] = b
			n++
		}
	}
	return s[:n]
}

func GetLive(backends []*Backend) (liveBackends []*Backend) {
	for _, b := range backends {
		if b.Live {
			liveBackends = append(liveBackends, b)
		}
	}
	return
}

func (b *Backend) Inc() int64 {
	b.Hits++
	return b.Hits
}
