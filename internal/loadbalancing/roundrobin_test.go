package loadbalancing

import (
	"fmt"
	"testing"

	"github.com/dalibormesaric/rplb/internal/backend"
)

const (
	bname string = "test"
	be1   string = "http://b:1234"
	be2   string = "http://b:1235"
	be3   string = "http://b:1236"
)

func TestRoundRobinSequence(t *testing.T) {
	createBackends := func() []*backend.Backend {
		r, _ := backend.CreateBackends(fmt.Sprintf("%s,%s,%s,%s,%s,%s", bname, be1, bname, be2, bname, be3))
		return r[bname]
	}()

	var test = struct {
		liveBackends []*backend.Backend
		expected     []string
	}{
		liveBackends: createBackends,
		expected:     []string{be1, be2, be3, be1, be2, be3, be1},
	}

	roundRobin := NewAlgorithm(RoundRobin)
	for _, expected := range test.expected {
		b := roundRobin.Get(nil, test.liveBackends)
		if b.URL.String() != expected {
			t.Errorf("wrong backend: want (%s) got (%s)", expected, b.URL.String())
		}
	}
}

func TestRoundRobinGetNil(t *testing.T) {
	var test = struct {
		liveBackends []*backend.Backend
		expected     *backend.Backend
	}{
		liveBackends: []*backend.Backend{},
		expected:     nil,
	}

	roundRobin := NewAlgorithm(RoundRobin)
	b := roundRobin.Get(nil, test.liveBackends)
	if b != test.expected {
		t.Errorf("wrong backend: want (%v) got (%v)", test.expected, b)
	}
}
