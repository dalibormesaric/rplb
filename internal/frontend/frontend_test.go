package frontend

import (
	"fmt"
	"testing"
)

func TestCreateFrontends(t *testing.T) {
	var tests = []struct {
		urlNamePair string
		key         string
		values      *Frontend
	}{
		{"a,b", "a", &Frontend{BackendName: "b"}},
	}

	for _, test := range tests {
		result, _ := NewFrontends(test.urlNamePair)
		f := result.Get(test.key)
		if f.BackendName != test.values.BackendName {
			t.Errorf("wrong backend name: want (%s) got (%s)\n", test.values.BackendName, f.BackendName)
		}
	}
}

func TestCreateFrontendsErrors(t *testing.T) {
	var tests = []struct {
		urlNamePair string
		err         error
	}{
		{"a", fmt.Errorf("frontends must be a comma-separated list containing even number of items")},
		{"a,b,c", fmt.Errorf("frontends must be a comma-separated list containing even number of items")},
		{",", fmt.Errorf("urlNamePair at index 0 must have a value")},
		{"a,", fmt.Errorf("urlNamePair at index 1 must have a value")},
		{"a,b,a,b", fmt.Errorf("frontend host has to be unique")},
		{"", nil},
		{" ", nil},
		{"a,b", nil},
	}

	for _, test := range tests {
		_, err := NewFrontends(test.urlNamePair)
		if test.err != nil && err.Error() != test.err.Error() {
			t.Errorf("Want: %s\nGot: %s\n", test.err.Error(), err.Error())
		}
		if test.err == nil && err != nil {
			t.Errorf("Was not expecting an error\n")
		}
	}
}
