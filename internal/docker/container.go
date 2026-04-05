package docker

import (
	"context"
	"fmt"
	"time"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/mount"
	"github.com/moby/moby/client"
	"github.com/sirupsen/logrus"
)

// Submission describes a grading container to run.
type Submission struct {
	ImageName     string
	ContainerName string
	BindedDir     string
	TargetDir     string
	Timeout       int // seconds
}

// RunContainer creates, starts, waits for, and cleans up a grading container.
func RunContainer(log *logrus.Logger, cli Client, sub *Submission) error {
	start := time.Now()

	timeout := time.Duration(sub.Timeout) * time.Second
	if timeout <= 0 {
		timeout = 10 * time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	resp, err := cli.ContainerCreate(ctx, client.ContainerCreateOptions{
		Name: sub.ContainerName,
		Config: &container.Config{
			Image:        sub.ImageName,
			Tty:          false,
			AttachStderr: true,
			AttachStdout: true,
		},
		HostConfig: &container.HostConfig{
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeBind,
					Source: sub.BindedDir,
					Target: sub.TargetDir,
				},
			},
			Resources: container.Resources{
				Memory:   256 * 1024 * 1024,
				NanoCPUs: 1e9,
			},
			NetworkMode: "none",
		},
	})
	if err != nil {
		return fmt.Errorf("creating container: %w", err)
	}

	containerID := resp.ID
	defer func() {
		rmCtx, rmCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer rmCancel()
		if _, err := cli.ContainerRemove(rmCtx, containerID, client.ContainerRemoveOptions{Force: true}); err != nil {
			log.WithError(err).Warn("failed to remove container")
		}
	}()

	if _, err := cli.ContainerStart(ctx, containerID, client.ContainerStartOptions{}); err != nil {
		return fmt.Errorf("starting container: %w", err)
	}

	waitResult := cli.ContainerWait(ctx, containerID, client.ContainerWaitOptions{
		Condition: container.WaitConditionNotRunning,
	})

	select {
	case err := <-waitResult.Error:
		if err != nil {
			return fmt.Errorf("container wait error: %w", err)
		}
	case status := <-waitResult.Result:
		if status.StatusCode != 0 {
			log.WithField("exit_code", status.StatusCode).Warn("container exited with non-zero status")
		}
	case <-ctx.Done():
		return fmt.Errorf("container timed out after %v", timeout)
	}

	log.WithField("duration", time.Since(start)).Info("container finished")
	return nil
}
