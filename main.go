package main

import (
	"github.com/atkrad/gateway/cmd"
	log "github.com/sirupsen/logrus"
	"os"
)

func init() {
	log.SetOutput(os.Stdout)
}

func main() {
	cmd.Execute()
}
