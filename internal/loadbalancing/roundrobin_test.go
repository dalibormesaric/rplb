package loadbalancing

import (
	"fmt"
	"testing"

	"github.com/dalibormesaric/rplb/internal/backend"
)

const (
	roundRobinPool1 string = RoundRobin + "1"
	roundRobinPool2 string = RoundRobin + "2"
	roundRobinB1    string = "http://a:1234"
	roundRobinB2    string = "http://b:1234"
	roundRobinB3    string = "http://c:1234"
)

func TestRoundRobinSequence(t *testing.T) {
	getBackendsForPool := func(poolName string) []*backend.Backend {
		bp, _ := backend.NewBackendPool(fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s,%s,%s,%s",
			roundRobinPool1, roundRobinB1, roundRobinPool1, roundRobinB2, roundRobinPool1, roundRobinB3,
			roundRobinPool2, roundRobinB1, roundRobinPool2, roundRobinB2))
		return bp[poolName]
	}

	var tests = []struct {
		bs       []*backend.Backend
		expected []string
	}{
		{
			bs:       getBackendsForPool(roundRobinPool1),
			expected: []string{roundRobinB1, roundRobinB2, roundRobinB3, roundRobinB1, roundRobinB2, roundRobinB3, roundRobinB1},
		},
		{
			bs:       getBackendsForPool(roundRobinPool2),
			expected: []string{roundRobinB1, roundRobinB2, roundRobinB1, roundRobinB2, roundRobinB1},
		},
	}

	for _, test := range tests {
		roundRobin, _ := NewAlgorithm(RoundRobin)
		for _, expected := range test.expected {
			b, _ := roundRobin.GetNext("", test.bs)
			if b.URL.String() != expected {
				t.Errorf("wrong backend: want (%s) got (%s)", expected, b.URL.String())
			}
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
		b, _ := roundRobin.GetNext("", test.bs)
		if b != test.expected {
			t.Errorf("wrong backend: want (%v) got (%v)", test.expected, b)
		}
	}
}
