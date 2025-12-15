package main

import (
	"os"

	"github.com/lnobach/gonrg/cli"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetOutput(os.Stderr)
	log.SetLevel(log.WarnLevel)
	err := cli.Start()
	if err != nil {
		panic(err)
	}
}
