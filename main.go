package main

import (
	"fmt"
	"autograder/server"
	"os"
	"os/signal"
)



func main() {
	htmlServer := server.StartServer("config.yaml")
	defer htmlServer.Stop()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan

	fmt.Println("main : shutting down")

}
