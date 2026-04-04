package config

import (
	"autograder/internal/models"
	"fmt"
	"os"
	"strconv"

	"gopkg.in/yaml.v2"
)

// Config is the top-level application configuration.
type Config struct {
	Labs    []Lab `yaml:"labs"`
	Redis   RedisConfig
	Server  ServerConfig `yaml:"server"`
	Logging LogConfig    `yaml:"logging"`
	Marker  MarkerConfig `yaml:"marker"`
}

// ServerConfig holds HTTP server settings.
type ServerConfig struct {
	Port         string `yaml:"server_port"`
	Host         string `yaml:"host"`
	WriteTimeout int32  `yaml:"write_timeout"`
	ReadTimeout  int32  `yaml:"read_timeout"`
}

// RedisConfig holds Redis connection and rate-limiting settings.
// Configured entirely via environment variables, not the YAML config file.
type RedisConfig struct {
	MaxRetry    int
	RateLimiter string
	RedisServer string
}

// LogConfig holds logging settings.
type LogConfig struct {
	LogDir   string `yaml:"log_dir"`
	LogLevel string `yaml:"log_level"`
}

// MarkerConfig holds Docker/marker settings.
type MarkerConfig struct {
	DockerfilePath    string   `yaml:"dockerfile_path"`
	SubmissionsFolder string   `yaml:"submissions_folder"`
	MountPath         string   `yaml:"mount_path"`
	Command           string   `yaml:"command"`
	ImageName         string   `yaml:"image_name"`
	AllowedExtensions []string `yaml:"allowed_extensions"`
	ContainerTimeout  int      `yaml:"container_timeout"`
	// HostSubmissionsFolder is the path to the submissions folder as seen by
	// the Docker daemon. When the server runs inside Docker with a mounted
	// socket, this must be set to the host-side volume/bind path so that
	// grading containers can access submission files. When empty, defaults
	// to SubmissionsFolder (assumes the server runs directly on the host).
	HostSubmissionsFolder string `yaml:"host_submissions_folder"`
}

// Lab represents a lab assignment.
type Lab struct {
	Name             string     `yaml:"name" json:"name"`
	ID               string     `yaml:"id" json:"id"`
	ProblemStatement string     `yaml:"problem_statement" json:"problem_statement"`
	Testcase         []Testcase `yaml:"testcase" json:"-"`
}

// Testcase defines a single test within a lab.
type Testcase struct {
	Expected  []Expected `yaml:"expected"`
	Type      string     `yaml:"type"`
	Name      string     `yaml:"name"`
	Functions []Function `yaml:"functions"`
}

// Expected defines grading criteria.
type Expected struct {
	Feedback string   `yaml:"feedback"`
	Points   float32  `yaml:"points"`
	Values   []string `yaml:"values"`
}

// Function describes a function call to test.
type Function struct {
	FunctionName string        `yaml:"function_name" json:"function_name"`
	FunctionArgs []FunctionArg `yaml:"function_args" json:"function_args"`
	TestcaseName string        `json:"testcase_name"`
}

// FunctionArg is a typed argument for a function call.
type FunctionArg struct {
	Type  string      `yaml:"type" json:"type"`
	Value interface{} `yaml:"value" json:"value"`
}

// LabToModel converts a config Lab to the generated FlatBuffers model.
func LabToModel(l *Lab) *models.LabT {
	return &models.LabT{
		Id:               l.ID,
		Name:             l.Name,
		ProblemStatement: l.ProblemStatement,
	}
}

// Load reads and parses the config file, applying defaults and env overrides.
func Load(path string) (*Config, error) {
	cfg := defaults()

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config: %w", err)
	}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	applyEnvOverrides(&cfg)
	return &cfg, nil
}

func defaults() Config {
	return Config{
		Redis: RedisConfig{
			MaxRetry:    3,
			RedisServer: "redis://0.0.0.0:6379",
			RateLimiter: "50-H",
		},
		Server: ServerConfig{
			Host:         "0.0.0.0",
			Port:         "9090",
			ReadTimeout:  5,
			WriteTimeout: 5,
		},
		Marker: MarkerConfig{
			AllowedExtensions: []string{"py"},
			ContainerTimeout:  10,
			SubmissionsFolder: "./files/",
			MountPath:         "/mnt/",
		},
	}
}

func applyEnvOverrides(cfg *Config) {
	if v := os.Getenv("AUTOGRADER_HOST"); v != "" {
		cfg.Server.Host = v
	}
	if v := os.Getenv("AUTOGRADER_PORT"); v != "" {
		cfg.Server.Port = v
	}
	if v := os.Getenv("AUTOGRADER_HOST_FILES_DIR"); v != "" {
		cfg.Marker.HostSubmissionsFolder = v
	}
	if v := os.Getenv("AUTOGRADER_MARKER_IMAGE"); v != "" {
		cfg.Marker.ImageName = v
	}

	// Redis — env only, with defaults
	if v := os.Getenv("REDIS_URL"); v != "" {
		cfg.Redis.RedisServer = v
	}
	if v := os.Getenv("REDIS_RATE_LIMIT"); v != "" {
		cfg.Redis.RateLimiter = v
	}
	if v := os.Getenv("REDIS_MAX_RETRY"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			cfg.Redis.MaxRetry = n
		}
	}
}
