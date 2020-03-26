package main

import (
	//"fmt"
	"autograder/server"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
)

var log = logrus.New()

func main() {
	log.Out = os.Stdout
	log.SetFormatter(&logrus.JSONFormatter{
		DisableTimestamp: false,
	})
	htmlServer := server.StartServer("config.yaml")
	defer htmlServer.Stop()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan

	log.Info("HTTP server shutdown")

}
