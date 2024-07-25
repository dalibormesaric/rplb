package frontend

import (
	"fmt"
	"log"
	"strings"
	"sync/atomic"
)

type Host string

type Frontends map[Host]*Frontend

type Frontend struct {
	BackendName string
	hits        atomic.Uint64
}

// CreateFrontends parses url name pairs and returns created Frontends.
func CreateFrontends(urlNamePair string) (Frontends, error) {
	frontends := make(Frontends)

	if strings.TrimSpace(urlNamePair) == "" {
		log.Println("No frontends configured")
		return frontends, nil
	}

	split := strings.Split(urlNamePair, ",")
	if len(split)%2 != 0 {
		return nil, fmt.Errorf("frontends must be a comma-separated list containing even number of items")
	}

	for i, v := range split {
		if v == "" {
			return nil, fmt.Errorf("urlNamePair at index %d must have a value", i)
		}

		if (i+1)%2 == 0 {
			host := Host(split[i-1])

			_, ok := frontends[host]
			if ok {
				return nil, fmt.Errorf("frontend host has to be unique")
			}

			backendName := split[i]
			frontends[host] = &Frontend{
				BackendName: backendName,
			}
			log.Printf("Added frontend host (%s) for (%s)\n", host, backendName)
		}
	}

	return frontends, nil
}

// Get returns Frontend for a host.
func (f Frontends) Get(host string) *Frontend {
	return f[Host(host)]
}

// GetHits returns number of hits for Frontend.
// Concurrency-safe.
func (f *Frontend) GetHits() uint64 {
	return f.hits.Load()
}

// IncHits increases and returns number of hits for Frontend.
// Concurrency-safe.
func (f *Frontend) IncHits() uint64 {
	return f.hits.Add(1)
}
