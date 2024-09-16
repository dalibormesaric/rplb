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
	bs := func() []*backend.Backend {
		bp, _ := backend.NewBackendPool(fmt.Sprintf("%s,%s,%s,%s,%s,%s", randomBpName, randomB1, randomBpName, randomB2, randomBpName, randomB3))
		return bp[randomBpName]
	}()

	var test = struct {
		bs []*backend.Backend
	}{
		bs: bs,
	}

	random, _ := NewAlgorithm(Random)
	for range 7 {
		b, _ := random.Get("", test.bs)
		if b.URL.String() != randomB1 && b.URL.String() != randomB2 && b.URL.String() != randomB3 {
			t.Errorf("wrong backend: want (%s, %s or %s) got (%s)", randomB1, randomB2, randomB3, b.URL.String())
		}
	}
}

func TestRandomGetNil(t *testing.T) {
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
		random, _ := NewAlgorithm(Random)
		b, _ := random.Get("", test.bs)
		if b != test.expected {
			t.Errorf("wrong backend: want (%v) got (%v)", test.expected, b)
		}
	}
}
