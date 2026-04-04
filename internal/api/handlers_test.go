package api

import (
	"autograder/internal/config"
	"autograder/internal/grader"
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sirupsen/logrus"
)

func quiet() *logrus.Logger {
	l := logrus.New()
	l.SetLevel(logrus.PanicLevel)
	return l
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

func TestHasAllowedExtension(t *testing.T) {
	tests := []struct {
		name string
		exts []string
		want bool
	}{
		{"script.py", []string{"py"}, true},
		{"script.js", []string{"py"}, false},
		{"a.tar.gz", []string{"gz"}, true},
		{"noext", []string{"py"}, false},
		{"", []string{"py"}, false},
	}
	for _, tt := range tests {
		if got := hasAllowedExtension(tt.name, tt.exts); got != tt.want {
			t.Errorf("hasAllowedExtension(%q, %v) = %v, want %v", tt.name, tt.exts, got, tt.want)
		}
	}
}
