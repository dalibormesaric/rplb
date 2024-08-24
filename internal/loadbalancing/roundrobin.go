package loadbalancing

import (
	"net/http"
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

func (algo *roundRobin) Get(r *http.Request, liveBackends []*backend.Backend) *backend.Backend {
	algo.state.mu.Lock()
	defer algo.state.mu.Unlock()

	n := len(liveBackends)
	if n == 0 {
		return nil
	}

	if algo.state.n >= n {
		algo.state.n = 0
	}
	liveBackend := liveBackends[algo.state.n]
	algo.state.n++
	return liveBackend
}
