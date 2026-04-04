package config

import (
	"os"
	"path/filepath"
	"testing"
)

var envKeys = []string{"AUTOGRADER_HOST", "AUTOGRADER_PORT", "AUTOGRADER_HOST_FILES_DIR", "AUTOGRADER_MARKER_IMAGE", "REDIS_URL", "REDIS_RATE_LIMIT", "REDIS_MAX_RETRY"}

func clearEnv(t *testing.T) {
	for _, k := range envKeys {
		t.Setenv(k, "")
	}
}

func TestLoad_Valid(t *testing.T) {
	clearEnv(t)
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yml")
	os.WriteFile(path, []byte(`
server:
  host: 127.0.0.1
  server_port: "8080"
labs:
  - id: lab_1
    name: Test Lab
marker:
  image_name: testimg
`), 0644)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Server.Host != "127.0.0.1" {
		t.Errorf("expected host 127.0.0.1, got %q", cfg.Server.Host)
	}
	if cfg.Server.Port != "8080" {
		t.Errorf("expected port 8080, got %q", cfg.Server.Port)
	}
	if len(cfg.Labs) != 1 {
		t.Errorf("expected 1 lab, got %d", len(cfg.Labs))
	}
	// Redis defaults should be set even without yaml
	if cfg.Redis.RedisServer != "redis://0.0.0.0:6379" {
		t.Errorf("expected default redis URL, got %q", cfg.Redis.RedisServer)
	}
}

func TestLoad_Defaults(t *testing.T) {
	clearEnv(t)
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yml")
	os.WriteFile(path, []byte("{}"), 0644)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Server.Port != "9090" {
		t.Errorf("expected default port 9090, got %q", cfg.Server.Port)
	}
	if cfg.Marker.ContainerTimeout != 10 {
		t.Errorf("expected default timeout 10, got %d", cfg.Marker.ContainerTimeout)
	}
	if len(cfg.Marker.AllowedExtensions) != 1 || cfg.Marker.AllowedExtensions[0] != "py" {
		t.Errorf("expected default extensions [py], got %v", cfg.Marker.AllowedExtensions)
	}
}

func TestLoad_EnvOverrides(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yml")
	os.WriteFile(path, []byte("{}"), 0644)

	t.Setenv("AUTOGRADER_HOST", "10.0.0.1")
	t.Setenv("AUTOGRADER_PORT", "3000")
	t.Setenv("REDIS_URL", "redis://redis:6379")
	t.Setenv("REDIS_RATE_LIMIT", "100-H")
	t.Setenv("REDIS_MAX_RETRY", "5")

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Server.Host != "10.0.0.1" {
		t.Errorf("expected host 10.0.0.1, got %q", cfg.Server.Host)
	}
	if cfg.Server.Port != "3000" {
		t.Errorf("expected port 3000, got %q", cfg.Server.Port)
	}
	if cfg.Redis.RedisServer != "redis://redis:6379" {
		t.Errorf("expected redis URL override, got %q", cfg.Redis.RedisServer)
	}
	if cfg.Redis.RateLimiter != "100-H" {
		t.Errorf("expected rate limit 100-H, got %q", cfg.Redis.RateLimiter)
	}
	if cfg.Redis.MaxRetry != 5 {
		t.Errorf("expected max retry 5, got %d", cfg.Redis.MaxRetry)
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := Load("/nonexistent/config.yml")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLoad_InvalidYAML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yml")
	os.WriteFile(path, []byte(":::invalid"), 0644)

	_, err := Load(path)
	if err == nil {
		t.Error("expected error for invalid YAML")
	}
}
