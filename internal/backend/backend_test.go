package backend

import (
	"fmt"
	"testing"
)

func TestCreateBackends(t *testing.T) {
	var tests = []struct {
		nameUrlPairs string
		key          string
		host         string
	}{
		{"a,http://b:1234", "a", "b:1234"},
	}

	for _, test := range tests {
		backends, err := CreateBackends(test.nameUrlPairs)
		if err != nil {
			t.Error(err)
		}
		b := backends[test.key][0]
		if b.URL.Host != test.host {
			t.Errorf("wrong backend url: want (%s) got (%s)\n", test.host, b.URL.Host)
		}
	}
}

func TestCreateBackendsErrors(t *testing.T) {
	var tests = []struct {
		nameUrlPairs string
		err          error
	}{
		{"a", fmt.Errorf("backends must be a comma-separated list containing even number of items")},
		{"a,b,c", fmt.Errorf("backends must be a comma-separated list containing even number of items")},
		{",", fmt.Errorf("nameUrlPair at index 0 must have a value")},
		{"a,", fmt.Errorf("nameUrlPair at index 1 must have a value")},
		{"a,b,a,b", fmt.Errorf("empty host for url (b) in backend (a)")},
		{"a,http://b:1234,a,http://b:1234", fmt.Errorf("url (http://b:1234) already exist in backend (a)")},
		{"", nil},
		{" ", nil},
		{"a,http://b:1234", nil},
	}

	for _, test := range tests {
		_, err := CreateBackends(test.nameUrlPairs)
		if test.err != nil && err.Error() != test.err.Error() {
			t.Errorf("was expecting an error: want (%s)\ngot: (%s)\n", test.err.Error(), err.Error())
		}
		if test.err == nil && err != nil {
			t.Errorf("was not expecting an error\n")
		}
	}
}

func TestStripString(t *testing.T) {
	var tests = []struct {
		input          string
		expectedOutput string
	}{
		{"", ""},
		{" ", ""},
		{"%-_=#\"", ""},
		{"abc", "abc"},
		{"ABC", "ABC"},
		{"123", "123"},
		{"a,http://b:1234", "ahttpb1234"},
	}

	for _, test := range tests {
		output := stripString(test.input)
		if output != test.expectedOutput {
			t.Errorf("wrong stripString output: want (%s) got (%s)\n", test.expectedOutput, output)
		}
	}
}
