package loadbalancing

import (
	"math/rand/v2"

	"github.com/dalibormesaric/rplb/internal/backend"
)

type random struct {
}

var _ Algorithm = (*random)(nil)

func (_ *random) Get(_ string, backends []*backend.Backend) *backend.Backend {
	n := len(backends)
	if n == 0 {
		return nil
	}

	randBackend := rand.IntN(n)
	liveBackend := backends[randBackend]
	return liveBackend
}
