package loadbalancing

import (
	"math/rand/v2"

	"github.com/dalibormesaric/rplb/internal/backend"
)

type random struct {
}

var _ Algorithm = (*random)(nil)

func (*random) GetNext(_ string, backends []*backend.Backend) (backend *backend.Backend, _ func()) {
	l := len(backends)
	if l == 0 {
		return nil, nil
	}

	randBackend := rand.IntN(l)
	liveBackend := backends[randBackend]
	return liveBackend, nil
}
