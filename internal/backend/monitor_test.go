package backend

import (
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func TestHealthCheckBackend(t *testing.T) {
	backend := httptest.NewServer(nil)

	url, err := url.Parse(backend.URL)
	if err != nil {
		t.Error(err)
	}
	latency := healthCheck(url.Host)
	if latency == 0 {
		t.Errorf("wrong latency: want (>0) got (%v)\n", latency)
	}
}

func TestHealthCheckNoBackend(t *testing.T) {
	latency := healthCheck("notexists.local:1234")
	if latency > 0 {
		t.Errorf("wrong latency: want (0) got (%v)\n", latency)
	}
}

func TestGetColorCodes(t *testing.T) {
	var tests = []struct {
		latencyString string
		colorCode     int64
	}{
		{"0µs", 0},
		{"79.892µs", 5},
		{"1.123ms", 5},
		{"1.123ms", 5},
		{"2.123ms", 5},
		{"3.123ms", 5},
		{"4.123ms", 5},
		{"5.123ms", 10},
		{"6.123ms", 10},
		{"7.123ms", 10},
		{"8.123ms", 10},
		{"9.123ms", 10},
		{"10.123ms", 10},
		{"11.123ms", 10},
		{"12.123ms", 10},
		{"22.123ms", 20},
		{"32.123ms", 30},
		{"92.123ms", 90},
		{"102.123ms", 100},
		{"112.123ms", 100},
		{"192.123ms", 100},
		{"202.123ms", 200},
		{"902.123ms", 900},
		{"1002.123ms", 1000},
		{"1012.123ms", 1000},
		{"1102.123ms", 10000},
		{"1502.123ms", 10000},
		{"5.123s", 10000},
	}

	for _, test := range tests {
		latency, err := time.ParseDuration(test.latencyString)
		if err != nil {
			t.Error(err)
		}
		colorCode := getColorCode(latency)
		if colorCode != test.colorCode {
			t.Errorf("wrong colorCode: want (%d) got (%d)\n", test.colorCode, colorCode)
		}
	}
}

func TestLast20(t *testing.T) {
	var tests = []struct {
		monitorFrames []MonitorFrame
		expectedLen   int
	}{
		{nil, 0},
		{[]MonitorFrame{}, 0},
		{[]MonitorFrame{{}}, 1},
		{[]MonitorFrame{{}, {}}, 2},
		{[]MonitorFrame{{}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}}, 19},
		{[]MonitorFrame{{}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}}, 20},
		{[]MonitorFrame{{}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}}, 20},
	}

	for _, test := range tests {
		last20 := last20(test.monitorFrames)
		if len(last20) != test.expectedLen {
			t.Errorf("wrong monitor frames length: want (%d) got (%d)\n", test.expectedLen, len(last20))
		}
	}
}
