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

	backends, err := backend.CreateBackends(*be)
	if err != nil {
		log.Fatalf("Create backends: %s", err)
	}
	messages := backends.Monitor()

	algo, err := loadbalancing.NewAlgorithm(loadbalancing.Sticky)
	if err != nil {
		log.Fatalf("New algorithm: %s", err)
	}
	go reverseproxy.ListenAndServe(frontends, backends, algo, messages)

	dashboard.ListenAndServe(frontends, backends, messages, config.Version)
}
