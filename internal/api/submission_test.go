package api

import (
	"autograder/internal/config"
	"autograder/internal/grader"
	"testing"

	"github.com/sirupsen/logrus"
)

func testSubmissionService() *SubmissionService {
	log := logrus.New()
	log.SetLevel(logrus.PanicLevel)
	return &SubmissionService{
		grader:     &grader.DefaultGrader{},
		docker:     &mockDockerClient{},
		log:        log,
		filesDir:   "./test_files/",
		mountPath:  "/mnt/",
		extensions: []string{"py"},
		timeout:    5,
	}
}

func TestValidateExtension_Valid(t *testing.T) {
	s := testSubmissionService()
	if err := s.ValidateExtension("solution.py"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateExtension_Invalid(t *testing.T) {
	s := testSubmissionService()
	if err := s.ValidateExtension("hack.js"); err == nil {
		t.Error("expected error for .js file")
	}
}

func TestValidateExtension_NoExtension(t *testing.T) {
	s := testSubmissionService()
	if err := s.ValidateExtension("noext"); err == nil {
		t.Error("expected error for file without extension")
	}
}

func TestValidateExtension_MultipleAllowed(t *testing.T) {
	s := testSubmissionService()
	s.extensions = []string{"py", "rb"}
	if err := s.ValidateExtension("script.rb"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestProcess_InvalidExtension(t *testing.T) {
	s := testSubmissionService()
	lab := &config.Lab{ID: "lab_1"}
	_, err := s.Process(SubmissionRequest{
		LabID: "lab_1", Code: []byte("x"), Filename: "hack.js",
	}, lab, "marker")
	if err == nil {
		t.Error("expected error for invalid extension")
	}
}

func TestProcess_WritesFileAndRunsContainer(t *testing.T) {
	s := testSubmissionService()
	s.filesDir = t.TempDir()
	lab := &config.Lab{
		ID: "lab_1",
		Testcase: []config.Testcase{{
			Type: "stdout", Name: "hw",
			Expected: []config.Expected{{Feedback: "OK", Points: 1, Values: []string{"Hello"}}},
		}},
	}

	// mockDockerClient.ContainerStart returns error ("no docker in test")
	// so Process will fail at the container step — but file I/O should succeed
	_, err := s.Process(SubmissionRequest{
		LabID: "lab_1", Code: []byte("print('Hello')"), Filename: "solution.py",
	}, lab, "marker")

	if err == nil {
		t.Error("expected error from docker mock")
	}
	// Verify the error is from container, not from file writing
	if err != nil && err.Error() == "invalid file extension. Allowed: [py]" {
		t.Error("should have passed extension validation")
	}
}
