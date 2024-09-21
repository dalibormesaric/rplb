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
	mu                  sync.Mutex
	clientIpBackendHost map[string]string
	n                   int
}

var _ Algorithm = (*sticky)(nil)

func (algo *sticky) GetNext(remoteAddr string, backends []*backend.Backend) (backend *backend.Backend, _ func()) {
	algo.state.mu.Lock()
	defer algo.state.mu.Unlock()

	n := len(backends)
	if n == 0 {
		return nil, nil
	}

	host, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		return nil, nil
	}
	clientIp := host

	backendHost, ok := algo.state.clientIpBackendHost[clientIp]
	if ok {
		for _, b := range backends {
			if backendHost == b.URL.Host {
				return b, nil
			}
		}
	} else {
		defer func() { algo.state.n++ }()
	}

	if algo.state.n >= n {
		algo.state.n = 0
	}
	b := backends[algo.state.n]
	algo.state.clientIpBackendHost[clientIp] = b.URL.Host
	return b, nil
}
