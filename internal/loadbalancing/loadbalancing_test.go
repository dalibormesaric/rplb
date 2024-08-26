package loadbalancing

import (
	"reflect"
	"testing"
)

func TestNewAlgorithm(t *testing.T) {
	var tests = []struct {
		name     string
		expected string
	}{
		{
			name:     Sticky,
			expected: reflect.TypeOf(&sticky{}).String(),
		},
		{
			name:     RoundRobin,
			expected: reflect.TypeOf(&roundRobin{}).String(),
		},
		{
			name:     Random,
			expected: reflect.TypeOf(&random{}).String(),
		},
	}

	for _, test := range tests {
		algo := NewAlgorithm(test.name)
		got := reflect.TypeOf(algo).String()
		if got != test.expected {
			t.Errorf("wrong type for (%s): want (%s) got (%s)\n", test.name, test.expected, got)
		}
	}
}
