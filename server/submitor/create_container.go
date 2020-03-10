package submitor

import (
	"context"
	"github.com/docker/docker/api/types"
	"io"
	"os"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

type Submission struct {
	ImageName     string
	ContainerName string
	BindedDir     string
	Command       []string
	TargetDir     string
}

type ContainerLog struct {
	ErrorType string
	Error     error
	Message   string
}

func CreateContainer(submission *Submission) {
	start := time.Now()
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatal(err)
	}

	res, err := cli.ContainerCreate(ctx,
		&container.Config{
			Image: submission.ImageName,
			Cmd:   submission.Command,
		},
		&container.HostConfig{
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeBind,
					Source: submission.BindedDir,
					Target: submission.TargetDir,
				},
			},
		}, nil, submission.ContainerName)

	if err != nil {
		log.Info(err)
	}

	err = cli.ContainerStart(ctx, res.ID, types.ContainerStartOptions{})

	reader, err := cli.ContainerLogs(ctx, res.ID, types.ContainerLogsOptions{})
	if err != nil {
		log.Info(err)
	}

	_, err = io.Copy(os.Stdout, reader)
	if err != nil && err != io.EOF {
		log.Info(err)
	}
	t := time.Now()

	log.Info(t.Sub(start))
}
