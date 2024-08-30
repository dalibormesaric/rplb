package loadbalancing

import (
	"fmt"
	"log"

	"github.com/dalibormesaric/rplb/internal/backend"
)

const (
	Random     string = "random"
	First      string = "first"
	RoundRobin string = "roundrobin"
	Sticky     string = "sticky"
)

type Algorithm interface {
	// Get returns next available backend according to the algorithm
	Get(remoteAddr string, backends []*backend.Backend) *backend.Backend
}

func NewAlgorithm(name string) (algo Algorithm, err error) {
	defer func() {
		if err == nil {
			log.Printf("Using algorithm (%s)\n", name)
		}
	}()

	switch name {
	case Sticky:
		return &sticky{
			state: &stickyState{
				clientIpBackendHost: make(map[string]string),
			},
		}, nil
	case RoundRobin:
		return &roundRobin{
			state: &roundRobinState{},
		}, nil
	case First:
		return &first{}, nil
	case Random:
		return &random{}, nil
	default:
		return nil, fmt.Errorf("unknown algorithm type (%s)\n", name)
	}
}
