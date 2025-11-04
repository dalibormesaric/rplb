package reverseproxy

import (
	"embed"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/dalibormesaric/rplb/internal/backend"
	"github.com/dalibormesaric/rplb/internal/dashboard"
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

	slog.Info("Reverse Proxy listening on :8080")
	http.ListenAndServe(":8080", rpMux)
}

//go:embed static/*.html
var static embed.FS

func (rp *reverseProxy) reverseProxyAndLoadBalance(w http.ResponseWriter, r *http.Request) {
	host := hostname(r)
	f := rp.frontends.Get(host)
	if f == nil {
		slog.Warn("Unknown frontend host", "host", host)
		http.ServeFileFS(w, r, static, "static/404.html")
		return
	}

	if rp.messages != nil {
		tf := TrafficFrame{Type: "traffic-fe", Name: host, Hits: f.IncHits()}
		rp.messages <- tf
	}
	dashboard.FrontendHits.WithLabelValues(host).Inc()

	retryTimeout := 500 * time.Millisecond
	retryAmount := 5
	for retryCurrent := range retryAmount {
		liveBackends := backend.GetLive(rp.bp[f.BackendName])
		liveBackend, afterBackendResponse := rp.loadbalancing.GetNext(r.RemoteAddr, liveBackends)
		if liveBackend == nil || retryCurrent+1 == retryAmount {
			if liveBackend == nil {
				slog.Warn("No live backends", "host", host)
			}
			if retryCurrent+1 == retryAmount {
				slog.Warn("Backends unavailable after retries", "host", host)
			}
			http.ServeFileFS(w, r, static, "static/503.html")
			break
		}
		// TODO: Consider X-Forwarded-For or similar header
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
			dashboard.BackendHits.WithLabelValues(liveBackend.GetPoolName(), liveBackend.URL.String()).Inc()
			break
		}

		time.Sleep(retryTimeout)
		retryTimeout *= 2
		dashboard.BackendRetries.Inc()
	}
}

func hostname(r *http.Request) string {
	host := strings.TrimSpace(r.Host)

	if strings.Contains(host, "://") {
		if u, err := url.Parse(host); err == nil {
			host = u.Host
		}
	}

	h, _, err := net.SplitHostPort(host)
	if err == nil {
		return h
	}

	return host
}
