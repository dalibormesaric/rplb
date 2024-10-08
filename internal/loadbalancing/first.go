package loadbalancing

import (
	"github.com/dalibormesaric/rplb/internal/backend"
)

type first struct {
}

var _ Algorithm = (*first)(nil)

func (*first) GetNext(_ string, backends []*backend.Backend) (backend *backend.Backend, _ func()) {
	l := len(backends)
	if l == 0 {
		return nil, nil
	}

	liveBackend := backends[0]
	return liveBackend, nil
}
