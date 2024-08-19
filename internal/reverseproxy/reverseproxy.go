package reverseproxy

import (
	"embed"
	"log"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/dalibormesaric/rplb/internal/backend"
	"github.com/dalibormesaric/rplb/internal/frontend"
)

type reverseProxy struct {
	frontends       frontend.Frontends
	backends        backend.Backends
	messages        chan interface{}
	roundRobinState *roundRobinState
	stickyState     *stickyState
}

type TrafficFrame struct {
	Type string
	Name string
	Hits uint64
}

type TrafficBackendFrame struct {
	*TrafficFrame
	FrontendName string
}

func ListenAndServe(frontends frontend.Frontends, backends backend.Backends, messages chan interface{}) {
	rp := &reverseProxy{
		frontends:       frontends,
		backends:        backends,
		messages:        messages,
		roundRobinState: &roundRobinState{},
		stickyState: &stickyState{
			clientIpBackendHost: make(map[string]string),
		},
	}
	rpMux := http.NewServeMux()
	rpMux.HandleFunc("/", rp.reverseProxyAndLoadBalance)
	http.ListenAndServe(":8080", rpMux)
}

//go:embed static/*.html
var static embed.FS

func (rp *reverseProxy) reverseProxyAndLoadBalance(w http.ResponseWriter, r *http.Request) {
	host, _, _ := net.SplitHostPort(r.Host)
	f := rp.frontends.Get(host)
	if f == nil {
		log.Printf("Unknown frontend host (%s)\n", host)
		http.ServeFileFS(w, r, static, "static/404.html")
		return
	}

	if rp.messages != nil {
		tf := TrafficFrame{Type: "traffic-fe", Name: host, Hits: f.IncHits()}
		rp.messages <- tf
	}

	retryTimeout := 500 * time.Millisecond
	retryAmount := 5
	for range retryAmount {
		liveBackend := rp.sticky(r.RemoteAddr, f.BackendName)
		// liveBackend := rp.roundRobin(f.BackendName)
		// liveBackend := rp.random(f.BackendName)
		if liveBackend == nil {
			log.Printf("No live backends for host (%s)\n", host)
			http.ServeFileFS(w, r, static, "static/503.html")
			break
		}
		// rw.Header().Add("proxy-url", liveBackends[randBackend].Url)
		liveBackend.Proxy.ServeHTTP(w, r)
		rplb, _ := strconv.Atoi(w.Header().Get("RPLB-Backend-StatusCode"))
		w.Header().Del("RPLB-Backend-StatusCode")

		if rplb < http.StatusInternalServerError {
			if rp.messages != nil {
				tf := TrafficBackendFrame{TrafficFrame: &TrafficFrame{Type: "traffic-be", Name: liveBackend.Name, Hits: liveBackend.IncHits()}, FrontendName: host}
				rp.messages <- tf
			}
			break
		}

		time.Sleep(retryTimeout)
		retryTimeout *= 2
	}
}

type stickyState struct {
	mu                  sync.Mutex
	clientIpBackendHost map[string]string
}

func (rp *reverseProxy) sticky(remoteAddr, backendUrl string) *backend.Backend {
	rp.stickyState.mu.Lock()
	defer rp.stickyState.mu.Unlock()
	liveBackends := backend.GetLive(rp.backends[backendUrl])

	n := len(liveBackends)
	if n == 0 {
		return nil
	}

	host, _, _ := net.SplitHostPort(remoteAddr)
	clientIp := host

	backendHost, ok := rp.stickyState.clientIpBackendHost[clientIp]
	if ok {
		for _, b := range liveBackends {
			if backendHost == b.URL.Host {
				return b
			}
		}
	}

	b := liveBackends[0]
	rp.stickyState.clientIpBackendHost[clientIp] = b.URL.Host
	return b
}

type roundRobinState struct {
	mu sync.Mutex
	n  int
}

func (rp *reverseProxy) roundRobin(backendName string) *backend.Backend {
	rp.roundRobinState.mu.Lock()
	defer rp.roundRobinState.mu.Unlock()
	liveBackends := backend.GetLive(rp.backends[backendName])

	n := len(liveBackends)
	if n == 0 {
		return nil
	}

	if rp.roundRobinState.n >= n {
		rp.roundRobinState.n = 0
	}
	liveBackend := liveBackends[rp.roundRobinState.n]
	rp.roundRobinState.n++
	return liveBackend
}

func (rp *reverseProxy) random(backendName string) *backend.Backend {
	liveBackends := backend.GetLive(rp.backends[backendName])

	n := len(liveBackends)
	if n == 0 {
		return nil
	}

	randBackend := rand.Intn(n)
	liveBackend := liveBackends[randBackend]
	return liveBackend
}
