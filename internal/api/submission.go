package api

import (
	"autograder/internal/config"
	"autograder/internal/docker"
	"autograder/internal/grader"
	"autograder/internal/models"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// ValidationError represents a client-side input error.
type ValidationError struct {
	Msg string
}

func (e *ValidationError) Error() string { return e.Msg }

// SubmissionRequest represents a parsed submission from either file or code.
type SubmissionRequest struct {
	LabID    string
	Code     []byte
	Filename string
}

// SubmissionService handles the grading pipeline: file I/O → container → evaluation.
type SubmissionService struct {
	grader     grader.Grader
	docker     docker.Client
	log        *logrus.Logger
	filesDir   string
	mountPath  string
	extensions []string
	timeout    int
	hostDir    string // HostSubmissionsFolder for Docker bind mounts
}

// NewSubmissionService creates a SubmissionService from config.
func NewSubmissionService(cfg *config.Config, log *logrus.Logger, dc docker.Client, g grader.Grader) *SubmissionService {
	exts := cfg.Marker.AllowedExtensions
	if len(exts) == 0 {
		exts = []string{"py"}
	}
	timeout := cfg.Marker.ContainerTimeout
	if timeout <= 0 {
		timeout = 10
	}
	return &SubmissionService{
		grader:     g,
		docker:     dc,
		log:        log,
		filesDir:   cfg.Marker.SubmissionsFolder,
		mountPath:  cfg.Marker.MountPath,
		extensions: exts,
		timeout:    timeout,
		hostDir:    cfg.Marker.HostSubmissionsFolder,
	}
}

// ValidateExtension checks if the filename has an allowed extension.
func (s *SubmissionService) ValidateExtension(filename string) error {
	parts := strings.Split(filename, ".")
	if len(parts) < 2 {
		return &ValidationError{fmt.Sprintf("Invalid file extension. Allowed: %v", s.extensions)}
	}
	ext := parts[len(parts)-1]
	for _, e := range s.extensions {
		if ext == e {
			return nil
		}
	}
	return &ValidationError{fmt.Sprintf("Invalid file extension. Allowed: %v", s.extensions)}
}

// Process runs the full grading pipeline for a submission.
func (s *SubmissionService) Process(req SubmissionRequest, lab *config.Lab, imageName string) (*models.GradeResultT, error) {
	// Validate extension
	if err := s.ValidateExtension(req.Filename); err != nil {
		return nil, err
	}

	// Create submission directory
	id := time.Now().Format("20060102150405") + strconv.Itoa(rand.Intn(10000))
	subDir := filepath.Join(s.filesDir, id)
	if err := os.MkdirAll(subDir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("creating submission dir: %w", err)
	}
	defer os.RemoveAll(subDir)

	// Write file
	destPath := filepath.Join(subDir, req.Filename)
	if err := os.WriteFile(destPath, req.Code, 0644); err != nil {
		return nil, fmt.Errorf("writing submission file: %w", err)
	}

	absDir, err := filepath.Abs(subDir)
	if err != nil {
		return nil, fmt.Errorf("resolving path: %w", err)
	}

	// Build marker input
	input := s.grader.BuildInput(lab)
	input.Filename = s.mountPath + req.Filename
	if err := grader.WriteInput(input, filepath.Join(absDir, "input.json")); err != nil {
		return nil, fmt.Errorf("writing input config: %w", err)
	}

	// Resolve bind mount path
	bindDir := absDir
	if s.hostDir != "" {
		bindDir = filepath.Join(s.hostDir, id)
	}

	// Run container
	sub := &docker.Submission{
		ContainerName: id,
		ImageName:     imageName,
		TargetDir:     s.mountPath,
		BindedDir:     bindDir,
		Timeout:       s.timeout,
	}
	if err := docker.RunContainer(s.log, s.docker, sub); err != nil {
		return nil, fmt.Errorf("grading failed: %w", err)
	}

	// Read & evaluate output
	output, err := s.grader.ReadOutput(filepath.Join(absDir, "output.json"))
	if err != nil {
		return nil, fmt.Errorf("reading results: %w", err)
	}

	return s.grader.Evaluate(lab, output)
}
