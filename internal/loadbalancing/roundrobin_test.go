package loadbalancing

import (
	"fmt"
	"testing"

	"github.com/dalibormesaric/rplb/internal/backend"
)

const (
	roundRobinBpName string = RoundRobin
	roundRobinB1     string = "http://a:1234"
	roundRobinB2     string = "http://b:1234"
	roundRobinB3     string = "http://c:1234"
)

func TestRoundRobinSequence(t *testing.T) {
	bs := func() []*backend.Backend {
		bp, _ := backend.NewBackendPool(fmt.Sprintf("%s,%s,%s,%s,%s,%s", roundRobinBpName, roundRobinB1, roundRobinBpName, roundRobinB2, roundRobinBpName, roundRobinB3))
		return bp[roundRobinBpName]
	}()

	var test = struct {
		bs       []*backend.Backend
		expected []string
	}{
		bs:       bs,
		expected: []string{roundRobinB1, roundRobinB2, roundRobinB3, roundRobinB1, roundRobinB2, roundRobinB3, roundRobinB1},
	}

	roundRobin, _ := NewAlgorithm(RoundRobin)
	for _, expected := range test.expected {
		b, _ := roundRobin.Get("", test.bs)
		if b.URL.String() != expected {
			t.Errorf("wrong backend: want (%s) got (%s)", expected, b.URL.String())
		}
	}
}

func TestRoundRobinGetNil(t *testing.T) {
	var tests = []struct {
		bs       []*backend.Backend
		expected *backend.Backend
	}{
		{
			bs:       nil,
			expected: nil,
		},
		{
			bs:       []*backend.Backend{},
			expected: nil,
		},
	}

	for _, test := range tests {
		roundRobin, _ := NewAlgorithm(RoundRobin)
		b, _ := roundRobin.Get("", test.bs)
		if b != test.expected {
			t.Errorf("wrong backend: want (%v) got (%v)", test.expected, b)
		}
	}
}
