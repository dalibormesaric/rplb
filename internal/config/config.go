package config

import (
	"fmt"
	"log/slog"
)

var (
	Version string
)

func init() {
	s := "Starting RPLB"
	if Version != "" {
		s += fmt.Sprintf(" %s", Version)
	}
	slog.Info(s + " ...")
}
