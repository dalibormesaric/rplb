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
	loadForBackend        map[string]int
	// track number of ongoing requests per backend
}

var _ (Algorithm) = (*leastLoaded)(nil)

func (algo *leastLoaded) GetNext(_ string, backends []*backend.Backend) (backend *backend.Backend, afterBackendResponse func()) {
	algo.state.mu.Lock()
	defer algo.state.mu.Unlock()

	if len(backends) == 0 {
		return nil, nil
	}

	// ensure all backends are in state with initial load n = 0
	algo.ensureLoadForBackendInState(backends)

	// go through all live backends and find all with least n
	// and do round robin for that n
	backend = algo.getLeastLoad(backends)

	algo.state.loadForBackend[backend.Name]++

	afterBackendResponse = func(name string) func() {
		return func() {
			algo.state.mu.Lock()
			defer algo.state.mu.Unlock()
			algo.state.loadForBackend[name]--
		}
	}(backend.Name)

	return backend, afterBackendResponse
	// find backend with least load (number of requests)
	// increase number of requests for backend
	// call proxy
	// decrease number of requests for backend
}

func (algo *leastLoaded) ensureLoadForBackendInState(backends []*backend.Backend) {
	// TODO: to avoid calling this on every request, pass backends when calling NewAlgorithm
	for _, b := range backends {
		_, ok := algo.state.loadForBackend[b.Name]
		if !ok {
			algo.state.loadForBackend[b.Name] = 0
		}
	}
}

func (algo *leastLoaded) getLeastLoad(backends []*backend.Backend) (ba *backend.Backend) {
	minLoadForBackend := 999

	for _, b := range backends {
		loadForBackend := algo.state.loadForBackend[b.Name]
		if loadForBackend < minLoadForBackend {
			minLoadForBackend = loadForBackend
		}
	}

	roundRobinForLoad, ok := algo.state.roundRobinForPoolLoad[getPoolLoad(backends[0], minLoadForBackend)]
	if !ok {
		roundRobinForLoad, algo.state.roundRobinForPoolLoad[getPoolLoad(backends[0], minLoadForBackend)] = 0, 0
	}

	iForLoad := 0
	var firstBackend *backend.Backend
	for _, b := range backends {
		if minLoadForBackend == algo.state.loadForBackend[b.Name] {
			if firstBackend == nil {
				firstBackend = b
			}
			if iForLoad == roundRobinForLoad {
				ba = b
				algo.state.roundRobinForPoolLoad[getPoolLoad(backends[0], minLoadForBackend)]++
				break
			} else {
				iForLoad++
			}
		}
	}

	if ba == nil {
		algo.state.roundRobinForPoolLoad[getPoolLoad(backends[0], minLoadForBackend)] = 1
		ba = firstBackend
	}

	return ba
}

func getPoolLoad(backend *backend.Backend, load int) string {
	return fmt.Sprintf("%s%d", backend.GetPoolName(), load)
}
