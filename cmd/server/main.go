package main

import (
	"embed"
	"flag"
	"html/template"
	"log"
	"math/rand"
	"net"
	"net/http"
	"time"

	"github.com/dalibormesaric/rplb/internal/backend"
	"github.com/dalibormesaric/rplb/internal/frontend"
	"github.com/dalibormesaric/rplb/internal/server"
)

//go:embed template/*.html
var content embed.FS

//go:embed assets/*.css assets/*.js
var assets embed.FS

var (
	fe    = flag.String("f", "", "frontends")
	be    = flag.String("b", "", "backends")
	since = time.Now()
)

func main() {
	log.Println("Starting RPLB...")
	flag.Parse()

	frontends, err := frontend.CreateFrontends(*fe)
	if err != nil {
		log.Fatalf("Create frontends: %s", err)
	}

	backends, err := backend.CreateBackends(*be)
	if err != nil {
		log.Fatalf("Create backends: %s", err)
	}

	// reverse proxy
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

		if len(liveBackends) > 0 {
			randBackend := rand.Intn(len(liveBackends))

			// rw.Header().Add("proxy-url", liveBackends[randBackend].Url)
			liveBackends[randBackend].Proxy.ServeHTTP(w, r)
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
		}
	})
	go http.ListenAndServe(":8080", rpMux)

	// dashboard
	wsServer := server.New()
	http.HandleFunc("/ws", wsServer.WsHandler)
	go wsServer.Broadcaster()

	server := server.NewServer(backends, wsServer.Messages)
	go server.Monitor()

	http.Handle("/assets/", http.StripPrefix("/", http.FileServer(http.FS(assets))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		a, _ := template.
			New("index").
			Funcs(template.FuncMap{
				"printSince": func() string {
					return since.Format("2006-01-02 15:04:05")
				}}).
			ParseFS(content, "template/*.html")
		a.ExecuteTemplate(w, "monitor.html", backends)
	})

	http.ListenAndServe(":8000", nil)
}
