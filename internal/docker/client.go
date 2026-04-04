package docker

import (
	"context"
	"io"

	"github.com/docker/docker/api/types"
	containertypes "github.com/docker/docker/api/types/container"
	imagetypes "github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

// Client abstracts the Docker API for testing.
type Client interface {
	ImageList(ctx context.Context, options imagetypes.ListOptions) ([]imagetypes.Summary, error)
	ImageRemove(ctx context.Context, imageID string, options imagetypes.RemoveOptions) ([]imagetypes.DeleteResponse, error)
	ImageBuild(ctx context.Context, buildContext io.Reader, options types.ImageBuildOptions) (types.ImageBuildResponse, error)
	ContainerCreate(ctx context.Context, config *containertypes.Config, hostConfig *containertypes.HostConfig, containerName string) (containertypes.CreateResponse, error)
	ContainerStart(ctx context.Context, containerID string, options containertypes.StartOptions) error
	ContainerWait(ctx context.Context, containerID string, condition containertypes.WaitCondition) (<-chan containertypes.WaitResponse, <-chan error)
	ContainerLogs(ctx context.Context, containerID string, options containertypes.LogsOptions) (io.ReadCloser, error)
	ContainerRemove(ctx context.Context, containerID string, options containertypes.RemoveOptions) error
}

// RealClient wraps the Docker SDK client.
type RealClient struct {
	cli *client.Client
}

func (r *RealClient) ensureClient() error {
	if r.cli == nil {
		c, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			return err
		}
		r.cli = c
	}
	return nil
}

func (r *RealClient) ImageList(ctx context.Context, options imagetypes.ListOptions) ([]imagetypes.Summary, error) {
	if err := r.ensureClient(); err != nil {
		return nil, err
	}
	return r.cli.ImageList(ctx, options)
}

func (r *RealClient) ImageRemove(ctx context.Context, imageID string, options imagetypes.RemoveOptions) ([]imagetypes.DeleteResponse, error) {
	if err := r.ensureClient(); err != nil {
		return nil, err
	}
	return r.cli.ImageRemove(ctx, imageID, options)
}

func (r *RealClient) ImageBuild(ctx context.Context, buildContext io.Reader, options types.ImageBuildOptions) (types.ImageBuildResponse, error) {
	if err := r.ensureClient(); err != nil {
		return types.ImageBuildResponse{}, err
	}
	return r.cli.ImageBuild(ctx, buildContext, options)
}

func (r *RealClient) ContainerCreate(ctx context.Context, config *containertypes.Config, hostConfig *containertypes.HostConfig, containerName string) (containertypes.CreateResponse, error) {
	if err := r.ensureClient(); err != nil {
		return containertypes.CreateResponse{}, err
	}
	return r.cli.ContainerCreate(ctx, config, hostConfig, nil, nil, containerName)
}

func (r *RealClient) ContainerStart(ctx context.Context, containerID string, options containertypes.StartOptions) error {
	if err := r.ensureClient(); err != nil {
		return err
	}
	return r.cli.ContainerStart(ctx, containerID, options)
}

func (r *RealClient) ContainerWait(ctx context.Context, containerID string, condition containertypes.WaitCondition) (<-chan containertypes.WaitResponse, <-chan error) {
	if r.ensureClient() != nil {
		errCh := make(chan error, 1)
		errCh <- r.ensureClient()
		return nil, errCh
	}
	return r.cli.ContainerWait(ctx, containerID, condition)
}

func (r *RealClient) ContainerLogs(ctx context.Context, containerID string, options containertypes.LogsOptions) (io.ReadCloser, error) {
	if err := r.ensureClient(); err != nil {
		return nil, err
	}
	return r.cli.ContainerLogs(ctx, containerID, options)
}

func (r *RealClient) ContainerRemove(ctx context.Context, containerID string, options containertypes.RemoveOptions) error {
	if err := r.ensureClient(); err != nil {
		return err
	}
	return r.cli.ContainerRemove(ctx, containerID, options)
}
