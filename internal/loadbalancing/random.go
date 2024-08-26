package loadbalancing

import (
	"math/rand/v2"
	"net/http"

	"github.com/dalibormesaric/rplb/internal/backend"
)

type random struct {
}

var _ Algorithm = (*random)(nil)

func (_ *random) Get(r *http.Request, liveBackends []*backend.Backend) *backend.Backend {
	n := len(liveBackends)
	if n == 0 {
		return nil
	}

	randBackend := rand.IntN(n)
	liveBackend := liveBackends[randBackend]
	return liveBackend
}
