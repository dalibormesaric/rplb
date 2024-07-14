package backend

import (
	"fmt"
	"log"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

type Backends map[string][]*Backend

type Backend struct {
	Name    string
	Url     string
	Proxy   *httputil.ReverseProxy
	Live    bool
	Monitor []MonitorFrame
	Hits    int64
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

func CreateBackends(nameUrlPairs string) (Backends, error) {
	split := strings.Split(nameUrlPairs, ",")
	if len(split)%2 != 0 {
		return nil, fmt.Errorf("backends must be a comma-separated list containing even number of items")
	}

	backends := make(Backends)

	for i, v := range split {
		if v == "" {
			return nil, fmt.Errorf("nameUrlPair at index %d must have a value", i)
		}

		if (i+1)%2 == 0 {
			k := split[i-1]
			_, ok := backends[k]
			backendUrl := split[i]
			b, _ := createBackend(k, backendUrl)
			if !ok {
				backends[k] = []*Backend{b}
			} else {
				for _, existingBackend := range backends[k] {
					if existingBackend.Url == backendUrl {
						return nil, fmt.Errorf("url (%s) already exist in backend (%s)", backendUrl, k)
					}
				}
				backends[k] = append(backends[k], b)
			}
			log.Printf("Added backend url (%s) for (%s)\n", backendUrl, k)
		}
	}

	return backends, nil
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
