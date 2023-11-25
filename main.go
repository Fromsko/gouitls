package main

import (
	"log"

	"github.com/Fromsko/gouitls/docker"
)

func main() {
	d, err := docker.NewDockerDash("alpine")
	if err != nil {
		return
	}
	dash := &docker.Dash{
		Command:      []string{"sh", "-c", "while true; do echo hello; sleep 1; done"},
		DockerDasher: d,
	}

	err = docker.PullAndRunAlpine(dash)
	if err != nil {
		log.Fatal(err)
		return
	}
}
