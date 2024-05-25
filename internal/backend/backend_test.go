package backend

import (
	"fmt"
	"testing"
)

func TestGetUrlsForNames(t *testing.T) {
	var tests = []struct {
		nameUrlPairs string
		key          string
		values       []string
	}{
		{"a,b", "a", []string{"b"}},
	}

	for _, test := range tests {
		result, _ := GetUrlsForNames(test.nameUrlPairs)
		for i, v := range result[test.key] {
			if v != test.values[i] {
				t.Errorf("error")
			}
		}
	}
}

func TestGetUrlsForNamesErrors(t *testing.T) {
	var tests = []struct {
		nameUrlPairs string
		err          error
	}{
		{"", fmt.Errorf("unable to split nameUrlPairs")},
		{"a", fmt.Errorf("unable to split nameUrlPairs")},
		{"a,b,c", fmt.Errorf("unable to split nameUrlPairs")},
		{",", fmt.Errorf("nameUrlPair at index 0 must have a value")},
		{"a,", fmt.Errorf("nameUrlPair at index 1 must have a value")},
		{"a,b", nil},
	}

	for _, test := range tests {
		_, err := GetUrlsForNames(test.nameUrlPairs)
		if test.err != nil && err.Error() != test.err.Error() {
			t.Errorf("Want: %s\nGot: %s\n", test.err.Error(), err.Error())
		}
		if test.err == nil && err != nil {
			t.Errorf("Was not expecting an error")
		}
	}
}
