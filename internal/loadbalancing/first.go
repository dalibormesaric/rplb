package loadbalancing

import (
	"net/http"

	"github.com/dalibormesaric/rplb/internal/backend"
)

type first struct {
}

var _ Algorithm = (*first)(nil)

func (_ *first) Get(r *http.Request, liveBackends []*backend.Backend) *backend.Backend {
	n := len(liveBackends)
	if n == 0 {
		return nil
	}

	liveBackend := liveBackends[0]
	return liveBackend
}
