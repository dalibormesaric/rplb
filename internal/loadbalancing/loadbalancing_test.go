package loadbalancing

import (
	"fmt"
	"reflect"
	"testing"
)

func TestNewAlgorithm(t *testing.T) {
	var tests = []struct {
		name     string
		expected string
		err      error
	}{
		{
			name:     Sticky,
			expected: reflect.TypeOf(&sticky{}).String(),
			err:      nil,
		},
		{
			name:     RoundRobin,
			expected: reflect.TypeOf(&roundRobin{}).String(),
			err:      nil,
		},
		{
			name:     First,
			expected: reflect.TypeOf(&first{}).String(),
			err:      nil,
		},
		{
			name:     Random,
			expected: reflect.TypeOf(&random{}).String(),
			err:      nil,
		},
		{
			name:     "foo",
			expected: "",
			err:      fmt.Errorf("unknown algorithm type (foo)"),
		},
	}

	for _, test := range tests {
		algo, err := NewAlgorithm(test.name)
		if err != nil && err.Error() != test.err.Error() {
			t.Errorf("wrong error for (%s): want (%s) got (%s)\n", test.name, test.err.Error(), err.Error())
		}
		if err == nil {
			got := reflect.TypeOf(algo).String()
			if got != test.expected {
				t.Errorf("wrong type for (%s): want (%s) got (%s)\n", test.name, test.expected, got)
			}
		}
	}
}
