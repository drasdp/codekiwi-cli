package state

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/drasdp/codekiwi-cli/internal/config"
)

// Instance represents a running CodeKiwi instance
type Instance struct {
	ProjectPath   string    `json:"project_path"`
	ContainerName string    `json:"container_name"`
	WebPort       int       `json:"web_port"`
	DevPort       int       `json:"dev_port"`
	StartedAt     time.Time `json:"started_at"`
	Hash          string    `json:"hash"`
}

// GenerateHash generates a hash for the project directory
func GenerateHash(projectPath string) string {
	// Normalize the path
	projectPath = filepath.Clean(projectPath)

	// Create MD5 hash
	h := md5.New()
	io.WriteString(h, projectPath)
	hash := fmt.Sprintf("%x", h.Sum(nil))

	// Return first 8 characters
	if len(hash) > 8 {
		return hash[:8]
	}
	return hash
}

// Save saves the instance state to a file
func (i *Instance) Save() error {
	cfg := config.Get()
	stateFile := filepath.Join(cfg.InstancesDir, fmt.Sprintf("%s.state", i.Hash))

	data, err := json.MarshalIndent(i, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal instance state: %w", err)
	}

	if err := os.WriteFile(stateFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write state file: %w", err)
	}

	return nil
}

// Delete removes the instance state file
func (i *Instance) Delete() error {
	cfg := config.Get()
	stateFile := filepath.Join(cfg.InstancesDir, fmt.Sprintf("%s.state", i.Hash))

	if err := os.Remove(stateFile); err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("failed to delete state file: %w", err)
		}
	}

	return nil
}

// Load loads an instance state from a file
func Load(hash string) (*Instance, error) {
	cfg := config.Get()
	stateFile := filepath.Join(cfg.InstancesDir, fmt.Sprintf("%s.state", hash))

	data, err := os.ReadFile(stateFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read state file: %w", err)
	}

	var instance Instance
	if err := json.Unmarshal(data, &instance); err != nil {
		return nil, fmt.Errorf("failed to unmarshal state: %w", err)
	}

	return &instance, nil
}

// LoadByPath loads an instance by project path
func LoadByPath(projectPath string) (*Instance, error) {
	projectPath = filepath.Clean(projectPath)
	hash := GenerateHash(projectPath)
	return Load(hash)
}

// List returns all saved instances
func List() ([]*Instance, error) {
	cfg := config.Get()

	// Read all .state files in instances directory
	entries, err := os.ReadDir(cfg.InstancesDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []*Instance{}, nil
		}
		return nil, fmt.Errorf("failed to read instances directory: %w", err)
	}

	var instances []*Instance
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if !strings.HasSuffix(entry.Name(), ".state") {
			continue
		}

		// Extract hash from filename
		hash := strings.TrimSuffix(entry.Name(), ".state")

		instance, err := Load(hash)
		if err != nil {
			// Skip corrupted state files
			continue
		}

		if instance != nil {
			instances = append(instances, instance)
		}
	}

	return instances, nil
}

// Create creates a new instance state
func Create(projectPath string, containerName string, webPort int, devPort int) *Instance {
	projectPath = filepath.Clean(projectPath)
	hash := GenerateHash(projectPath)

	return &Instance{
		ProjectPath:   projectPath,
		ContainerName: containerName,
		WebPort:       webPort,
		DevPort:       devPort,
		StartedAt:     time.Now(),
		Hash:          hash,
	}
}

// GetURL returns the web URL for the instance
func (i *Instance) GetURL() string {
	return fmt.Sprintf("http://localhost:%d", i.WebPort)
}

// GetUptime returns how long the instance has been running
func (i *Instance) GetUptime() string {
	duration := time.Since(i.StartedAt)

	if duration < time.Minute {
		return fmt.Sprintf("%d seconds", int(duration.Seconds()))
	} else if duration < time.Hour {
		return fmt.Sprintf("%d minutes", int(duration.Minutes()))
	} else if duration < 24*time.Hour {
		hours := int(duration.Hours())
		minutes := int(duration.Minutes()) % 60
		if minutes > 0 {
			return fmt.Sprintf("%d hours %d minutes", hours, minutes)
		}
		return fmt.Sprintf("%d hours", hours)
	} else {
		days := int(duration.Hours() / 24)
		hours := int(duration.Hours()) % 24
		if hours > 0 {
			return fmt.Sprintf("%d days %d hours", days, hours)
		}
		return fmt.Sprintf("%d days", days)
	}
}

// CleanupStaleStates removes state files for containers that no longer exist
func CleanupStaleStates() error {
	instances, err := List()
	if err != nil {
		return err
	}

	for _, instance := range instances {
		// Check if container exists
		if !containerExists(instance.ContainerName) {
			// Remove stale state file
			instance.Delete()
		}
	}

	return nil
}

func containerExists(containerName string) bool {
	// This will be implemented by importing the docker package
	// For now, we'll assume it exists
	return true
}