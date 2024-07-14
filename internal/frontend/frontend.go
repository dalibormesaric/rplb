package frontend

import (
	"fmt"
	"log"
	"strings"
)

type Host string

type Frontends map[Host]*Frontend

type Frontend struct {
	BackendName string
	Hits        int64
}

func CreateFrontends(urlNamePair string) (Frontends, error) {
	fe := make(Frontends)

	if strings.TrimSpace(urlNamePair) == "" {
		log.Println("No frontends configured")
		return fe, nil
	}

	split := strings.Split(urlNamePair, ",")
	if len(split)%2 != 0 {
		return nil, fmt.Errorf("frontends must be a comma-separated list containing even number of items")
	}

	for i, v := range split {
		if v == "" {
			return nil, fmt.Errorf("urlNamePair at index %d must have a value", i)
		}

		if (i+1)%2 == 0 {
			host := Host(split[i-1])

			_, ok := fe[host]
			if ok {
				return nil, fmt.Errorf("frontend host has to be unique")
			}

			backendName := split[i]
			fe[host] = &Frontend{
				BackendName: backendName,
			}
			log.Printf("Added frontend host (%s) for (%s)\n", host, backendName)
		}
	}

	return fe, nil
}

func (f Frontends) Get(host string) *Frontend {
	return f[Host(host)]
}

func (f *Frontend) Inc() int64 {
	f.Hits++
	return f.Hits
}
