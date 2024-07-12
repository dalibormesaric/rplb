package reverseproxy

import (
	"log"
	"math/rand"
	"net"
	"net/http"

	"github.com/dalibormesaric/rplb/internal/backend"
	"github.com/dalibormesaric/rplb/internal/frontend"
)

type reverseProxy struct {
	frontends frontend.Frontends
	backends  backend.Backends
	messages  chan interface{}
}

type TrafficFrame struct {
	Type string
	Name string
	Hits int64
}

type TrafficBackendFrame struct {
	*TrafficFrame
	FrontendName string
}

func ListenAndServe(frontends frontend.Frontends, backends backend.Backends, messages chan interface{}) {
	rp := &reverseProxy{
		frontends: frontends,
		backends:  backends,
		messages:  messages,
	}
	rpMux := http.NewServeMux()
	rpMux.HandleFunc("/", rp.ServeHTTP)
	http.ListenAndServe(":8080", rpMux)
}

func (rp *reverseProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	host, _, _ := net.SplitHostPort(r.Host)
	f := rp.frontends.Get(host)
	if f == nil {
		log.Printf("Unknown frontend host (%s)\n", host)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	tf := TrafficFrame{Type: "traffic-fe", Name: host, Hits: f.Inc()}
	rp.messages <- tf
	// log.Println(f.BackendName)
	liveBackends := backend.GetLive(rp.backends[f.BackendName])

	n := len(liveBackends)
	if n > 0 {
		randBackend := rand.Intn(n)

		// rw.Header().Add("proxy-url", liveBackends[randBackend].Url)
		liveBackend := liveBackends[randBackend]
		liveBackend.Proxy.ServeHTTP(w, r)
		tf := TrafficBackendFrame{TrafficFrame: &TrafficFrame{Type: "traffic-be", Name: liveBackend.Name, Hits: liveBackend.Inc()}, FrontendName: host}
		rp.messages <- tf
	} else {
		log.Printf("No live backends for host (%s)\n", host)
		w.WriteHeader(http.StatusServiceUnavailable)
	}
}
