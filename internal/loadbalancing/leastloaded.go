package loadbalancing

import (
	"sync"

	"github.com/dalibormesaric/rplb/internal/backend"
)

type leastLoaded struct {
	state *leastLoadedState
}

type leastLoadedState struct {
	mu                sync.Mutex
	roundRobinForLoad map[int]int
	loadForBackend    map[string]int
	// track number of ongoing requests per backend
}

var _ (Algorithm) = (*leastLoaded)(nil)

func (algo *leastLoaded) Get(_ string, backends []*backend.Backend) *backend.Backend {
	// find backend with least load (number of requests)
	// increase number of requests for backend
	// call proxy
	// decrease number of requests for backend
	return nil
}

func (algo *leastLoaded) ensureBackendInState(backends []*backend.Backend) {
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

	roundRobinForLoad, ok := algo.state.roundRobinForLoad[minLoadForBackend]
	if !ok {
		roundRobinForLoad, algo.state.roundRobinForLoad[minLoadForBackend] = 0, 0
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
				algo.state.roundRobinForLoad[minLoadForBackend]++
				break
			} else {
				iForLoad++
			}
		}
	}

	if ba == nil {
		algo.state.roundRobinForLoad[minLoadForBackend] = 0
		ba = firstBackend
	}

	return ba
}

func (algo *leastLoaded) Get2(_ string, backends []*backend.Backend) (b *backend.Backend, afterBackendResponse func()) {
	algo.state.mu.Lock()
	defer algo.state.mu.Unlock()

	// ensure all backends are in state with default n = 0
	algo.ensureBackendInState(backends)

	// go through all live backends and find all with least n
	// and do round robin for that n

	b = algo.getLeastLoad(backends)
	// roundRobinForRequests := make(map[int]int)

	// fmt.Printf("%s %d\n", b.Name, algo.state.loadForBackend[b.Name])
	algo.state.loadForBackend[b.Name]++
	// fmt.Printf("%s %d\n", b.Name, algo.state.loadForBackend[b.Name])
	// find backend with least load (number of requests)
	// t := make(map[string]int)

	// max := 999
	// for k, v := range t {
	// 	if v < max {
	// 		max = v
	// 	}
	// }

	// backend.Name

	afterBackendResponse = func(name string) func() {
		return func() {
			algo.state.mu.Lock()
			defer algo.state.mu.Unlock()
			algo.state.loadForBackend[name]--
			// fmt.Printf("%s %d\n", name, algo.state.loadForBackend[name])
		}
	}(b.Name)

	return b, afterBackendResponse
}
