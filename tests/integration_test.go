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
		algoirthm        string
		expectedBackends []string
		expectedFrontend string
		expectedRetries  string
	}{
		{
			algoirthm:        loadbalancing.First,
			expectedBackends: []string{"rplb_backend_hits{backend_bool_name=\"example\",backend_url=\"http://172.17.0.1:8081\"} 10"},
			expectedFrontend: "rplb_frontend_hits{frontend_host=\"localhost\"} 10",
			expectedRetries:  "rplb_backend_retries 0",
		},
		{
			algoirthm: loadbalancing.Random,
			expectedBackends: []string{
				"rplb_backend_hits{backend_bool_name=\"example\",backend_url=\"http://172.17.0.1:8081\"}",
				"rplb_backend_hits{backend_bool_name=\"example\",backend_url=\"http://172.17.0.1:8082\"}",
				"rplb_backend_hits{backend_bool_name=\"example\",backend_url=\"http://172.17.0.1:8083\"}",
			},
			expectedFrontend: "rplb_frontend_hits{frontend_host=\"localhost\"} 10",
			expectedRetries:  "rplb_backend_retries 0",
		},
		{
			algoirthm: loadbalancing.RoundRobin,
			expectedBackends: []string{
				"rplb_backend_hits{backend_bool_name=\"example\",backend_url=\"http://172.17.0.1:8081\"} 4",
				"rplb_backend_hits{backend_bool_name=\"example\",backend_url=\"http://172.17.0.1:8082\"} 3",
				"rplb_backend_hits{backend_bool_name=\"example\",backend_url=\"http://172.17.0.1:8083\"} 3",
			},
			expectedFrontend: "rplb_frontend_hits{frontend_host=\"localhost\"} 10",
			expectedRetries:  "rplb_backend_retries 0",
		},
		{
			algoirthm:        loadbalancing.Sticky,
			expectedBackends: []string{"rplb_backend_hits{backend_bool_name=\"example\",backend_url=\"http://172.17.0.1:8081\"} 10"},
			expectedFrontend: "rplb_frontend_hits{frontend_host=\"localhost\"} 10",
			expectedRetries:  "rplb_backend_retries 0",
		},
		{
			algoirthm: loadbalancing.LeastLoaded,
			expectedBackends: []string{
				"rplb_backend_hits{backend_bool_name=\"example\",backend_url=\"http://172.17.0.1:8081\"} 4",
				"rplb_backend_hits{backend_bool_name=\"example\",backend_url=\"http://172.17.0.1:8082\"} 3",
				"rplb_backend_hits{backend_bool_name=\"example\",backend_url=\"http://172.17.0.1:8083\"} 3",
			},
			expectedFrontend: "rplb_frontend_hits{frontend_host=\"localhost\"} 10",
			expectedRetries:  "rplb_backend_retries 0",
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

			for _, metric := range test.expectedBackends {
				if strings.Contains(line, metric) {
					tested = append(tested, metric)
				}
			}

			if strings.Contains(line, test.expectedFrontend) {
				tested = append(tested, test.expectedFrontend)
			}

			if strings.Contains(line, test.expectedRetries) {
				tested = append(tested, test.expectedRetries)
			}
		}

		// remove 2 because of expectedFrontend and expectedRetries
		if len(tested)-2 != len(test.expectedBackends) {
			t.Errorf("Did not find all expected test cases.\n")
		}
		for _, testExpected := range test.expectedBackends {
			if !slices.Contains(tested, testExpected) {
				t.Errorf("Algorithm (%s): (%s)", test.algoirthm, testExpected)
			}
		}
		if !slices.Contains(tested, test.expectedFrontend) {
			t.Errorf("Algorithm (%s): (%s)", test.algoirthm, test.expectedFrontend)
		}
		if !slices.Contains(tested, test.expectedRetries) {
			t.Errorf("Algorithm (%s): (%s)", test.algoirthm, test.expectedRetries)
		}

		tearDown()
	}
}
