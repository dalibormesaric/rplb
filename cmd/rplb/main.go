package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/dalibormesaric/rplb/internal/backend"
	"github.com/dalibormesaric/rplb/internal/config"
	"github.com/dalibormesaric/rplb/internal/dashboard"
	"github.com/dalibormesaric/rplb/internal/frontend"
	"github.com/dalibormesaric/rplb/internal/loadbalancing"
	"github.com/dalibormesaric/rplb/internal/reverseproxy"
)

var (
	fe = flag.String("f", "", "Comma-separated list of Frontend Hostname and BackendPool Name pairs. (example \"frontend.local,backend\")")
	be = flag.String("b", "", "Comma-separated list of BackendPool Name and URL pairs. (example \"backend,10.0.0.1:1234\")")
	a  = flag.String("a", loadbalancing.Sticky, fmt.Sprintf("Algorithm used for loadbalancing. Choose from: %s, %s, %s or %s.", loadbalancing.First, loadbalancing.Random, loadbalancing.RoundRobin, loadbalancing.Sticky))
)

func main() {
	flag.Parse()

	frontends, err := frontend.CreateFrontends(*fe)
	if err != nil {
		log.Fatalf("Create frontends: %s", err)
	}

	bp, err := backend.NewBackendPool(*be)
	if err != nil {
		log.Fatalf("NewBackendPool: %s", err)
	}
	messages := bp.Monitor()

	algo, err := loadbalancing.NewAlgorithm(*a)
	if err != nil {
		log.Fatalf("NewAlgorithm: %s", err)
	}
	go reverseproxy.ListenAndServe(frontends, bp, algo, messages)

	dashboard.ListenAndServe(frontends, bp, messages, config.Version)
}
