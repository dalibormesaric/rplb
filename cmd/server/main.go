package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/dalibormesaric/rplb/internal/backend"
	"github.com/dalibormesaric/rplb/internal/dashboard"
	"github.com/dalibormesaric/rplb/internal/frontend"
	"github.com/dalibormesaric/rplb/internal/reverseproxy"
	"github.com/dalibormesaric/rplb/internal/server"
)

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

	dashboard.ListenAndServe(frontends, backends)
}
