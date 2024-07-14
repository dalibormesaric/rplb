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
	URL     *url.URL
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
	be := make(Backends)

	if strings.TrimSpace(nameUrlPairs) == "" {
		log.Println("No backends configured")
		return be, nil
	}

	split := strings.Split(nameUrlPairs, ",")
	if len(split)%2 != 0 {
		return nil, fmt.Errorf("backends must be a comma-separated list containing even number of items")
	}

	for i, v := range split {
		if v == "" {
			return nil, fmt.Errorf("nameUrlPair at index %d must have a value", i)
		}

		if (i+1)%2 == 0 {
			k := split[i-1]

			backendUrl := split[i]
			b, err := createBackend(k, backendUrl)
			if err != nil {
				return nil, err
			}

			_, ok := be[k]
			if !ok {
				be[k] = []*Backend{b}
			} else {
				for _, existingBackend := range be[k] {
					if existingBackend.URL.Host == b.URL.Host {
						return nil, fmt.Errorf("url (%s) already exist in backend (%s)", backendUrl, k)
					}
				}
				be[k] = append(be[k], b)
			}
			log.Printf("Added backend url (%s) for (%s)\n", backendUrl, k)
		}
	}

	return be, nil
}

func createBackend(key, urlString string) (*Backend, error) {
	urlParsed, err := url.Parse(urlString)
	if err != nil {
		return nil, err
	}
	if urlParsed.Host == "" {
		return nil, fmt.Errorf("empty host for url (%s) in backend (%s)", urlString, key)
	}

	proxy := httputil.NewSingleHostReverseProxy(urlParsed)

	return &Backend{
		Name:    string(strip([]byte(fmt.Sprintf("%s%s", key, urlString)))),
		URL:     urlParsed,
		Proxy:   proxy,
		Live:    false,
		Monitor: []MonitorFrame{},
		Hits:    0,
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
