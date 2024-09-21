package loadbalancing

import (
	"fmt"
	"testing"

	"github.com/dalibormesaric/rplb/internal/backend"
)

const (
	firstPool1 string = First + "1"
	firstPool2 string = First + "2"
	firstB1    string = "http://a:1234"
	firstB2    string = "http://b:1234"
	firstB3    string = "http://c:1234"
)

func TestFirstSequence(t *testing.T) {
	getBackendsForPool := func(poolName string) []*backend.Backend {
		backendPool, _ := backend.NewBackendPool(fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s,%s,%s,%s",
			firstPool1, firstB1, firstPool1, firstB2, firstPool1, firstB3,
			firstPool2, firstB2, firstPool2, firstB3))
		return backendPool[poolName]
	}

	var tests = []struct {
		backends []*backend.Backend
		expected []string
	}{
		{
			backends: getBackendsForPool(firstPool1),
			expected: []string{firstB1, firstB1, firstB1, firstB1, firstB1, firstB1, firstB1},
		},

		{
			backends: getBackendsForPool(firstPool2),
			expected: []string{firstB2, firstB2, firstB2, firstB2, firstB2, firstB2, firstB2},
		},
	}

	for _, test := range tests {
		first, _ := NewAlgorithm(First)
		for _, expected := range test.expected {
			b, _ := first.GetNext("", test.backends)
			if b.URL.String() != expected {
				t.Errorf("wrong backend: want (%s) got (%s)", expected, b.URL.String())
			}
		}
	}
}

func TestFirstGetNil(t *testing.T) {
	var tests = []struct {
		backends []*backend.Backend
		expected *backend.Backend
	}{
		{
			backends: nil,
			expected: nil,
		},
		{
			backends: []*backend.Backend{},
			expected: nil,
		},
	}

	for _, test := range tests {
		first, _ := NewAlgorithm(First)
		b, _ := first.GetNext("", test.backends)
		if b != test.expected {
			t.Errorf("wrong backend: want (%v) got (%v)", test.expected, b)
		}
	}
}
