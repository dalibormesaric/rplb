package dashboard

import (
	"embed"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/dalibormesaric/rplb/internal/backend"
	"github.com/dalibormesaric/rplb/internal/frontend"
)

type BaseModel struct {
	SelectedMenu string
	Version      string
}

type MonitorModel struct {
	BaseModel
	Backends backend.BackendPool
}

type TrafficModel struct {
	BaseModel
	Frontends frontend.Frontends
	Backends  backend.BackendPool
}

//go:embed templates/*.html
var content embed.FS

//go:embed assets/*.css assets/*.js assets/*.ico
var assets embed.FS

var (
	since = time.Now()

	FrontendHits = promauto.NewCounter(prometheus.CounterOpts{
		Name: "rplb_frontend_hits",
		Help: "The total number of frontend hits",
	})
	BackendRetries = promauto.NewCounter(prometheus.CounterOpts{
		Name: "rplb_backend_retries",
		Help: "The total number of backend retries",
	})
	BackendHits = promauto.NewCounter(prometheus.CounterOpts{
		Name: "rplb_backend_hits",
		Help: "The total number of backend hits",
	})
)

func ListenAndServe(frontends frontend.Frontends, bp backend.BackendPool, messages chan interface{}, version string) {
	// TODO: wsServer should produce chan messages?
	wsServer := NewWsServer(messages)
	http.HandleFunc("/ws", wsServer.WsHandler)
	go wsServer.Broadcaster()

	http.Handle("/assets/", http.StripPrefix("/", http.FileServer(http.FS(assets))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RequestURI() == "/favicon.ico" {
			http.ServeFileFS(w, r, assets, "assets/favicon.ico")
			return
		}

		http.Redirect(w, r, "/monitor", http.StatusPermanentRedirect)
	})

	http.HandleFunc("/monitor", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		a, _ := template.
			New("index").
			Funcs(getFuncMap()).
			ParseFS(content, "templates/index.html", "templates/monitor.html")
		a.ExecuteTemplate(w, "monitor.html", MonitorModel{BaseModel: BaseModel{SelectedMenu: "monitor", Version: version}, Backends: bp})
	})

	http.HandleFunc("/traffic", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		a, _ := template.
			New("index").
			Funcs(getFuncMap()).
			ParseFS(content, "templates/index.html", "templates/traffic.html")
		a.ExecuteTemplate(w, "traffic.html", TrafficModel{BaseModel: BaseModel{SelectedMenu: "traffic", Version: version}, Frontends: frontends, Backends: bp})
	})

	http.Handle("/metrics", promhttp.Handler())

	log.Printf("Dashboard listening on :8000\n")
	http.ListenAndServe(":8000", nil)
}

func getFuncMap() template.FuncMap {
	return template.FuncMap{
		"printSince": func() string {
			return since.Format("2006-01-02 15:04:05 MST")
		},
		"menuSelected": func(selectedMenu, menuItem string) string {
			if selectedMenu == menuItem {
				return " menu-selected"
			}
			return ""
		}}
}
