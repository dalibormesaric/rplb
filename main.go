package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"math"
	"math/rand"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

type monitorFrame struct {
	Live    bool
	Latency time.Duration
}

type backend struct {
	Url     string
	client  http.Client
	live    bool
	proxy   *httputil.ReverseProxy
	Monitor []monitorFrame
}

type rplb struct {
	Frontends map[string]string
	Backends  map[string][]*backend
}

var (
	fe = flag.String("f", "", "frontends")
	be = flag.String("b", "", "backends")
)

func main() {
	flag.Parse()

	createBackend := func(urlString string) (*backend, error) {
		urlParsed, err := url.Parse(urlString)
		if err != nil {
			return nil, err
		}

		proxy := httputil.NewSingleHostReverseProxy(urlParsed)

		return &backend{
			urlString,
			http.Client{Timeout: 2 * time.Second},
			false,
			proxy,
			[]monitorFrame{},
		}, nil
	}

	frontendName := strings.Split(*fe, ",")[0]
	frontendBackend := strings.Split(*fe, ",")[1]

	backendName := strings.Split(*be, ",")[0]
	backendUrls := strings.Split(*be, ",")[1:]

	backends := []*backend{}
	for _, b := range backendUrls {
		backend, err := createBackend(b)
		if err != nil {
			break
		}
		backends = append(backends, backend)
	}

	rplb := &rplb{
		Frontends: make(map[string]string),
		Backends:  make(map[string][]*backend),
	}

	rplb.Frontends[frontendName] = frontendBackend
	rplb.Backends[backendName] = backends

	go func(backends []*backend) {
		for {
			for _, b := range backends {
				go func(b *backend) {
					start := time.Now()
					res, err := b.client.Head(b.Url)
					latency := time.Since(start)
					if err != nil {
						b.live = false
						_ = fmt.Sprintf("%s\n\terror: %v\n", b.Url, err)
					} else {
						b.live = true
						_ = fmt.Sprintf("%s\n\tstatus code: %v\n\tlatency: %v\n", b.Url, res.StatusCode, latency)
					}
					b.Monitor = last20(append(b.Monitor, monitorFrame{b.live, latency}))
				}(b)
			}
			time.Sleep(1 * time.Second)
		}
	}(backends)

	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		host, _, _ := net.SplitHostPort(r.Host)
		f := rplb.Frontends[host]

		backends := rplb.Backends[f]

		var liveBackends []*backend
		for _, b := range backends {
			// fmt.Printf("%s\n\tlive: %v\n", b.url, b.live)
			if b.live {
				liveBackends = append(liveBackends, b)
			}
		}

		if len(liveBackends) > 0 {
			randBackend := rand.Intn(len(liveBackends))

			rw.Header().Add("proxy-url", liveBackends[randBackend].Url)
			liveBackends[randBackend].proxy.ServeHTTP(rw, r)
		} else {
			rw.WriteHeader(http.StatusServiceUnavailable)
		}
	})
	go http.ListenAndServe(":8080", nil)

	const tpl = `
<script>setTimeout('window.location.reload();', 1000);</script>
{{range $key, $value := .Backends}}
	<h2>{{$key}}</h2>
	{{range $value}}
		<div>
			<h3>{{.Url}}</h3>
			{{range .Monitor}}
				{{if .Live}}
					<span style="background-color: green;width: 20px;height: 20px;display: inline-block;" title="{{.Latency}}"></span>
				{{else}}
					<span style="background-color: red;width: 20px;height: 20px;display: inline-block;" title="{{.Latency}}"></span>
				{{end}}
			{{end}}
		</div>
	{{end}}
{{end}}
	`

	dashboardMux := http.NewServeMux()
	dashboardMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		t, err := template.New("dashboard").
			Parse(tpl)
		if err != nil {
			log.Fatal(err)
		}
		err = t.Execute(w, rplb)
		if err != nil {
			log.Fatal(err)
		}
	})
	dashboard := &http.Server{Addr: ":8000", Handler: dashboardMux}
	dashboard.ListenAndServe()
}

func last20(m []monitorFrame) []monitorFrame {
	l := int(math.Max(0, float64(len(m)-20)))
	return m[l:]
}
