package reverseproxy

import (
	"log"
	"math/rand"
	"net"
	"net/http"

	"github.com/dalibormesaric/rplb/internal/backend"
	"github.com/dalibormesaric/rplb/internal/frontend"
)

func ListenAndServe(frontends frontend.Frontends, backends backend.Backends) {
	rpMux := http.NewServeMux()
	rpMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		host, _, _ := net.SplitHostPort(r.Host)
		f := frontends.Get(host)
		if f == nil {
			log.Printf("Unknown frontend host: %s\n", host)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		f.Inc()
		// log.Println(f.BackendName)
		liveBackends := backend.GetAlive(backends[f.BackendName])

		n := len(liveBackends)
		if n > 0 {
			randBackend := rand.Intn(n)

			// rw.Header().Add("proxy-url", liveBackends[randBackend].Url)
			liveBackends[randBackend].Proxy.ServeHTTP(w, r)
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
		}
	})
	http.ListenAndServe(":8080", rpMux)
}
