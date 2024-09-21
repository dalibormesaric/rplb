package loadbalancing

import (
	"fmt"
	"testing"

	"github.com/dalibormesaric/rplb/internal/backend"
)

const (
	randomBpName string = Random
	randomB1     string = "http://a:1234"
	randomB2     string = "http://b:1234"
	randomB3     string = "http://c:1234"
)

func TestRandomSequence(t *testing.T) {
	getBackends := func() []*backend.Backend {
		bp, _ := backend.NewBackendPool(fmt.Sprintf("%s,%s,%s,%s,%s,%s", randomBpName, randomB1, randomBpName, randomB2, randomBpName, randomB3))
		return bp[randomBpName]
	}()

	var test = struct {
		backends []*backend.Backend
	}{
		backends: getBackends,
	}

	random, _ := NewAlgorithm(Random)
	for range 7 {
		b, _ := random.GetNext("", test.backends)
		if b.URL.String() != randomB1 && b.URL.String() != randomB2 && b.URL.String() != randomB3 {
			t.Errorf("wrong backend: want (%s, %s or %s) got (%s)", randomB1, randomB2, randomB3, b.URL.String())
		}
	}
}

func TestRandomGetNil(t *testing.T) {
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
		random, _ := NewAlgorithm(Random)
		b, _ := random.GetNext("", test.backends)
		if b != test.expected {
			t.Errorf("wrong backend: want (%v) got (%v)", test.expected, b)
		}
	}
}
