package loadbalancing

import (
	"fmt"
	"testing"

	"github.com/dalibormesaric/rplb/internal/backend"
)

const (
	leastLoadedBpName string = RoundRobin
	leastLoadedB1     string = "http://a:1234"
	leastLoadedB2     string = "http://b:1234"
	leastLoadedB3     string = "http://c:1234"
)

func TestGet(t *testing.T) {
	bsf := func() []*backend.Backend {
		bp, _ := backend.NewBackendPool(fmt.Sprintf("%s,%s,%s,%s,%s,%s", leastLoadedBpName, leastLoadedB1, leastLoadedBpName, leastLoadedB2, leastLoadedBpName, leastLoadedB3))
		return bp[leastLoadedBpName]
	}

	var tests = []struct {
		callbackNever    bool
		loadForBackend   []int
		expectedBackends []string
		expectedStates   [][]int
	}{
		{
			callbackNever:    true,
			loadForBackend:   []int{2, 1, 1},
			expectedBackends: []string{leastLoadedB2, leastLoadedB3, leastLoadedB1, leastLoadedB3, leastLoadedB2},
			expectedStates:   [][]int{{2, 2, 1}, {2, 2, 2}, {3, 2, 2}, {3, 2, 3}, {3, 3, 3}},
		},
		{
			callbackNever:    false,
			loadForBackend:   []int{0, 0, 0},
			expectedBackends: []string{leastLoadedB1, leastLoadedB2, leastLoadedB3, leastLoadedB1, leastLoadedB2},
			expectedStates:   [][]int{{1, 0, 0}, {0, 1, 0}, {0, 0, 1}, {1, 0, 0}, {0, 1, 0}},
		},
		{
			callbackNever:    false,
			loadForBackend:   []int{2, 1, 1},
			expectedBackends: []string{leastLoadedB2, leastLoadedB3, leastLoadedB2, leastLoadedB3, leastLoadedB2},
			expectedStates:   [][]int{{2, 2, 1}, {2, 1, 2}, {2, 2, 1}, {2, 1, 2}, {2, 2, 1}},
		},
	}

	for _, test := range tests {
		bs := bsf()
		leastloaded := &leastLoaded{state: &leastLoadedState{
			loadForBackend:        make(map[string]int),
			roundRobinForPoolLoad: make(map[string]int),
		}}
		for j := range len(test.loadForBackend) {
			leastloaded.state.loadForBackend[bs[j].Name] = test.loadForBackend[j]
		}

		for i, expectedBackend := range test.expectedBackends {
			b, f := leastloaded.GetNext("", bs)
			if b.URL.String() != expectedBackend {
				t.Errorf("Wrong backend at step (%d): want (%s) got (%s)\n", i, expectedBackend, b.URL.String())
			}
			for j := range len(test.loadForBackend) {
				if test.expectedStates[i][j] != leastloaded.state.loadForBackend[bs[j].Name] {
					t.Errorf("Wrong loadForBackend[(%d)] at step (%d): want (%d) got (%d)\n", j, i, test.expectedStates[i][j], leastloaded.state.loadForBackend[bs[j].Name])
				}
			}

			if !test.callbackNever {
				f()
			}
		}
	}
}
