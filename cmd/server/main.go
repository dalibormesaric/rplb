package main

import (
	"flag"
	"log"

	"github.com/dalibormesaric/rplb/internal/backend"
	"github.com/dalibormesaric/rplb/internal/dashboard"
	"github.com/dalibormesaric/rplb/internal/frontend"
	"github.com/dalibormesaric/rplb/internal/reverseproxy"
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

	dashboard.ListenAndServe(frontends, backends, monitor.Messages)
}
