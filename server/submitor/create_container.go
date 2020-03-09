package submitor

import (
	"context"
	//"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	log "github.com/sirupsen/logrus"
)

type Submission struct {
	ImageName     string
	ContainerName string
	BindedDir     string
	Command       []string
}

type ContainerLog struct {
	ErrorType string
	Error     error
	Message   string
}

func CreateContainer(submission *Submission) {
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
					Target: "/input",
				},
			},
		}, nil, submission.ContainerName)

	if err != nil {
		//error_log := &ContainerLog{Error: err, ErrorType: err.Error(), Message: "failed to create container"}
		//log.WithField(error_log).Info()
	}

	fmt.Println(res)

}
