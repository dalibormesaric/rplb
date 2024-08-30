package loadbalancing

import (
	"fmt"
	"testing"

	"github.com/dalibormesaric/rplb/internal/backend"
)

const (
	stickyBpName string = Sticky
	b1           string = "http://b:1234"
	b2           string = "http://b:1235"
	b3           string = "http://b:1236"
	c1           string = "192.168.0.1:1234"
	c2           string = "192.168.0.2:1235"
)

func TestStickySequence(t *testing.T) {
	bs := func() []*backend.Backend {
		bp, _ := backend.NewBackendPool(fmt.Sprintf("%s,%s,%s,%s,%s,%s", stickyBpName, b1, stickyBpName, b2, stickyBpName, b3))
		return bp[stickyBpName]
	}()

	var test = struct {
		bs       []*backend.Backend
		clients  []string
		expected []string
	}{
		bs:       bs,
		clients:  []string{c1, c1, c2, c2, c1, c2, c1},
		expected: []string{b1, b1, b2, b2, b1, b2, b1},
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
