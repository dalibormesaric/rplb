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
	Alive         bool
	MonitorClient http.Client
	Monitor       []MonitorFrame
}

type MonitorFrame struct {
	Live      bool
	Latency   time.Duration
	ColorCode int64
}

type LiveMonitorFrame struct {
	Name      string
	Alive     bool
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
	}, nil
}

// func GetUrlsForNames(nameUrlPairs string) (map[string][]string, error) {
// 	split := strings.Split(nameUrlPairs, ",")
// 	if len(split)%2 != 0 {
// 		return nil, fmt.Errorf("unable to split nameUrlPairs")
// 	}

// 	result := make(map[string][]string)
// 	for i, v := range split {
// 		if v == "" {
// 			return nil, fmt.Errorf("nameUrlPair at index %d must have a value", i)
// 		}

// 		if (i+1)%2 == 0 {
// 			k := split[i-1]
// 			_, ok := result[k]
// 			if !ok {
// 				result[k] = []string{split[i]}
// 			} else {
// 				result[k] = append(result[k], split[i])
// 			}
// 		}
// 	}

// 	return result, nil
// }

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

func GetAlive(backends []*Backend) (liveBackends []*Backend) {
	for _, b := range backends {
		if b.Alive {
			liveBackends = append(liveBackends, b)
		}
	}
	return
}
