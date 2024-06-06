package frontend

import (
	"fmt"
	"strings"
)

type Host string

type Frontends map[Host]*Frontend

type Frontend struct {
	BackendName string
	Hits        int64
}

func CreateFrontends(urlNamePair string) (Frontends, error) {
	split := strings.Split(urlNamePair, ",")
	if len(split)%2 != 0 {
		return nil, fmt.Errorf("unable to split nameUrlPairs")
	}

	fe := make(Frontends)
	for i := range split {
		if (i+1)%2 == 0 {
			fe[Host(split[i-1])] = &Frontend{
				BackendName: split[i],
			}
		}
	}

	return fe, nil
}

func (f Frontends) Get(host string) *Frontend {
	return f[Host(host)]
}

func (f *Frontend) Inc() {
	f.Hits++
	fmt.Printf("%s %d\n", f.BackendName, f.Hits)
}
