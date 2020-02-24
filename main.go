package main

import (
	//"fmt"
	"autograder/server"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
)

func main() {
	htmlServer := server.StartServer("config.yaml")
	defer htmlServer.Stop()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan

	log.Info("Stopping : autograder")

}
