package reverseproxy

import (
	"embed"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/dalibormesaric/rplb/internal/backend"
	"github.com/dalibormesaric/rplb/internal/frontend"
	"github.com/dalibormesaric/rplb/internal/loadbalancing"
)

type reverseProxy struct {
	frontends     frontend.Frontends
	bp            backend.BackendPool
	messages      chan interface{}
	loadbalancing loadbalancing.Algorithm
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

func ListenAndServe(frontends frontend.Frontends, bp backend.BackendPool, loadbalancing loadbalancing.Algorithm, messages chan interface{}) {
	rp := &reverseProxy{
		frontends:     frontends,
		bp:            bp,
		messages:      messages,
		loadbalancing: loadbalancing,
	}
	rpMux := http.NewServeMux()
	rpMux.HandleFunc("/", rp.reverseProxyAndLoadBalance)

	log.Printf("Reverse Proxy listening on :8080\n")
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
		liveBackends := backend.GetLive(rp.bp[f.BackendName])
		liveBackend, afterBackendResponse := rp.loadbalancing.GetNext(r.RemoteAddr, liveBackends)
		if liveBackend == nil {
			log.Printf("No live backends for host (%s)\n", host)
			http.ServeFileFS(w, r, static, "static/503.html")
			break
		}
		// rw.Header().Add("proxy-url", liveBackends[randBackend].Url)
		liveBackend.Proxy.ServeHTTP(w, r)
		if afterBackendResponse != nil {
			afterBackendResponse()
		}
		rplbBackendStatusCode, _ := strconv.Atoi(w.Header().Get(backend.RPLBBackendStatusCode))
		w.Header().Del(backend.RPLBBackendStatusCode)

		if rplbBackendStatusCode < http.StatusInternalServerError {
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
