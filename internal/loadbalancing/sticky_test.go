package loadbalancing

import (
	"fmt"
	"testing"

	"github.com/dalibormesaric/rplb/internal/backend"
)

const (
	stickyPool1 string = Sticky + "1"
	stickyPool2 string = Sticky + "2"
	stickyPool3 string = Sticky + "3"
	stickyB1    string = "http://a:1234"
	stickyB2    string = "http://b:1234"
	stickyB3    string = "http://c:1234"
	stickyB4    string = "http://d:1234"
	stickyB5    string = "http://e:1234"
	stickyC1    string = "192.168.0.10:1234"
	stickyC2    string = "192.168.0.11:1234"
)

func TestStickySequence(t *testing.T) {
	getBackendsForPool := func(poolName string) []*backend.Backend {
		backendPool, _ := backend.NewBackendPool(fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s",
			stickyPool1, stickyB1, stickyPool1, stickyB2, stickyPool1, stickyB3,
			stickyPool2, stickyB4, stickyPool2, stickyB5,
			stickyPool3, stickyB1))
		return backendPool[poolName]
	}

	var tests = []struct {
		otherPool []*backend.Backend
		backends  []*backend.Backend
		clients   []string
		expected  []string
	}{
		{
			otherPool: getBackendsForPool(stickyPool2),
			backends:  getBackendsForPool(stickyPool1),
			clients:   []string{stickyC1, stickyC1, stickyC2, stickyC2, stickyC1, stickyC2, stickyC1},
			expected:  []string{stickyB1, stickyB1, stickyB2, stickyB2, stickyB1, stickyB2, stickyB1},
		},
		{
			otherPool: getBackendsForPool(stickyPool1),
			backends:  getBackendsForPool(stickyPool2),
			clients:   []string{stickyC1, stickyC1, stickyC2, stickyC2, stickyC1, stickyC2, stickyC1},
			expected:  []string{stickyB4, stickyB4, stickyB5, stickyB5, stickyB4, stickyB5, stickyB4},
		},
		{
			otherPool: getBackendsForPool(stickyPool1),
			backends:  getBackendsForPool(stickyPool3),
			clients:   []string{stickyC1, stickyC1, stickyC2, stickyC2, stickyC1, stickyC2, stickyC1},
			expected:  []string{stickyB1, stickyB1, stickyB1, stickyB1, stickyB1, stickyB1, stickyB1},
		},
		{
			otherPool: getBackendsForPool(stickyPool2),
			backends:  getBackendsForPool(stickyPool3),
			clients:   []string{stickyC1, stickyC1, stickyC2, stickyC2, stickyC1, stickyC2, stickyC1},
			expected:  []string{stickyB1, stickyB1, stickyB1, stickyB1, stickyB1, stickyB1, stickyB1},
		},
	}

	for _, test := range tests {
		sticky, _ := NewAlgorithm(Sticky)
		for i, expected := range test.expected {
			// trigger other pool backends to test that it does not affect testing backends
			sticky.GetNext(test.clients[i], test.otherPool)
			b, _ := sticky.GetNext(test.clients[i], test.backends)
			if b.URL.String() != expected {
				t.Errorf("wrong backend for client (%s): want (%s) got (%s)", test.clients[i], expected, b.URL.String())
			}
		}
	}
}

func TestStickyGetNil(t *testing.T) {
	var tests = []struct {
		backends   []*backend.Backend
		remoteAddr string
		expected   *backend.Backend
	}{
		{
			backends:   nil,
			remoteAddr: stickyC1,
			expected:   nil,
		},
		{
			backends:   []*backend.Backend{},
			remoteAddr: stickyC1,
			expected:   nil,
		},
		{
			backends:   []*backend.Backend{{}},
			remoteAddr: "",
			expected:   nil,
		},
		{
			backends:   []*backend.Backend{{}},
			remoteAddr: "wrong",
			expected:   nil,
		},
		{
			backends:   []*backend.Backend{{}},
			remoteAddr: "1234",
			expected:   nil,
		},
		{
			backends:   []*backend.Backend{{}},
			remoteAddr: "10.0.0.1",
			expected:   nil,
		},
	}

	for _, test := range tests {
		sticky, _ := NewAlgorithm(Sticky)
		b, _ := sticky.GetNext(test.remoteAddr, test.backends)
		if b != test.expected {
			t.Errorf("wrong backend: want (%v) got (%v)", test.expected, b)
		}
	}
}
