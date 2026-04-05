package api

import (
	"autograder/internal/config"
	"autograder/internal/grader"
	"autograder/internal/models"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/moby/moby/client"
	"github.com/sirupsen/logrus"
)

func quiet() *logrus.Logger {
	l := logrus.New()
	l.SetLevel(logrus.PanicLevel)
	return l
}

// mockDockerClient satisfies docker.Client for handler tests.
type mockDockerClient struct{}

func (m *mockDockerClient) ImageList(_ context.Context, _ client.ImageListOptions) (client.ImageListResult, error) {
	return client.ImageListResult{}, nil
}
func (m *mockDockerClient) ImageRemove(_ context.Context, _ string, _ client.ImageRemoveOptions) (client.ImageRemoveResult, error) {
	return client.ImageRemoveResult{}, nil
}
func (m *mockDockerClient) ImageBuild(_ context.Context, _ io.Reader, _ client.ImageBuildOptions) (client.ImageBuildResult, error) {
	return client.ImageBuildResult{Body: io.NopCloser(io.Reader(nil))}, nil
}
func (m *mockDockerClient) ContainerCreate(_ context.Context, _ client.ContainerCreateOptions) (client.ContainerCreateResult, error) {
	return client.ContainerCreateResult{ID: "test"}, nil
}
func (m *mockDockerClient) ContainerStart(_ context.Context, _ string, _ client.ContainerStartOptions) (client.ContainerStartResult, error) {
	return client.ContainerStartResult{}, fmt.Errorf("no docker in test")
}
func (m *mockDockerClient) ContainerWait(_ context.Context, _ string, _ client.ContainerWaitOptions) client.ContainerWaitResult {
	return client.ContainerWaitResult{}
}
func (m *mockDockerClient) ContainerRemove(_ context.Context, _ string, _ client.ContainerRemoveOptions) (client.ContainerRemoveResult, error) {
	return client.ContainerRemoveResult{}, nil
}

func testHandler() *handler {
	return &handler{
		cfg: &config.Config{
			Labs: []config.Lab{
				{ID: "lab_1", Name: "Lab 1"},
				{ID: "lab_2", Name: "Lab 2"},
			},
			Marker: config.MarkerConfig{
				AllowedExtensions: []string{"py"},
				SubmissionsFolder: "./test_files/",
				MountPath:         "/mnt/",
				ImageName:         "autograder",
				ContainerTimeout:  5,
			},
		},
		log:    quiet(),
		grader: &grader.DefaultGrader{},
		docker: &mockDockerClient{},
		submission: &SubmissionService{
			grader:     &grader.DefaultGrader{},
			docker:     &mockDockerClient{},
			log:        quiet(),
			filesDir:   "./test_files/",
			mountPath:  "/mnt/",
			extensions: []string{"py"},
			timeout:    5,
		},
	}
}

func TestListLabs(t *testing.T) {
	h := testHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/labs", nil)
	rr := httptest.NewRecorder()

	h.listLabs(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}

	var labs []config.Lab
	json.NewDecoder(rr.Body).Decode(&labs)
	if len(labs) != 2 {
		t.Errorf("expected 2 labs, got %d", len(labs))
	}
	if labs[0].ID != "lab_1" {
		t.Errorf("expected lab_1, got %q", labs[0].ID)
	}
}

func TestGetLab_Found(t *testing.T) {
	h := testHandler()
	router := mux.NewRouter()
	router.HandleFunc("/api/labs/{id}", h.getLab).Methods("GET")

	req := httptest.NewRequest(http.MethodGet, "/api/labs/lab_1", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}

	var lab models.LabT
	json.NewDecoder(rr.Body).Decode(&lab)
	if lab.Id != "lab_1" {
		t.Errorf("expected lab_1, got %q", lab.Id)
	}
	if lab.Name != "Lab 1" {
		t.Errorf("expected 'Lab 1', got %q", lab.Name)
	}
}

func TestGetLab_NotFound(t *testing.T) {
	h := testHandler()
	router := mux.NewRouter()
	router.HandleFunc("/api/labs/{id}", h.getLab).Methods("GET")

	req := httptest.NewRequest(http.MethodGet, "/api/labs/nonexistent", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rr.Code)
	}

	var resp models.APIErrorT
	json.NewDecoder(rr.Body).Decode(&resp)
	if resp.Error == "" {
		t.Error("expected error message")
	}
}

