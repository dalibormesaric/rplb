package reverseproxy

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dalibormesaric/rplb/internal/frontend"
)

func TestReverseProxyWithNoFrontends(t *testing.T) {
	rp := &reverseProxy{}
	ts := httptest.NewServer(http.HandlerFunc(rp.reverseProxyAndLoadBalance))
	defer ts.Close()

	res, err := ts.Client().Get(ts.URL)
	if err != nil {
		t.Error(err)
	}
	b, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Error(err)
	}
	content404 := "<h1>404</h1>"
	if !strings.Contains(string(b), content404) {
		t.Errorf("wrong content: want to contain (%s) got (%s)\n", content404, b)
	}
	expectedStatusCode := 200
	if res.StatusCode != expectedStatusCode {
		t.Errorf("wrong status code: want (%d) got (%d)\n", expectedStatusCode, res.StatusCode)
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
	ts := httptest.NewServer(http.HandlerFunc(rp.reverseProxyAndLoadBalance))
	defer ts.Close()

	res, err := ts.Client().Get(ts.URL)
	if err != nil {
		t.Error(err)
	}
	b, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Error(err)
	}
	content503 := "<h1>503</h1>"
	if !strings.Contains(string(b), content503) {
		t.Errorf("wrong content: want to contain (%s) got (%s)\n", content503, b)
	}
	expectedStatusCode := 200
	if res.StatusCode != expectedStatusCode {
		t.Errorf("wrong status code: want (%d) got (%d)\n", expectedStatusCode, res.StatusCode)
	}
}
