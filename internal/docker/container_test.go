package docker

import (
	"context"
	"fmt"
	"testing"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
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
		ContainerCreateFn: func(_ context.Context, _ client.ContainerCreateOptions) (client.ContainerCreateResult, error) {
			return client.ContainerCreateResult{}, fmt.Errorf("boom")
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
		ContainerStartFn: func(context.Context, string, client.ContainerStartOptions) (client.ContainerStartResult, error) {
			return client.ContainerStartResult{}, fmt.Errorf("start fail")
		},
		ContainerRemoveFn: func(context.Context, string, client.ContainerRemoveOptions) (client.ContainerRemoveResult, error) {
			removed = true
			return client.ContainerRemoveResult{}, nil
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
		ContainerWaitFn: func(context.Context, string, client.ContainerWaitOptions) client.ContainerWaitResult {
			ch := make(chan container.WaitResponse, 1)
			ch <- container.WaitResponse{StatusCode: 1}
			return client.ContainerWaitResult{Result: ch}
		},
	}
	if err := RunContainer(quiet(), mock, &Submission{Timeout: 5}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunContainer_WaitError(t *testing.T) {
	mock := &MockClient{
		ContainerWaitFn: func(context.Context, string, client.ContainerWaitOptions) client.ContainerWaitResult {
			errCh := make(chan error, 1)
			errCh <- fmt.Errorf("wait fail")
			return client.ContainerWaitResult{Error: errCh}
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
		ContainerRemoveFn: func(_ context.Context, _ string, opts client.ContainerRemoveOptions) (client.ContainerRemoveResult, error) {
			forced = opts.Force
			return client.ContainerRemoveResult{}, nil
		},
	}
	RunContainer(quiet(), mock, &Submission{Timeout: 5})
	if !forced {
		t.Error("expected force removal")
	}
}
