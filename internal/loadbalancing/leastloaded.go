package loadbalancing

import (
	"fmt"
	"sync"

	"github.com/dalibormesaric/rplb/internal/backend"
)

type leastLoaded struct {
	state *leastLoadedState
}

type leastLoadedState struct {
	mu                    sync.Mutex
	roundRobinForPoolLoad map[string]int
	loadForBackend        map[string]int // Track number of ongoing requests per backend.
}

var _ (Algorithm) = (*leastLoaded)(nil)

func (algo *leastLoaded) GetNext(_ string, backends []*backend.Backend) (backend *backend.Backend, afterBackendResponse func()) {
	algo.state.mu.Lock()
	defer algo.state.mu.Unlock()

	if len(backends) == 0 {
		return nil, nil
	}

	// go through all backends and find all with least n
	// and do round robin for that n
	backend = algo.getLeastLoaded(backends)

	algo.state.loadForBackend[backend.Name]++

	afterBackendResponse = func(name string) func() {
		return func() {
			algo.state.mu.Lock()
			defer algo.state.mu.Unlock()
			algo.state.loadForBackend[name]--
		}
	}(backend.Name)

	return backend, afterBackendResponse
}

func (algo *leastLoaded) getLeastLoaded(backends []*backend.Backend) (be *backend.Backend) {
	minLoadForBackend := 999

	for _, b := range backends {
		loadForBackend := algo.state.loadForBackend[b.Name]
		if loadForBackend < minLoadForBackend {
			minLoadForBackend = loadForBackend
		}
	}

	roundRobinForPoolLoad, ok := algo.state.roundRobinForPoolLoad[getPoolLoad(backends[0], minLoadForBackend)]
	if !ok {
		roundRobinForPoolLoad, algo.state.roundRobinForPoolLoad[getPoolLoad(backends[0], minLoadForBackend)] = 0, 0
	}

	nForPoolLoad := 0
	var firstBackend *backend.Backend
	for _, b := range backends {
		if minLoadForBackend == algo.state.loadForBackend[b.Name] {
			if firstBackend == nil {
				firstBackend = b
			}
			if nForPoolLoad == roundRobinForPoolLoad {
				be = b
				algo.state.roundRobinForPoolLoad[getPoolLoad(backends[0], minLoadForBackend)]++
				break
			} else {
				nForPoolLoad++
			}
		}
	}

	if be == nil {
		algo.state.roundRobinForPoolLoad[getPoolLoad(backends[0], minLoadForBackend)] = 1
		be = firstBackend
	}

	return be
}

func getPoolLoad(backend *backend.Backend, load int) string {
	return fmt.Sprintf("%s%d", backend.GetPoolName(), load)
}
