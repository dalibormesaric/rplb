package loadbalancing

import (
	"fmt"
	"testing"

	"github.com/dalibormesaric/rplb/internal/backend"
)

const (
	beName string = "test"
	b1     string = "http://b:1234"
	b2     string = "http://b:1235"
	b3     string = "http://b:1236"
)

func TestRoundRobinSequence(t *testing.T) {
	bs := func() []*backend.Backend {
		be, _ := backend.CreateBackends(fmt.Sprintf("%s,%s,%s,%s,%s,%s", beName, b1, beName, b2, beName, b3))
		return be[beName]
	}()

	var test = struct {
		bs       []*backend.Backend
		expected []string
	}{
		bs:       bs,
		expected: []string{b1, b2, b3, b1, b2, b3, b1},
	}

	roundRobin, _ := NewAlgorithm(RoundRobin)
	for _, expected := range test.expected {
		b := roundRobin.Get(nil, test.bs)
		if b.URL.String() != expected {
			t.Errorf("wrong backend: want (%s) got (%s)", expected, b.URL.String())
		}
	}
}

func TestRoundRobinGetNil(t *testing.T) {
	var test = struct {
		bs       []*backend.Backend
		expected *backend.Backend
	}{
		bs:       []*backend.Backend{},
		expected: nil,
	}

	roundRobin, _ := NewAlgorithm(RoundRobin)
	b := roundRobin.Get(nil, test.bs)
	if b != test.expected {
		t.Errorf("wrong backend: want (%v) got (%v)", test.expected, b)
	}
}
