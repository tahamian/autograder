package docker

import (
	"context"
	"io"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
)

// MockClient implements Client for testing.
type MockClient struct {
	ImageListFn       func(ctx context.Context, opts client.ImageListOptions) (client.ImageListResult, error)
	ImageRemoveFn     func(ctx context.Context, id string, opts client.ImageRemoveOptions) (client.ImageRemoveResult, error)
	ImageBuildFn      func(ctx context.Context, buildContext io.Reader, opts client.ImageBuildOptions) (client.ImageBuildResult, error)
	ContainerCreateFn func(ctx context.Context, opts client.ContainerCreateOptions) (client.ContainerCreateResult, error)
	ContainerStartFn  func(ctx context.Context, id string, opts client.ContainerStartOptions) (client.ContainerStartResult, error)
	ContainerWaitFn   func(ctx context.Context, id string, opts client.ContainerWaitOptions) client.ContainerWaitResult
	ContainerRemoveFn func(ctx context.Context, id string, opts client.ContainerRemoveOptions) (client.ContainerRemoveResult, error)
}

func (m *MockClient) ImageList(ctx context.Context, opts client.ImageListOptions) (client.ImageListResult, error) {
	if m.ImageListFn != nil {
		return m.ImageListFn(ctx, opts)
	}
	return client.ImageListResult{}, nil
}

func (m *MockClient) ImageRemove(ctx context.Context, id string, opts client.ImageRemoveOptions) (client.ImageRemoveResult, error) {
	if m.ImageRemoveFn != nil {
		return m.ImageRemoveFn(ctx, id, opts)
	}
	return client.ImageRemoveResult{}, nil
}

func (m *MockClient) ImageBuild(ctx context.Context, buildContext io.Reader, opts client.ImageBuildOptions) (client.ImageBuildResult, error) {
	if m.ImageBuildFn != nil {
		return m.ImageBuildFn(ctx, buildContext, opts)
	}
	return client.ImageBuildResult{Body: io.NopCloser(io.Reader(nil))}, nil
}

func (m *MockClient) ContainerCreate(ctx context.Context, opts client.ContainerCreateOptions) (client.ContainerCreateResult, error) {
	if m.ContainerCreateFn != nil {
		return m.ContainerCreateFn(ctx, opts)
	}
	return client.ContainerCreateResult{ID: "mock-id"}, nil
}

func (m *MockClient) ContainerStart(ctx context.Context, id string, opts client.ContainerStartOptions) (client.ContainerStartResult, error) {
	if m.ContainerStartFn != nil {
		return m.ContainerStartFn(ctx, id, opts)
	}
	return client.ContainerStartResult{}, nil
}

func (m *MockClient) ContainerWait(ctx context.Context, id string, opts client.ContainerWaitOptions) client.ContainerWaitResult {
	if m.ContainerWaitFn != nil {
		return m.ContainerWaitFn(ctx, id, opts)
	}
	ch := make(chan container.WaitResponse, 1)
	ch <- container.WaitResponse{StatusCode: 0}
	return client.ContainerWaitResult{Result: ch}
}

func (m *MockClient) ContainerRemove(ctx context.Context, id string, opts client.ContainerRemoveOptions) (client.ContainerRemoveResult, error) {
	if m.ContainerRemoveFn != nil {
		return m.ContainerRemoveFn(ctx, id, opts)
	}
	return client.ContainerRemoveResult{}, nil
}
