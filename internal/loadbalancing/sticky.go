package loadbalancing

import (
	"net"
	"sync"

	"github.com/dalibormesaric/rplb/internal/backend"
)

type sticky struct {
	state *stickyState
}

type stickyState struct {
	mu                         sync.Mutex
	backendHostForPoolClientIp map[string]string
	// keeps track of round robin per pool
	nForPool map[string]int
}

var _ Algorithm = (*sticky)(nil)

func (algo *sticky) GetNext(remoteAddr string, backends []*backend.Backend) (backend *backend.Backend, _ func()) {
	algo.state.mu.Lock()
	defer algo.state.mu.Unlock()

	l := len(backends)
	if l == 0 {
		return nil, nil
	}

	clientIp, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		return nil, nil
	}

	poolName := backends[0].GetPoolName()

	backendHost, ok := algo.state.backendHostForPoolClientIp[getPoolClientIp(poolName, clientIp)]
	if ok {
		for _, b := range backends {
			if backendHost == b.URL.Host {
				return b, nil
			}
		}
	}
	n := algo.state.nForPool[poolName]

	// if current round robin target is larger then number of backends
	if n >= l {
		// we start from beginning
		n = 0
	}
	b := backends[n]
	algo.state.nForPool[poolName] = n + 1
	algo.state.backendHostForPoolClientIp[getPoolClientIp(poolName, clientIp)] = b.URL.Host
	return b, nil
}

func getPoolClientIp(pool, clientIp string) string {
	return pool + clientIp
}
