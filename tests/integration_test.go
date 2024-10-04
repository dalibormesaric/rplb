package integration

import (
	"bufio"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/dalibormesaric/rplb/internal/loadbalancing"
)

func TestIntegration(t *testing.T) {
	setUp(loadbalancing.First)

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
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "rplb_backend_hits") {
			if !strings.Contains(line, "rplb_backend_hits{backend_name=\"http://172.17.0.1:8081\"} 10") {
				t.Error(line)
			}
		}
	}

	tearDown()
}
