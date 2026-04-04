package docker

import (
	"context"
	"fmt"
	"time"

	containertypes "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
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

	resp, err := cli.ContainerCreate(ctx,
		&containertypes.Config{
			Image:        sub.ImageName,
			Tty:          false,
			AttachStderr: true,
			AttachStdout: true,
		},
		&containertypes.HostConfig{
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeBind,
					Source: sub.BindedDir,
					Target: sub.TargetDir,
				},
			},
			Resources: containertypes.Resources{
				Memory:   256 * 1024 * 1024, // 256 MB
				NanoCPUs: 1e9,               // 1 CPU
			},
			NetworkMode: "none",
		},
		sub.ContainerName,
	)
	if err != nil {
		return fmt.Errorf("creating container: %w", err)
	}

	containerID := resp.ID
	defer func() {
		rmCtx, rmCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer rmCancel()
		if err := cli.ContainerRemove(rmCtx, containerID, containertypes.RemoveOptions{Force: true}); err != nil {
			log.WithError(err).Warn("failed to remove container")
		}
	}()

	if err := cli.ContainerStart(ctx, containerID, containertypes.StartOptions{}); err != nil {
		return fmt.Errorf("starting container: %w", err)
	}

	statusCh, errCh := cli.ContainerWait(ctx, containerID, containertypes.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return fmt.Errorf("waiting for container: %w", err)
		}
	case status := <-statusCh:
		if status.StatusCode != 0 {
			log.WithField("exit_code", status.StatusCode).Warn("container exited with non-zero status")
		}
	case <-ctx.Done():
		return fmt.Errorf("container timed out after %v", timeout)
	}

	log.WithField("duration", time.Since(start)).Info("container finished")
	return nil
}
