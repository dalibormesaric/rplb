package loadbalancing

import (
	"sync"

	"github.com/dalibormesaric/rplb/internal/backend"
)

type roundRobin struct {
	state *roundRobinState
}

type roundRobinState struct {
	mu sync.Mutex
	n  int
}

var _ Algorithm = (*roundRobin)(nil)

func (algo *roundRobin) Get(_ string, backends []*backend.Backend) *backend.Backend {
	algo.state.mu.Lock()
	defer algo.state.mu.Unlock()

	n := len(backends)
	if n == 0 {
		return nil
	}

	if algo.state.n >= n {
		algo.state.n = 0
	}
	b := backends[algo.state.n]
	algo.state.n++
	return b
}
