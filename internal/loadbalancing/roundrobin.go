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
	// keeps track of round robin per pool
	nForPool map[string]int
}

var _ Algorithm = (*roundRobin)(nil)

func (algo *roundRobin) GetNext(_ string, backends []*backend.Backend) (backend *backend.Backend, _ func()) {
	algo.state.mu.Lock()
	defer algo.state.mu.Unlock()

	l := len(backends)
	if l == 0 {
		return nil, nil
	}

	poolName := backends[0].GetPoolName()
	n := algo.state.nForPool[poolName]

	// if current round robin target is larger then number of backends
	if n >= l {
		// we start from beginning
		n = 0
	}
	b := backends[n]
	algo.state.nForPool[poolName] = n + 1
	return b, nil
}
