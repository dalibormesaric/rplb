package reverseproxy

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dalibormesaric/rplb/internal/frontend"
)

func TestReverseProxyWithNoFrontends(t *testing.T) {
	rp := &reverseProxy{}
	server := httptest.NewServer(http.HandlerFunc(rp.reverseProxyAndLoadBalance))
	defer server.Close()

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != 404 {
		t.Errorf("wrong status code: want (404) got (%d)\n", resp.StatusCode)
	}
}

func TestReverseProxyWithFrontends(t *testing.T) {
	f, err := frontend.CreateFrontends("127.0.0.1,b")
	if err != nil {
		t.Error(err)
	}
	rp := &reverseProxy{
		frontends: f,
	}
	server := httptest.NewServer(http.HandlerFunc(rp.reverseProxyAndLoadBalance))
	defer server.Close()

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != 503 {
		t.Errorf("wrong status code: want (503) got (%d)\n", resp.StatusCode)
	}
}
