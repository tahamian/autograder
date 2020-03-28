package submitor

import (
	"context"
	"os"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"

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
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatal(err)
	}

	resp, err := cli.ContainerCreate(ctx,
		&container.Config{
			Image: submission.ImageName,
			Cmd:   submission.Command,
		},
		&container.HostConfig{
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeBind,
					Source: submission.TargetDir,
					Target: submission.BindedDir,
				},
			},
		}, nil, submission.ContainerName)

	if err != nil {
		log.Info(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			panic(err)
		}
	case <-statusCh:
	}

	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		panic(err)
	}

	_, err = stdcopy.StdCopy(os.Stdout, os.Stderr, out)

	if err != nil {
		log.Warn("error while ")
	}

	//err = cli.ContainerStart(ctx, res.ID, types.ContainerStartOptions{})
	//if err != nil {
	//	log.Info(err)
	//}

	//reader, err := cli.ContainerLogs(ctx, res.ID, types.ContainerLogsOptions{})
	//if err != nil {
	//	log.Info(err)
	//}
	//
	//_, err = io.Copy(os.Stdout, reader)
	//if err != nil && err != io.EOF {
	//	log.Info(err)
	//}

	t := time.Now()
	log.Info(t.Sub(start))
}
