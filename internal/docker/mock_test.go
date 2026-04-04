package docker

import (
	"context"
	"io"

	"github.com/docker/docker/api/types"
	containertypes "github.com/docker/docker/api/types/container"
	imagetypes "github.com/docker/docker/api/types/image"
)

// MockClient implements Client for testing.
type MockClient struct {
	ImageListFn       func(ctx context.Context, options imagetypes.ListOptions) ([]imagetypes.Summary, error)
	ImageRemoveFn     func(ctx context.Context, imageID string, options imagetypes.RemoveOptions) ([]imagetypes.DeleteResponse, error)
	ImageBuildFn      func(ctx context.Context, buildContext io.Reader, options types.ImageBuildOptions) (types.ImageBuildResponse, error)
	ContainerCreateFn func(ctx context.Context, config *containertypes.Config, hostConfig *containertypes.HostConfig, containerName string) (containertypes.CreateResponse, error)
	ContainerStartFn  func(ctx context.Context, containerID string, options containertypes.StartOptions) error
	ContainerWaitFn   func(ctx context.Context, containerID string, condition containertypes.WaitCondition) (<-chan containertypes.WaitResponse, <-chan error)
	ContainerLogsFn   func(ctx context.Context, containerID string, options containertypes.LogsOptions) (io.ReadCloser, error)
	ContainerRemoveFn func(ctx context.Context, containerID string, options containertypes.RemoveOptions) error
}

func (m *MockClient) ImageList(ctx context.Context, options imagetypes.ListOptions) ([]imagetypes.Summary, error) {
	if m.ImageListFn != nil {
		return m.ImageListFn(ctx, options)
	}
	return nil, nil
}

func (m *MockClient) ImageRemove(ctx context.Context, imageID string, options imagetypes.RemoveOptions) ([]imagetypes.DeleteResponse, error) {
	if m.ImageRemoveFn != nil {
		return m.ImageRemoveFn(ctx, imageID, options)
	}
	return nil, nil
}

func (m *MockClient) ImageBuild(ctx context.Context, buildContext io.Reader, options types.ImageBuildOptions) (types.ImageBuildResponse, error) {
	if m.ImageBuildFn != nil {
		return m.ImageBuildFn(ctx, buildContext, options)
	}
	return types.ImageBuildResponse{Body: io.NopCloser(io.Reader(nil))}, nil
}

func (m *MockClient) ContainerCreate(ctx context.Context, config *containertypes.Config, hostConfig *containertypes.HostConfig, containerName string) (containertypes.CreateResponse, error) {
	if m.ContainerCreateFn != nil {
		return m.ContainerCreateFn(ctx, config, hostConfig, containerName)
	}
	return containertypes.CreateResponse{ID: "mock-id"}, nil
}

func (m *MockClient) ContainerStart(ctx context.Context, containerID string, options containertypes.StartOptions) error {
	if m.ContainerStartFn != nil {
		return m.ContainerStartFn(ctx, containerID, options)
	}
	return nil
}

func (m *MockClient) ContainerWait(ctx context.Context, containerID string, condition containertypes.WaitCondition) (<-chan containertypes.WaitResponse, <-chan error) {
	if m.ContainerWaitFn != nil {
		return m.ContainerWaitFn(ctx, containerID, condition)
	}
	ch := make(chan containertypes.WaitResponse, 1)
	ch <- containertypes.WaitResponse{StatusCode: 0}
	return ch, nil
}

func (m *MockClient) ContainerLogs(ctx context.Context, containerID string, options containertypes.LogsOptions) (io.ReadCloser, error) {
	if m.ContainerLogsFn != nil {
		return m.ContainerLogsFn(ctx, containerID, options)
	}
	return io.NopCloser(io.Reader(nil)), nil
}

func (m *MockClient) ContainerRemove(ctx context.Context, containerID string, options containertypes.RemoveOptions) error {
	if m.ContainerRemoveFn != nil {
		return m.ContainerRemoveFn(ctx, containerID, options)
	}
	return nil
}
