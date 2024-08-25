package loadbalancing

import (
	"net/http"

	"github.com/dalibormesaric/rplb/internal/backend"
)

const (
	Random     string = "random"
	RoundRobin string = "roundrobin"
	Sticky     string = "sticky"
)

type Algorithm interface {
	// Get returns next available backend according to the algorithm
	Get(r *http.Request, liveBackends []*backend.Backend) *backend.Backend
}

func NewAlgorithm(name string) Algorithm {
	switch name {
	case Sticky:
		return &sticky{
			state: &stickyState{
				clientIpBackendHost: make(map[string]string),
			},
		}
	case RoundRobin:
		return &roundRobin{
			state: &roundRobinState{},
		}

	case Random:
		fallthrough
	default:
		return &random{}
	}
}
