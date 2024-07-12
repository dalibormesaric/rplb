package reverseproxy

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestReverseProxyWithNoFrontends(t *testing.T) {
	rp := &reverseProxy{}
	server := httptest.NewServer(rp)
	defer server.Close()

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != 404 {
		t.Errorf("wrong status code: want (404) got (%d)\n", resp.StatusCode)
	}
}
