package reverseproxy

import (
	"log"
	"math/rand"
	"net"
	"net/http"

	"github.com/dalibormesaric/rplb/internal/backend"
	"github.com/dalibormesaric/rplb/internal/frontend"
)

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
	rpMux := http.NewServeMux()
	rpMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		host, _, _ := net.SplitHostPort(r.Host)
		f := frontends.Get(host)
		if f == nil {
			log.Printf("Unknown frontend host (%s)\n", host)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		tf := TrafficFrame{Type: "traffic-fe", Name: host, Hits: f.Inc()}
		messages <- tf
		// log.Println(f.BackendName)
		liveBackends := backend.GetLive(backends[f.BackendName])

		n := len(liveBackends)
		if n > 0 {
			randBackend := rand.Intn(n)

			// rw.Header().Add("proxy-url", liveBackends[randBackend].Url)
			liveBackend := liveBackends[randBackend]
			liveBackend.Proxy.ServeHTTP(w, r)
			tf := TrafficBackendFrame{TrafficFrame: &TrafficFrame{Type: "traffic-be", Name: liveBackend.Name, Hits: liveBackend.Inc()}, FrontendName: host}
			messages <- tf
		} else {
			log.Printf("No live backends for host (%s)\n", host)
			w.WriteHeader(http.StatusServiceUnavailable)
		}
	})
	http.ListenAndServe(":8080", rpMux)
}
