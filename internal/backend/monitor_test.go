package backend

import (
	"testing"
	"time"
)

func TestGetColorCodes(t *testing.T) {
	var tests = []struct {
		latencyString string
		colorCode     int64
	}{
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
		latency, _ := time.ParseDuration(test.latencyString)
		colorCode := getColorCode(latency)
		if colorCode != test.colorCode {
			t.Errorf("wrong colorCode: want (%d) got (%d)\n", test.colorCode, colorCode)
		}
	}
}
