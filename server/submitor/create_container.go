package submitor

import (
	"context"
	"fmt"
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
			//Cmd:       submission.Command,
			Tty:          true,
			OpenStdin:    true,
			StdinOnce:    true,
			AttachStderr: true,
			AttachStdin:  true,
			AttachStdout: true,
		},
		&container.HostConfig{

			//AutoRemove: true,
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

	result, err := stdcopy.StdCopy(os.Stdout, os.Stderr, out)

	if err != nil {
		log.Warn("error while ")
	}
	fmt.Println(result)
	t := time.Now()
	log.Info(t.Sub(start))
}
