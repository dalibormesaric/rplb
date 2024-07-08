package dashboard

import (
	"embed"
	"html/template"
	"net/http"
	"time"

	"github.com/dalibormesaric/rplb/internal/backend"
	"github.com/dalibormesaric/rplb/internal/frontend"
)

//go:embed templates/*.html
var content embed.FS

//go:embed assets/*.css assets/*.js
var assets embed.FS

var (
	since = time.Now()
)

func ListenAndServe(frontends frontend.Frontends, backends backend.Backends, messages chan interface{}) {
	// TODO: wsServer should produce chan messages?
	wsServer := NewWsServer(messages)
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
			Funcs(getFuncMap()).
			ParseFS(content, "templates/index.html", "templates/monitor.html")
		a.ExecuteTemplate(w, "monitor.html", MonitorModel{BaseModel: &BaseModel{SelectedMenu: "monitor"}, Backends: backends})
	})

	http.HandleFunc("/traffic", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		a, _ := template.
			New("index").
			Funcs(getFuncMap()).
			ParseFS(content, "templates/index.html", "templates/traffic.html")
		a.ExecuteTemplate(w, "traffic.html", TrafficModel{BaseModel: &BaseModel{SelectedMenu: "traffic"}, Frontends: frontends, Backends: backends})
	})

	http.ListenAndServe(":8000", nil)
}

func getFuncMap() template.FuncMap {
	return template.FuncMap{
		"printSince": func() string {
			return since.Format("2006-01-02 15:04:05")
		},
		"menuSelected": func(selectedMenu, menuItem string) string {
			if selectedMenu == menuItem {
				return " menu-selected"
			}
			return ""
		}}
}
