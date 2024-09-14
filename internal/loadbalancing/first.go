package loadbalancing

import (
	"github.com/dalibormesaric/rplb/internal/backend"
)

type first struct {
}

var _ Algorithm = (*first)(nil)

func (*first) Get(_ string, backends []*backend.Backend) *backend.Backend {
	n := len(backends)
	if n == 0 {
		return nil
	}

	liveBackend := backends[0]
	return liveBackend
}
