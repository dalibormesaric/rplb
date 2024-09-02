package loadbalancing

import (
	"fmt"
	"testing"

	"github.com/dalibormesaric/rplb/internal/backend"
)

const (
	firstBpName string = First
	firstB1     string = "http://b:1234"
	firstB2     string = "http://b:1235"
	firstB3     string = "http://b:1236"
)

func TestFirstSequence(t *testing.T) {
	bs := func() []*backend.Backend {
		bp, _ := backend.NewBackendPool(fmt.Sprintf("%s,%s,%s,%s,%s,%s", firstBpName, firstB1, firstBpName, firstB2, firstBpName, firstB3))
		return bp[firstBpName]
	}()

	var test = struct {
		bs       []*backend.Backend
		expected []string
	}{
		bs:       bs,
		expected: []string{firstB1, firstB1, firstB1, firstB1, firstB1, firstB1, firstB1},
	}

	first, _ := NewAlgorithm(First)
	for _, expected := range test.expected {
		b := first.Get("", test.bs)
		if b.URL.String() != expected {
			t.Errorf("wrong backend: want (%s) got (%s)", expected, b.URL.String())
		}
	}
}

func TestFirstGetNil(t *testing.T) {
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
		first, _ := NewAlgorithm(First)
		b := first.Get("", test.bs)
		if b != test.expected {
			t.Errorf("wrong backend: want (%v) got (%v)", test.expected, b)
		}
	}
}
