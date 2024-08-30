package loadbalancing

import (
	"fmt"
	"testing"

	"github.com/dalibormesaric/rplb/internal/backend"
)

const (
	stickyBpName string = Sticky
	stickyB1     string = "http://b:1234"
	stickyB2     string = "http://b:1235"
	stickyB3     string = "http://b:1236"
	stickyC1     string = "10.0.0.1:1234"
	stickyC2     string = "10.0.0.2:1235"
)

func TestStickySequence(t *testing.T) {
	bs := func() []*backend.Backend {
		bp, _ := backend.NewBackendPool(fmt.Sprintf("%s,%s,%s,%s,%s,%s", stickyBpName, stickyB1, stickyBpName, stickyB2, stickyBpName, stickyB3))
		return bp[stickyBpName]
	}()

	var test = struct {
		bs       []*backend.Backend
		clients  []string
		expected []string
	}{
		bs:       bs,
		clients:  []string{stickyC1, stickyC1, stickyC2, stickyC2, stickyC1, stickyC2, stickyC1},
		expected: []string{stickyB1, stickyB1, stickyB2, stickyB2, stickyB1, stickyB2, stickyB1},
	}

	sticky, _ := NewAlgorithm(Sticky)
	for i, expected := range test.expected {
		b := sticky.Get(test.clients[i], test.bs)
		if b.URL.String() != expected {
			t.Errorf("wrong backend for client (%s): want (%s) got (%s)", test.clients[i], expected, b.URL.String())
		}
	}
}

func TestStickyGetNil(t *testing.T) {
	var test = struct {
		bs       []*backend.Backend
		expected *backend.Backend
	}{
		bs:       []*backend.Backend{},
		expected: nil,
	}

	sticky, _ := NewAlgorithm(Sticky)
	b := sticky.Get("", test.bs)
	if b != test.expected {
		t.Errorf("wrong backend: want (%v) got (%v)", test.expected, b)
	}
}
