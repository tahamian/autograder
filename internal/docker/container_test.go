package docker

import (
	"context"
	"fmt"
	"testing"

	containertypes "github.com/docker/docker/api/types/container"
	"github.com/sirupsen/logrus"
)

func quiet() *logrus.Logger {
	l := logrus.New()
	l.SetLevel(logrus.PanicLevel)
	return l
}

func TestRunContainer_Success(t *testing.T) {
	err := RunContainer(quiet(), &MockClient{}, &Submission{
		ImageName: "img", ContainerName: "c", BindedDir: "/tmp", TargetDir: "/mnt", Timeout: 5,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunContainer_CreateError(t *testing.T) {
	mock := &MockClient{
		ContainerCreateFn: func(ctx context.Context, config *containertypes.Config, hc *containertypes.HostConfig, name string) (containertypes.CreateResponse, error) {
			return containertypes.CreateResponse{}, fmt.Errorf("boom")
		},
	}
	err := RunContainer(quiet(), mock, &Submission{Timeout: 5})
	if err == nil || err.Error() != "creating container: boom" {
		t.Fatalf("expected create error, got: %v", err)
	}
}

func TestRunContainer_StartError_CleansUp(t *testing.T) {
	removed := false
	mock := &MockClient{
		ContainerStartFn: func(context.Context, string, containertypes.StartOptions) error {
			return fmt.Errorf("start fail")
		},
		ContainerRemoveFn: func(context.Context, string, containertypes.RemoveOptions) error {
			removed = true
			return nil
		},
	}
	err := RunContainer(quiet(), mock, &Submission{Timeout: 5})
	if err == nil {
		t.Fatal("expected error")
	}
	if !removed {
		t.Error("expected container removal after start failure")
	}
}

func TestRunContainer_NonZeroExit_NoError(t *testing.T) {
	mock := &MockClient{
		ContainerWaitFn: func(context.Context, string, containertypes.WaitCondition) (<-chan containertypes.WaitResponse, <-chan error) {
			ch := make(chan containertypes.WaitResponse, 1)
			ch <- containertypes.WaitResponse{StatusCode: 1}
			return ch, nil
		},
	}
	if err := RunContainer(quiet(), mock, &Submission{Timeout: 5}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunContainer_WaitError(t *testing.T) {
	mock := &MockClient{
		ContainerWaitFn: func(context.Context, string, containertypes.WaitCondition) (<-chan containertypes.WaitResponse, <-chan error) {
			ch := make(chan error, 1)
			ch <- fmt.Errorf("wait fail")
			return nil, ch
		},
	}
	err := RunContainer(quiet(), mock, &Submission{Timeout: 5})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestRunContainer_DefaultTimeout(t *testing.T) {
	if err := RunContainer(quiet(), &MockClient{}, &Submission{Timeout: 0}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunContainer_CleanupForceRemove(t *testing.T) {
	forced := false
	mock := &MockClient{
		ContainerRemoveFn: func(_ context.Context, _ string, opts containertypes.RemoveOptions) error {
			forced = opts.Force
			return nil
		},
	}
	RunContainer(quiet(), mock, &Submission{Timeout: 5})
	if !forced {
		t.Error("expected force removal")
	}
}
