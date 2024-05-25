package backend

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

type backend struct {
	Url           string
	proxy         *httputil.ReverseProxy
	alive         bool
	monitorClient http.Client
	Monitor       []monitorFrame
}

type monitorFrame struct {
	Live    bool
	Latency time.Duration
}

const (
	monitorTimeout = 2 * time.Second
)

func CreateBackends() *map[string][]*backend {
	return &map[string][]*backend{}
}

func CreateBackend(urlString string) (*backend, error) {
	urlParsed, err := url.Parse(urlString)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(urlParsed)

	return &backend{
		urlString,
		proxy,
		false,
		http.Client{Timeout: monitorTimeout},
		[]monitorFrame{},
	}, nil
}

func GetUrlsForNames(nameUrlPairs string) (map[string][]string, error) {
	split := strings.Split(nameUrlPairs, ",")
	if len(split)%2 != 0 {
		return nil, fmt.Errorf("unable to split nameUrlPairs")
	}

	result := make(map[string][]string)
	for i, v := range split {
		if v == "" {
			return nil, fmt.Errorf("nameUrlPair at index %d must have a value", i)
		}

		if (i+1)%2 == 0 {
			k := split[i-1]
			_, ok := result[k]
			if !ok {
				result[k] = []string{split[i]}
			} else {
				result[k] = append(result[k], split[i])
			}
		}
	}

	return result, nil
}
