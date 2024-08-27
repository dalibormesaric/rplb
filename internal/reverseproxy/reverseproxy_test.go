package reverseproxy

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dalibormesaric/rplb/internal/backend"
	"github.com/dalibormesaric/rplb/internal/frontend"
	"github.com/dalibormesaric/rplb/internal/loadbalancing"
)

func TestReverseProxyWithNoFrontends(t *testing.T) {
	rp := &reverseProxy{}
	ts := httptest.NewServer(http.HandlerFunc(rp.reverseProxyAndLoadBalance))
	defer ts.Close()

	res, err := ts.Client().Get(ts.URL)
	if err != nil {
		t.Error(err)
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Error(err)
	}
	content404 := "<h1>404</h1>"
	if !strings.Contains(string(body), content404) {
		t.Errorf("wrong content: want to contain (%s) got (%s)\n", content404, body)
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
	algo, _ := loadbalancing.NewAlgorithm(loadbalancing.Random)
	rp := &reverseProxy{
		frontends:     f,
		loadbalancing: algo,
	}
	ts := httptest.NewServer(http.HandlerFunc(rp.reverseProxyAndLoadBalance))
	defer ts.Close()

	res, err := ts.Client().Get(ts.URL)
	if err != nil {
		t.Error(err)
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Error(err)
	}
	content503 := "<h1>503</h1>"
	if !strings.Contains(string(body), content503) {
		t.Errorf("wrong content: want to contain (%s) got (%s)\n", content503, body)
	}
	expectedStatusCode := 200
	if res.StatusCode != expectedStatusCode {
		t.Errorf("wrong status code: want (%d) got (%d)\n", expectedStatusCode, res.StatusCode)
	}
}

func TestReverseProxyWithFrontendsAndWithBackends(t *testing.T) {
	f, err := frontend.CreateFrontends("127.0.0.1,b")
	if err != nil {
		t.Error(err)
	}
	bts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Backend")
	}))
	defer bts.Close()
	b, err := backend.CreateBackends("b," + bts.URL)
	if err != nil {
		t.Error(err)
	}
	b["b"][0].SetLive(true)
	algo, _ := loadbalancing.NewAlgorithm(loadbalancing.Random)
	rp := &reverseProxy{
		frontends:     f,
		backends:      b,
		loadbalancing: algo,
	}
	ts := httptest.NewServer(http.HandlerFunc(rp.reverseProxyAndLoadBalance))
	defer ts.Close()

	res, err := ts.Client().Get(ts.URL)
	if err != nil {
		t.Error(err)
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Error(err)
	}
	contentBackend := "Backend"
	if !strings.Contains(string(body), contentBackend) {
		t.Errorf("wrong content: want to contain (%s) got (%s)\n", contentBackend, body)
	}
	expectedStatusCode := 200
	if res.StatusCode != expectedStatusCode {
		t.Errorf("wrong status code: want (%d) got (%d)\n", expectedStatusCode, res.StatusCode)
	}
}
