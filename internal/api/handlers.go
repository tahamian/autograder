package api

import (
	"autograder/internal/config"
	"autograder/internal/docker"
	"autograder/internal/grader"
	"autograder/internal/models"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type handler struct {
	cfg        *config.Config
	log        *logrus.Logger
	docker     docker.Client
	grader     grader.Grader
	submission *SubmissionService
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

// GET /api/labs
func (h *handler) listLabs(w http.ResponseWriter, r *http.Request) {
	out := make([]*models.LabT, len(h.cfg.Labs))
	for i, l := range h.cfg.Labs {
		out[i] = config.LabToModel(&l)
	}
	jsonOK(w, out)
}

// GET /api/labs/{id}
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

	// Parse code or file
	var fileContent []byte
	var fileName string

	code := r.FormValue("code")
	if code != "" {
		fileContent = []byte(code)
		fileName = r.FormValue("filename")
		if fileName == "" {
			fileName = "solution.py"
		}
	} else {
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

	// Delegate to SubmissionService
	result, err := h.submission.Process(
		SubmissionRequest{LabID: labID, Code: fileContent, Filename: fileName},
		&lab,
		h.cfg.Marker.ImageName,
	)
	if err != nil {
		var ve *ValidationError
		if errors.As(err, &ve) {
			jsonError(w, http.StatusBadRequest, err.Error())
		} else {
			jsonError(w, http.StatusInternalServerError, err.Error())
		}
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
