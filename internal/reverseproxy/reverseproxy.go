package reverseproxy

import (
	"embed"
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
	Hits uint64
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
	// log.Println(f.BackendName)
	liveBackends := backend.GetLive(rp.backends[f.BackendName])

	n := len(liveBackends)
	if n == 0 {
		log.Printf("No live backends for host (%s)\n", host)
		http.ServeFileFS(w, r, static, "static/503.html")
		return
	}

	randBackend := rand.Intn(n)

	// rw.Header().Add("proxy-url", liveBackends[randBackend].Url)
	liveBackend := liveBackends[randBackend]
	liveBackend.Proxy.ServeHTTP(w, r)
	if rp.messages != nil {
		tf := TrafficBackendFrame{TrafficFrame: &TrafficFrame{Type: "traffic-be", Name: liveBackend.Name, Hits: liveBackend.IncHits()}, FrontendName: host}
		rp.messages <- tf
	}
}