func TestGetLab_AllLabs(t *testing.T) {
	h := testHandler()
	router := mux.NewRouter()
	router.HandleFunc("/api/labs/{id}", h.getLab).Methods("GET")

	for _, id := range []string{"lab_1", "lab_2"} {
		req := httptest.NewRequest(http.MethodGet, "/api/labs/"+id, nil)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("lab %s: expected 200, got %d", id, rr.Code)
		}

		var lab models.LabT
		json.NewDecoder(rr.Body).Decode(&lab)
		if lab.Id != id {
			t.Errorf("expected %s, got %q", id, lab.Id)
		}
	}
}

func TestSubmit_NoFile(t *testing.T) {
	h := testHandler()
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	w.WriteField("lab_id", "lab_1")
	w.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/submit", &buf)
	req.Header.Set("Content-Type", w.FormDataContentType())
	rr := httptest.NewRecorder()

	h.submit(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestSubmit_BadExtension(t *testing.T) {
	h := testHandler()
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	part, _ := w.CreateFormFile("file", "virus.exe")
	part.Write([]byte("bad"))
	w.WriteField("lab_id", "lab_1")
	w.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/submit", &buf)
	req.Header.Set("Content-Type", w.FormDataContentType())
	rr := httptest.NewRecorder()

	h.submit(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
	var resp map[string]string
	json.NewDecoder(rr.Body).Decode(&resp)
	if resp["error"] == "" {
		t.Error("expected error message")
	}
}

func TestSubmit_InvalidLab(t *testing.T) {
	h := testHandler()
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	part, _ := w.CreateFormFile("file", "script.py")
	part.Write([]byte("print('hi')"))
	w.WriteField("lab_id", "lab_999")
	w.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/submit", &buf)
	req.Header.Set("Content-Type", w.FormDataContentType())
	rr := httptest.NewRecorder()

	h.submit(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestSubmit_FileTooLarge(t *testing.T) {
	h := testHandler()
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	part, _ := w.CreateFormFile("file", "big.py")
	part.Write(make([]byte, 21*1024))
	w.WriteField("lab_id", "lab_1")
	w.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/submit", &buf)
	req.Header.Set("Content-Type", w.FormDataContentType())
	rr := httptest.NewRecorder()

	h.submit(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestSubmit_CodeField(t *testing.T) {
	h := testHandler()
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	w.WriteField("lab_id", "lab_1")
	w.WriteField("code", "print('hello')")
	w.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/submit", &buf)
	req.Header.Set("Content-Type", w.FormDataContentType())
	rr := httptest.NewRecorder()

	h.submit(rr, req)

	// Will fail at container run (no docker in test) but should get past validation
	// Check it didn't fail with "Missing file" or "Provide either" errors
	var resp map[string]interface{}
	json.NewDecoder(rr.Body).Decode(&resp)
	if errMsg, ok := resp["error"]; ok {
		msg := errMsg.(string)
		if msg == "Provide either a 'file' upload or 'code' text field" {
			t.Error("code field was not accepted")
		}
	}
}

func TestSubmit_CodeFieldCustomFilename(t *testing.T) {
	h := testHandler()
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	w.WriteField("lab_id", "lab_1")
	w.WriteField("code", "print('hello')")
	w.WriteField("filename", "my_script.py")
	w.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/submit", &buf)
	req.Header.Set("Content-Type", w.FormDataContentType())
	rr := httptest.NewRecorder()

	h.submit(rr, req)

	var resp map[string]interface{}
	json.NewDecoder(rr.Body).Decode(&resp)
	if errMsg, ok := resp["error"]; ok {
		msg := errMsg.(string)
		if msg == "Provide either a 'file' upload or 'code' text field" {
			t.Error("code field with custom filename was not accepted")
		}
	}
}

func TestSubmit_NoFileNoCode(t *testing.T) {
	h := testHandler()
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	w.WriteField("lab_id", "lab_1")
	w.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/submit", &buf)
	req.Header.Set("Content-Type", w.FormDataContentType())
	rr := httptest.NewRecorder()

	h.submit(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestSubmit_CodeFieldBadExtension(t *testing.T) {
	h := testHandler()
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	w.WriteField("lab_id", "lab_1")
	w.WriteField("code", "console.log('hi')")
	w.WriteField("filename", "hack.js")
	w.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/submit", &buf)
	req.Header.Set("Content-Type", w.FormDataContentType())
	rr := httptest.NewRecorder()

	h.submit(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestFindLab(t *testing.T) {
	labs := []config.Lab{{ID: "a"}, {ID: "b"}}

	l, err := findLab(labs, "b")
	if err != nil || l.ID != "b" {
		t.Errorf("expected lab b, got %+v, err=%v", l, err)
	}

	_, err = findLab(labs, "z")
	if err == nil {
		t.Error("expected error for missing lab")
	}
}
