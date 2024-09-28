package backend

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// The key type is unexported to prevent collisions with context keys defined in
// other packages.
type key string

const (
	// RPLBBackendStatusCode is  the context key to pass error status code from backend to reverse proxy and trigger a retry.
	RPLBBackendStatusCode key = "RPLB-Backend-StatusCode"
)

type BackendPool map[string][]*Backend

type Backend struct {
	mu            sync.RWMutex
	Name          string
	URL           *url.URL
	Proxy         *httputil.ReverseProxy
	live          bool
	monitorFrames []MonitorFrame
	Hits          atomic.Uint64
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

// NewBackendPool parses comma-separated list of Name and URL pairs and returns BackendPool.
func NewBackendPool(nameUrlPairs string) (BackendPool, error) {
	bp := make(BackendPool)

	if strings.TrimSpace(nameUrlPairs) == "" {
		log.Println("No backends configured")
		return bp, nil
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

			_, ok := bp[k]
			if !ok {
				bp[k] = []*Backend{b}
			} else {
				for _, existingBackend := range bp[k] {
					if existingBackend.URL.Host == b.URL.Host {
						return nil, fmt.Errorf("url (%s) already exist in backend pool (%s)", backendUrl, k)
					}
				}
				bp[k] = append(bp[k], b)
			}
			log.Printf("Added backend url (%s) for (%s)\n", backendUrl, k)
		}
	}

	return bp, nil
}

func createBackend(poolName, urlString string) (*Backend, error) {
	urlParsed, err := url.Parse(urlString)
	if err != nil {
		return nil, err
	}
	if urlParsed.Host == "" {
		return nil, fmt.Errorf("empty host for url (%s) in backend pool (%s)", urlString, poolName)
	}

	proxy := httputil.NewSingleHostReverseProxy(urlParsed)
	proxy.ModifyResponse = func(r *http.Response) error {
		if r.StatusCode >= http.StatusInternalServerError {
			return fmt.Errorf("%d", r.StatusCode)
		}
		return nil
	}
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		ctx := context.WithValue(r.Context(), RPLBBackendStatusCode, http.StatusBadGateway)
		req := r.WithContext(ctx)
		*r = *req
	}

	return &Backend{
		Name:          getName(poolName, urlParsed),
		URL:           urlParsed,
		Proxy:         proxy,
		live:          false,
		monitorFrames: []MonitorFrame{},
		Hits:          atomic.Uint64{},
	}, nil
}

// stripString is a string wrapper around strip.
func stripString(s string) string {
	return string(strip([]byte(s)))
}

// strip returns bytes a-z, A-Z and 0-9 and ignores the rest.
func strip(bytes []byte) []byte {
	n := 0
	for _, b := range bytes {
		if ('a' <= b && b <= 'z') ||
			('A' <= b && b <= 'Z') ||
			('0' <= b && b <= '9') {
			bytes[n] = b
			n++
		}
	}
	return bytes[:n]
}

func GetLive(backends []*Backend) (liveBackends []*Backend) {
	for _, b := range backends {
		if b.GetLive() {
			liveBackends = append(liveBackends, b)
		}
	}
	return
}

func getName(poolName string, url *url.URL) string {
	return stripString(fmt.Sprintf("%s%s", poolName, url.String()))
}

func (b *Backend) GetPoolName() string {
	return strings.TrimSuffix(b.Name, stripString(b.URL.String()))
}

// GetHits returns number of hits for Backend.
// Concurrency-safe.
func (b *Backend) GetHits() uint64 {
	return b.Hits.Load()
}

// IncHits increases and returns number of hits for Backend.
// Concurrency-safe.
func (b *Backend) IncHits() uint64 {
	return b.Hits.Add(1)
}

// SetLive sets if Backend is live.
// Concurrency-safe.
func (b *Backend) SetLive(live bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.live = live
}

// GetLive returns if Backend is live.
// Concurrency-safe.
func (b *Backend) GetLive() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.live
}

// SetMonitorFrames sets monitor frames.
// Concurrency-safe.
func (b *Backend) SetMonitorFrames(monitorFrames []MonitorFrame) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.monitorFrames = monitorFrames
}

// GetMonitorFrames returns monitor frames.
// Concurrency-safe.
func (b *Backend) GetMonitorFrames() []MonitorFrame {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.monitorFrames
}
