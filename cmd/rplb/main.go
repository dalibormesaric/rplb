package main

import (
	"flag"
	"log"

	"github.com/dalibormesaric/rplb/internal/backend"
	"github.com/dalibormesaric/rplb/internal/config"
	"github.com/dalibormesaric/rplb/internal/dashboard"
	"github.com/dalibormesaric/rplb/internal/frontend"
	"github.com/dalibormesaric/rplb/internal/loadbalancing"
	"github.com/dalibormesaric/rplb/internal/reverseproxy"
)

var (
	fe = flag.String("f", "", "frontends")
	be = flag.String("b", "", "backends")
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

	algo, err := loadbalancing.NewAlgorithm(loadbalancing.Sticky)
	if err != nil {
		log.Fatalf("NewAlgorithm: %s", err)
	}
	go reverseproxy.ListenAndServe(frontends, bp, algo, messages)

	dashboard.ListenAndServe(frontends, bp, messages, config.Version)
}
