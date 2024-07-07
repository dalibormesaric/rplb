package backend

import (
	"fmt"
	"testing"
)

func TestCreateBackends(t *testing.T) {
	var tests = []struct {
		nameUrlPairs string
		key          string
		values       []Backend
	}{
		{"a,b", "a", []Backend{{Url: "b"}}},
	}

	for _, test := range tests {
		result, _ := CreateBackends(test.nameUrlPairs)
		for i, v := range result[test.key] {
			if v.Url != test.values[i].Url {
				t.Errorf("wrong backend url: want (%s) got (%s)", test.values[i].Url, v.Url)
			}
		}
	}
}

func TestCreateBackendsErrors(t *testing.T) {
	var tests = []struct {
		nameUrlPairs string
		err          error
	}{
		{"", fmt.Errorf("backends must be a comma-separated list containing even number of items")},
		{"a", fmt.Errorf("backends must be a comma-separated list containing even number of items")},
		{"a,b,c", fmt.Errorf("backends must be a comma-separated list containing even number of items")},
		{",", fmt.Errorf("nameUrlPair at index 0 must have a value")},
		{"a,", fmt.Errorf("nameUrlPair at index 1 must have a value")},
		{"a,b,a,b", fmt.Errorf("url (b) already exist in backend (a)")},
		{"a,b", nil},
	}

	for _, test := range tests {
		_, err := CreateBackends(test.nameUrlPairs)
		if test.err != nil && err.Error() != test.err.Error() {
			t.Errorf("Want: %s\nGot: %s\n", test.err.Error(), err.Error())
		}
		if test.err == nil && err != nil {
			t.Errorf("Was not expecting an error")
		}
	}
}
