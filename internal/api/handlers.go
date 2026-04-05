package api

import (
	"autograder/internal/config"
	"autograder/internal/docker"
	"autograder/internal/grader"
	"autograder/internal/models"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type handler struct {
	cfg    *config.Config
	log    *logrus.Logger
	docker docker.Client
	grader grader.Grader
}

// --- JSON helpers ---

func jsonOK(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func jsonError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(&models.APIErrorT{Error: msg})
}

// --- Endpoints ---

// GET /api/labs — returns []models.LabT
func (h *handler) listLabs(w http.ResponseWriter, r *http.Request) {
	out := make([]*models.LabT, len(h.cfg.Labs))
	for i, l := range h.cfg.Labs {
		out[i] = config.LabToModel(&l)
	}
	jsonOK(w, out)
}

// GET /api/labs/{id} — returns a single models.LabT
func (h *handler) getLab(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	lab, err := findLab(h.cfg.Labs, id)
	if err != nil {
		jsonError(w, http.StatusNotFound, fmt.Sprintf("Lab %q not found", id))
		return
	}

	jsonOK(w, config.LabToModel(&lab))
}

// POST /api/submit
// Accepts either:
//   - multipart form with "file" field (file upload)
//   - multipart form with "code" field (inline code text)
//
// Both require "lab_id". When using "code", "filename" is optional (defaults to "solution.py").
func (h *handler) submit(w http.ResponseWriter, r *http.Request) {
	// 20 KB limit
	r.Body = http.MaxBytesReader(w, r.Body, 20*1024)
	if err := r.ParseMultipartForm(20 * 1024); err != nil {
		jsonError(w, http.StatusBadRequest, "Request too large or malformed")
		return
	}

	labID := r.FormValue("lab_id")
	lab, err := findLab(h.cfg.Labs, labID)
	if err != nil {
		jsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Resolve file content and name from either "file" upload or "code" text
	var fileContent []byte
	var fileName string

	code := r.FormValue("code")
	if code != "" {
		// Inline code submission
		fileContent = []byte(code)
		fileName = r.FormValue("filename")
		if fileName == "" {
			fileName = "solution.py"
		}
	} else {
		// File upload submission
		file, fh, err := r.FormFile("file")
		if err != nil {
			jsonError(w, http.StatusBadRequest, "Provide either a 'file' upload or 'code' text field")
			return
		}
		defer file.Close()

		fileContent, err = io.ReadAll(file)
		if err != nil {
			jsonError(w, http.StatusBadRequest, "Failed to read uploaded file")
			return
		}
		fileName = fh.Filename
	}

	// Validate extension
	exts := h.cfg.Marker.AllowedExtensions
	if len(exts) == 0 {
		exts = []string{"py"}
	}
	if !hasAllowedExtension(fileName, exts) {
		jsonError(w, http.StatusBadRequest, fmt.Sprintf("Invalid file extension. Allowed: %v", exts))
		return
	}

	// Create submission directory
	id := time.Now().Format("20060102150405") + strconv.Itoa(rand.Intn(10000))
	subDir := filepath.Join(h.cfg.Marker.SubmissionsFolder, id)
	if err := os.MkdirAll(subDir, os.ModePerm); err != nil {
		h.log.WithError(err).Error("failed to create submission dir")
		jsonError(w, http.StatusInternalServerError, "Server error")
		return
	}
	defer os.RemoveAll(subDir)

	// Write file
	destPath := filepath.Join(subDir, fileName)
	if err := os.WriteFile(destPath, fileContent, 0644); err != nil {
		h.log.WithError(err).Error("failed to write file")
		jsonError(w, http.StatusInternalServerError, "Server error")
		return
	}

	absDir, err := filepath.Abs(subDir)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "Server error")
		return
	}

	// Build marker input
	input := h.grader.BuildInput(&lab)
	input.Filename = h.cfg.Marker.MountPath + fileName
	if err := grader.WriteInput(input, filepath.Join(absDir, "input.json")); err != nil {
		h.log.WithError(err).Error("failed to write input")
		jsonError(w, http.StatusInternalServerError, "Server error")
		return
	}

	// Resolve bind mount path for Docker
	bindDir := absDir
	if h.cfg.Marker.HostSubmissionsFolder != "" {
		bindDir = filepath.Join(h.cfg.Marker.HostSubmissionsFolder, id)
	}

	// Run container
	timeout := h.cfg.Marker.ContainerTimeout
	if timeout <= 0 {
		timeout = 10
	}
	sub := &docker.Submission{
		ContainerName: id,
		ImageName:     h.cfg.Marker.ImageName,
		TargetDir:     h.cfg.Marker.MountPath,
		BindedDir:     bindDir,
		Timeout:       timeout,
	}
	if err := docker.RunContainer(h.log, h.docker, sub); err != nil {
		h.log.WithError(err).Error("container failed")
		jsonError(w, http.StatusInternalServerError, "Grading failed: "+err.Error())
		return
	}

	// Read & evaluate output
	output, err := h.grader.ReadOutput(filepath.Join(absDir, "output.json"))
	if err != nil {
		h.log.WithError(err).Error("failed to read output")
		jsonError(w, http.StatusInternalServerError, "Failed to read grading results")
		return
	}

	result, err := h.grader.Evaluate(&lab, output)
	if err != nil {
		h.log.WithError(err).Error("evaluation failed")
		jsonError(w, http.StatusInternalServerError, "Evaluation failed")
		return
	}

	jsonOK(w, result)
}

// --- helpers ---

func findLab(labs []config.Lab, id string) (config.Lab, error) {
	for _, l := range labs {
		if l.ID == id {
			return l, nil
		}
	}
	return config.Lab{}, fmt.Errorf("invalid lab ID: %q", id)
}

func hasAllowedExtension(filename string, exts []string) bool {
	parts := strings.Split(filename, ".")
	if len(parts) < 2 {
		return false
	}
	ext := parts[len(parts)-1]
	for _, e := range exts {
		if ext == e {
			return true
		}
	}
	return false
}
