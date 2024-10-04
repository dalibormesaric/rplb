package integration

import (
	"log"
	"os/exec"
	"strings"
)

func setUp(algo string) {
	out, err := exec.Command("docker", "-v").Output()
	if err != nil {
		log.Fatalf("Docker not running: (%v)", err)
	}
	if !strings.HasPrefix(string(out[:]), "Docker version") {
		log.Fatal("Docker not running")
	}

	err = exec.Command("docker", "compose", "--env-file", "."+algo+".env", "-f", "../example/compose.yaml", "up", "-d", "rplb", "--build").Run()
	if err != nil {
		log.Fatal(err)
	}
}

func tearDown() {
	err := exec.Command("docker", "compose", "-f", "../example/compose.yaml", "down").Run()
	if err != nil {
		log.Fatal(err)
	}
}
