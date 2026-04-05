package api

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func setupSPADir(t *testing.T) string {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "index.html"), []byte("<html>SPA</html>"), 0644)
	os.MkdirAll(filepath.Join(dir, "assets"), 0755)
	os.WriteFile(filepath.Join(dir, "assets", "app.js"), []byte("console.log('app')"), 0644)
	os.WriteFile(filepath.Join(dir, "assets", "style.css"), []byte("body{}"), 0644)
	return dir
}

func TestSPA_ServesStaticFile(t *testing.T) {
	dir := setupSPADir(t)
	spa := newSPAHandler(dir, "index.html")

	req := httptest.NewRequest("GET", "/assets/app.js", nil)
	rr := httptest.NewRecorder()
	spa.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
	if rr.Body.String() != "console.log('app')" {
		t.Errorf("unexpected body: %q", rr.Body.String())
	}
}

func TestSPA_ServesCSS(t *testing.T) {
	dir := setupSPADir(t)
	spa := newSPAHandler(dir, "index.html")

	req := httptest.NewRequest("GET", "/assets/style.css", nil)
	rr := httptest.NewRecorder()
	spa.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestSPA_FallsBackToIndex_UnknownRoute(t *testing.T) {
	dir := setupSPADir(t)
	spa := newSPAHandler(dir, "index.html")

	req := httptest.NewRequest("GET", "/labs/lab_1", nil)
	rr := httptest.NewRecorder()
	spa.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
	if rr.Body.String() != "<html>SPA</html>" {
		t.Errorf("expected index.html content, got: %q", rr.Body.String())
	}
}

func TestSPA_FallsBackToIndex_Root(t *testing.T) {
	dir := setupSPADir(t)
	spa := newSPAHandler(dir, "index.html")

	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()
	spa.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
	if rr.Body.String() != "<html>SPA</html>" {
		t.Errorf("expected index.html content, got: %q", rr.Body.String())
	}
}

func TestSPA_FallsBackToIndex_DeepRoute(t *testing.T) {
	dir := setupSPADir(t)
	spa := newSPAHandler(dir, "index.html")

	req := httptest.NewRequest("GET", "/some/deep/route", nil)
	rr := httptest.NewRecorder()
	spa.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestSPA_DirectoryDoesNotServeAsFile(t *testing.T) {
	dir := setupSPADir(t)
	spa := newSPAHandler(dir, "index.html")

	// /assets/ is a directory — should fallback to index, not list dir
	req := httptest.NewRequest("GET", "/assets/", nil)
	rr := httptest.NewRecorder()
	spa.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
	if rr.Body.String() != "<html>SPA</html>" {
		t.Errorf("expected index.html, got: %q", rr.Body.String())
	}
}
