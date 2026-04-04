package main

import (
	"autograder/internal/api"
	"autograder/internal/config"
	"autograder/internal/docker"
	"autograder/internal/grader"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()
	log.Out = os.Stdout
	log.SetFormatter(&logrus.JSONFormatter{
		DisableTimestamp: false,
	})

	cfg, err := config.Load("config.yml")
	if err != nil {
		log.WithError(err).Fatal("failed to load config")
	}

	dockerClient := &docker.RealClient{}
	if err := docker.EnsureImage(log, dockerClient, cfg.Marker.ImageName); err != nil {
		log.WithError(err).Fatal("failed to ensure marker image")
	}

	srv, err := api.NewServer(cfg, log, dockerClient, &grader.DefaultGrader{})
	if err != nil {
		log.WithError(err).Fatal("failed to start server")
	}
	defer srv.Stop()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Info("shutting down")
}
