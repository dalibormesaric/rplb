package integration

import (
	"bufio"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestIntegration(t *testing.T) {
	err := exec.Command("docker", "compose", "-f", "../example/compose.yaml", "up", "-d", "rplb").Run()
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup

	for range 10 {
		wg.Add(1)
		time.Sleep(100 * time.Millisecond)
		go func() {
			defer wg.Done()
			http.Get("http://host.docker.internal:8080")
		}()
	}
	wg.Wait()

	resp, _ := http.Get("http://host.docker.internal:8000/metrics")
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "rplb_backend_hits") {
			if !strings.Contains(line, "rplb_backend_hits 10") {
				t.Error(line)
			}
		}
	}

	err = exec.Command("docker", "compose", "-f", "../example/compose.yaml", "down").Run()
	if err != nil {
		log.Fatal(err)
	}
}
