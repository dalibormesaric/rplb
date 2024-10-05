package integration

import (
	"bufio"
	"net/http"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/dalibormesaric/rplb/internal/loadbalancing"
)

func TestIntegration(t *testing.T) {
	tests := []struct {
		algoirthm string
		expected  []string
	}{
		{
			algoirthm: loadbalancing.First,
			expected:  []string{"rplb_backend_hits{backend_name=\"http://172.17.0.1:8081\"} 10"},
		},
		{
			algoirthm: loadbalancing.Random,
			expected: []string{
				"rplb_backend_hits{backend_name=\"http://172.17.0.1:8081\"}",
				"rplb_backend_hits{backend_name=\"http://172.17.0.1:8082\"}",
				"rplb_backend_hits{backend_name=\"http://172.17.0.1:8083\"}",
			},
		},
		{
			algoirthm: loadbalancing.RoundRobin,
			expected: []string{
				"rplb_backend_hits{backend_name=\"http://172.17.0.1:8081\"} 4",
				"rplb_backend_hits{backend_name=\"http://172.17.0.1:8082\"} 3",
				"rplb_backend_hits{backend_name=\"http://172.17.0.1:8083\"} 3",
			},
		},
		{
			algoirthm: loadbalancing.Sticky,
			expected:  []string{"rplb_backend_hits{backend_name=\"http://172.17.0.1:8081\"} 10"},
		},
		{
			algoirthm: loadbalancing.LeastLoaded,
			expected: []string{
				"rplb_backend_hits{backend_name=\"http://172.17.0.1:8081\"} 4",
				"rplb_backend_hits{backend_name=\"http://172.17.0.1:8082\"} 3",
				"rplb_backend_hits{backend_name=\"http://172.17.0.1:8083\"} 3",
			},
		},
	}

	for _, test := range tests {
		setUp(test.algoirthm)

		var wg sync.WaitGroup

		for range 10 {
			wg.Add(1)
			time.Sleep(100 * time.Millisecond)
			go func() {
				defer wg.Done()
				http.Get("http://localhost:8080")
			}()
		}
		wg.Wait()

		resp, _ := http.Get("http://localhost:8000/metrics")
		scanner := bufio.NewScanner(resp.Body)
		var tested []string
		for scanner.Scan() {
			line := scanner.Text()

			for _, metric := range test.expected {
				if strings.Contains(line, metric) {
					tested = append(tested, metric)
				}
			}

			// TODO: test frontend hits is 10
			// TODO: retries is 0
		}

		if len(tested) != len(test.expected) {
			t.Errorf("Did not find all expected test cases.\n")
		}
		for _, testExpected := range test.expected {
			if !slices.Contains(tested, testExpected) {
				t.Errorf("Algorithm (%s): (%s)", test.algoirthm, testExpected)
			}
		}

		tearDown()
	}
}
