package main

import (
	"embed"
	"flag"
	"html/template"
	"log"
	"net/http"

	"github.com/dalibormesaric/rplb/internal/backend"
	"github.com/dalibormesaric/rplb/internal/dashboard"
	"github.com/dalibormesaric/rplb/internal/frontend"
	"github.com/dalibormesaric/rplb/internal/reverseproxy"
	"github.com/dalibormesaric/rplb/internal/server"
)

//go:embed template/*.html
var content embed.FS

//go:embed assets/*.css assets/*.js
var assets embed.FS

var (
	fe = flag.String("f", "", "frontends")
	be = flag.String("b", "", "backends")
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

	go reverseproxy.ListenAndServe(frontends, backends, monitor.Messages)

	// dashboard
	// move wsServer to dashboard package
	// TODO: wsServer should produce chan messages?
	wsServer := server.New(monitor.Messages)
	http.HandleFunc("/ws", wsServer.WsHandler)
	go wsServer.Broadcaster()

	http.Handle("/assets/", http.StripPrefix("/", http.FileServer(http.FS(assets))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/monitor", http.StatusPermanentRedirect)
	})

	http.HandleFunc("/monitor", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		a, _ := template.
			New("index").
			Funcs(dashboard.GetFuncMap()).
			ParseFS(content, "template/index.html", "template/monitor.html")
		a.ExecuteTemplate(w, "monitor.html", dashboard.MonitorModel{BaseModel: &dashboard.BaseModel{SelectedMenu: "monitor"}, Backends: backends})
	})

	http.HandleFunc("/traffic", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		a, _ := template.
			New("index").
			Funcs(dashboard.GetFuncMap()).
			ParseFS(content, "template/index.html", "template/traffic.html")
		a.ExecuteTemplate(w, "traffic.html", dashboard.TrafficModel{BaseModel: &dashboard.BaseModel{SelectedMenu: "traffic"}, Frontends: frontends, Backends: backends})
	})

	http.ListenAndServe(":8000", nil)
}
