package main

import (
	//"fmt"
	"autograder/server"
	"os"
	"os/signal"
	log "github.com/sirupsen/logrus"
)



func main() {
	htmlServer := server.StartServer("config.yaml")
	defer htmlServer.Stop()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan

	log.Info("Stopping : autograder")

}
