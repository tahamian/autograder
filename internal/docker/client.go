package docker

import (
	"context"
	"io"

	"github.com/moby/moby/client"
)

// Client abstracts the Docker API for testing.
type Client interface {
	ImageList(ctx context.Context, opts client.ImageListOptions) (client.ImageListResult, error)
	ImageRemove(ctx context.Context, id string, opts client.ImageRemoveOptions) (client.ImageRemoveResult, error)
	ImageBuild(ctx context.Context, buildContext io.Reader, opts client.ImageBuildOptions) (client.ImageBuildResult, error)
	ContainerCreate(ctx context.Context, opts client.ContainerCreateOptions) (client.ContainerCreateResult, error)
	ContainerStart(ctx context.Context, id string, opts client.ContainerStartOptions) (client.ContainerStartResult, error)
	ContainerWait(ctx context.Context, id string, opts client.ContainerWaitOptions) client.ContainerWaitResult
	ContainerRemove(ctx context.Context, id string, opts client.ContainerRemoveOptions) (client.ContainerRemoveResult, error)
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

func (r *RealClient) ImageList(ctx context.Context, opts client.ImageListOptions) (client.ImageListResult, error) {
	if err := r.ensureClient(); err != nil {
		return client.ImageListResult{}, err
	}
	return r.cli.ImageList(ctx, opts)
}

func (r *RealClient) ImageRemove(ctx context.Context, id string, opts client.ImageRemoveOptions) (client.ImageRemoveResult, error) {
	if err := r.ensureClient(); err != nil {
		return client.ImageRemoveResult{}, err
	}
	return r.cli.ImageRemove(ctx, id, opts)
}

func (r *RealClient) ImageBuild(ctx context.Context, buildContext io.Reader, opts client.ImageBuildOptions) (client.ImageBuildResult, error) {
	if err := r.ensureClient(); err != nil {
		return client.ImageBuildResult{}, err
	}
	return r.cli.ImageBuild(ctx, buildContext, opts)
}

func (r *RealClient) ContainerCreate(ctx context.Context, opts client.ContainerCreateOptions) (client.ContainerCreateResult, error) {
	if err := r.ensureClient(); err != nil {
		return client.ContainerCreateResult{}, err
	}
	return r.cli.ContainerCreate(ctx, opts)
}

func (r *RealClient) ContainerStart(ctx context.Context, id string, opts client.ContainerStartOptions) (client.ContainerStartResult, error) {
	if err := r.ensureClient(); err != nil {
		return client.ContainerStartResult{}, err
	}
	return r.cli.ContainerStart(ctx, id, opts)
}

func (r *RealClient) ContainerWait(ctx context.Context, id string, opts client.ContainerWaitOptions) client.ContainerWaitResult {
	if err := r.ensureClient(); err != nil {
		errCh := make(chan error, 1)
		errCh <- err
		return client.ContainerWaitResult{Error: errCh}
	}
	return r.cli.ContainerWait(ctx, id, opts)
}

func (r *RealClient) ContainerRemove(ctx context.Context, id string, opts client.ContainerRemoveOptions) (client.ContainerRemoveResult, error) {
	if err := r.ensureClient(); err != nil {
		return client.ContainerRemoveResult{}, err
	}
	return r.cli.ContainerRemove(ctx, id, opts)
}
