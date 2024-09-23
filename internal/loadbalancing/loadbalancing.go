package loadbalancing

import (
	"fmt"
	"log"

	"github.com/dalibormesaric/rplb/internal/backend"
)

const (
	Random      string = "random"
	First       string = "first"
	RoundRobin  string = "roundrobin"
	Sticky      string = "sticky"
	LeastLoaded string = "leastloaded"
)

type Algorithm interface {
	// GetNext returns next available backend according to the algorithm with optional [afterBackendResponse] callback
	GetNext(remoteAddr string, backends []*backend.Backend) (backend *backend.Backend, afterBackendResponse func())
}

// Returns new load balancing algorithm.
func NewAlgorithm(name string) (algo Algorithm, err error) {
	defer func() {
		if err == nil {
			log.Printf("Using algorithm (%s)\n", name)
		}
	}()

	switch name {
	case LeastLoaded:
		return &leastLoaded{
			state: &leastLoadedState{
				roundRobinForPoolLoad: make(map[string]int),
				loadForBackend:        make(map[string]int),
			},
		}, nil
	case Sticky:
		return &sticky{
			state: &stickyState{
				backendHostForPoolClientIp: make(map[string]string),
				nForPool:                   make(map[string]int),
			},
		}, nil
	case RoundRobin:
		return &roundRobin{
			state: &roundRobinState{
				nForPool: make(map[string]int),
			},
		}, nil
	case First:
		return &first{}, nil
	case Random:
		return &random{}, nil
	default:
		return nil, fmt.Errorf("unknown algorithm type (%s)", name)
	}
}
