package main

import (
	"embed"
	"flag"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/dalibormesaric/rplb/internal/backend"
	"github.com/dalibormesaric/rplb/internal/frontend"
	"github.com/dalibormesaric/rplb/internal/reverseproxy"
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
	monitor := backends.NewMonitor()
	go monitor.Run()

	go reverseproxy.ListenAndServe(frontends, backends)

	// dashboard
	// move wsServer to dashboard package
	wsServer := server.New(monitor.Messages)
	http.HandleFunc("/ws", wsServer.WsHandler)
	go wsServer.Broadcaster()

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
