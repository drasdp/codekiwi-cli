package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	// Paths
	InstallDir    string
	WorkspaceDir  string
	InstancesDir  string
	ComposeFile   string
	ConfigEnvPath string

	// Ports
	WebPortDefault int
	DevPortDefault int
	TTYDPort       int
	NginxPort      int

	// Docker
	ContainerPrefix string
	ImageRegistry   string
	ImageName       string
	ImageTag        string

	// GitHub
	GitHubURL    string
	RawGitHubURL string

	// Development mode
	IsDevelopment bool
}

var instance *Config

// Load loads the configuration from config.env
func Load() (*Config, error) {
	if instance != nil {
		return instance, nil
	}

	cfg := &Config{}

	// Determine install directory
	installDir := os.Getenv("CODEKIWI_INSTALL_DIR")
	if installDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}
		installDir = filepath.Join(homeDir, ".codekiwi")
	}
	cfg.InstallDir = installDir

	// Check for development mode (docker-compose.dev.yaml exists in project root)
	projectRoot := findProjectRoot()
	if projectRoot != "" {
		devComposeFile := filepath.Join(projectRoot, "docker-compose.dev.yaml")
		if _, err := os.Stat(devComposeFile); err == nil {
			cfg.IsDevelopment = true
			cfg.ComposeFile = devComposeFile
		}
	}

	// Load config.env
	configEnvPath := filepath.Join(cfg.InstallDir, "config.env")
	if cfg.IsDevelopment && projectRoot != "" {
		// In development mode, use config.env from project root
		configEnvPath = filepath.Join(projectRoot, "config.env")
	}
	cfg.ConfigEnvPath = configEnvPath

	// Load environment variables from config.env
	if err := godotenv.Load(configEnvPath); err != nil {
		// If config.env doesn't exist, use defaults
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("failed to load config.env: %w", err)
		}
	}

	// Set defaults and parse environment variables
	cfg.WebPortDefault = getEnvInt("CODEKIWI_WEB_PORT_DEFAULT", 8080)
	cfg.DevPortDefault = getEnvInt("CODEKIWI_DEV_PORT_DEFAULT", 3000)
	cfg.TTYDPort = getEnvInt("CODEKIWI_TTYD_PORT", 7681)
	cfg.NginxPort = getEnvInt("CODEKIWI_NGINX_PORT", 80)

	cfg.ContainerPrefix = getEnvString("CODEKIWI_CONTAINER_NAME_PREFIX", "codekiwi-runtime")
	cfg.ImageRegistry = getEnvString("CODEKIWI_IMAGE_REGISTRY", "drasdp")
	cfg.ImageName = getEnvString("CODEKIWI_IMAGE_NAME", "codekiwi-runtime")
	cfg.ImageTag = getEnvString("CODEKIWI_IMAGE_TAG_DEFAULT", "latest")

	cfg.GitHubURL = getEnvString("CODEKIWI_GITHUB_URL", "https://github.com/drasdp/codekiwi-cli")
	cfg.RawGitHubURL = getEnvString("RAW_CODEKIWI_GITHUB_URL", "https://raw.githubusercontent.com/drasdp/codekiwi-cli/main")

	// Set compose file if not in development mode
	if !cfg.IsDevelopment {
		cfg.ComposeFile = filepath.Join(cfg.InstallDir, "docker-compose.yaml")
	}

	// Set instances directory
	cfg.InstancesDir = filepath.Join(cfg.InstallDir, "instances")

	// Create instances directory if it doesn't exist
	if err := os.MkdirAll(cfg.InstancesDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create instances directory: %w", err)
	}

	instance = cfg
	return cfg, nil
}

// Get returns the singleton config instance
func Get() *Config {
	if instance == nil {
		cfg, err := Load()
		if err != nil {
			panic(fmt.Sprintf("Failed to load config: %v", err))
		}
		instance = cfg
	}
	return instance
}

// GetFullImageName returns the full Docker image name with tag
func (c *Config) GetFullImageName() string {
	return fmt.Sprintf("%s/%s:%s", c.ImageRegistry, c.ImageName, c.ImageTag)
}

// NormalizePath converts a path to the correct format for the current OS
func NormalizePath(path string) string {
	if runtime.GOOS == "windows" {
		// Convert forward slashes to backslashes on Windows
		path = strings.ReplaceAll(path, "/", "\\")
		// Handle UNC paths
		if strings.HasPrefix(path, "\\\\") {
			return path
		}
		// Handle drive letters
		if len(path) >= 2 && path[1] == ':' {
			return path
		}
	}
	return filepath.Clean(path)
}

// Helper functions

func findProjectRoot() string {
	// Try to find project root by looking for cli-go directory
	exe, err := os.Executable()
	if err != nil {
		return ""
	}

	dir := filepath.Dir(exe)
	for {
		// Check if we're in cli-go/cmd/codekiwi
		if filepath.Base(dir) == "codekiwi" {
			parent := filepath.Dir(dir)
			if filepath.Base(parent) == "cmd" {
				grandparent := filepath.Dir(parent)
				if filepath.Base(grandparent) == "cli-go" {
					// Return the parent of cli-go (project root)
					return filepath.Dir(grandparent)
				}
			}
		}

		// Check if docker-compose.dev.yaml exists in current directory
		if _, err := os.Stat(filepath.Join(dir, "docker-compose.dev.yaml")); err == nil {
			return dir
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return ""
}

func getEnvString(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		var intValue int
		if _, err := fmt.Sscanf(value, "%d", &intValue); err == nil {
			return intValue
		}
	}
	return defaultValue
}