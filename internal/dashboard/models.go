package dashboard

import (
	"github.com/dalibormesaric/rplb/internal/backend"
	"github.com/dalibormesaric/rplb/internal/frontend"
)

type BaseModel struct {
	SelectedMenu string
	Version      string
}

type MonitorModel struct {
	BaseModel
	Backends backend.BackendPool
}

type TrafficModel struct {
	BaseModel
	Frontends frontend.Frontends
	Backends  backend.BackendPool
}
