package config

import (
	"fmt"
	"log"
)

var (
	Version string
)

func init() {
	s := "Starting RPLB"
	if Version != "" {
		s += fmt.Sprintf(" %s", Version)
	}
	log.Println(s + " ...")
}
