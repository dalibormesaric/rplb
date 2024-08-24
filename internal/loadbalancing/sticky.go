package loadbalancing

import (
	"net"
	"net/http"
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

func (algo *sticky) Get(r *http.Request, liveBackends []*backend.Backend) *backend.Backend {
	algo.state.mu.Lock()
	defer algo.state.mu.Unlock()

	n := len(liveBackends)
	if n == 0 {
		return nil
	}

	host, _, _ := net.SplitHostPort(r.RemoteAddr)
	clientIp := host

	backendHost, ok := algo.state.clientIpBackendHost[clientIp]
	if ok {
		for _, b := range liveBackends {
			if backendHost == b.URL.Host {
				return b
			}
		}
	} else {
		algo.state.n++
	}

	if algo.state.n >= n {
		algo.state.n = 0
	}
	b := liveBackends[algo.state.n]
	algo.state.clientIpBackendHost[clientIp] = b.URL.Host
	return b
}
