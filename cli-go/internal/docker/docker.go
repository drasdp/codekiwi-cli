package docker

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/drasdp/codekiwi-cli/internal/config"
	"github.com/drasdp/codekiwi-cli/internal/platform"
)

// CheckDocker verifies Docker is installed and running
func CheckDocker() error {
	// Check if docker command exists
	if _, err := exec.LookPath("docker"); err != nil {
		return fmt.Errorf("Docker is not installed. Please install Docker from https://www.docker.com/get-started")
	}

	// Check if Docker daemon is running
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "docker", "info")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Docker daemon is not running. Please start Docker Desktop")
	}

	return nil
}

// CheckDockerCompose checks if docker-compose is available
func CheckDockerCompose() error {
	// Try docker compose (v2) first
	cmd := exec.Command("docker", "compose", "version")
	if err := cmd.Run(); err == nil {
		return nil
	}

	// Fall back to docker-compose (v1)
	if _, err := exec.LookPath("docker-compose"); err != nil {
		return fmt.Errorf("docker-compose is not installed")
	}

	return nil
}

// ComposeUp starts the Docker container using docker-compose
func ComposeUp(projectPath string, containerName string, webPort int, devPort int) error {
	cfg := config.Get()

	// Set environment variables
	env := os.Environ()
	env = append(env,
		fmt.Sprintf("WORKSPACE_DIR=%s", projectPath),
		fmt.Sprintf("CONTAINER_NAME=%s", containerName),
		fmt.Sprintf("WEB_PORT=%d", webPort),
		fmt.Sprintf("DEV_PORT=%d", devPort),
		fmt.Sprintf("CODEKIWI_WORKSPACE_DIR=/workspace"),
		fmt.Sprintf("CODEKIWI_INSTALL_DIR_NAME=.codekiwi"),
		fmt.Sprintf("CODEKIWI_IMAGE_NAME=%s", cfg.GetFullImageName()),
		fmt.Sprintf("AUTH_DIR=%s", filepath.Join(cfg.InstallDir, "opencode")),
	)

	// Determine docker-compose command
	composeCmd := getDockerComposeCommand()

	// Build command arguments
	args := []string{}
	if composeCmd == "docker" {
		args = append(args, "compose")
	}

	args = append(args,
		"-f", cfg.ComposeFile,
		"-p", containerName,
		"up", "-d",
	)

	// Add --build flag in development mode
	if cfg.IsDevelopment {
		args = append(args, "--build")
	}

	// Execute docker-compose up
	cmd := exec.Command(composeCmd, args...)
	cmd.Env = env
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	platform.PrintInfo(fmt.Sprintf("Starting container %s...", containerName))

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to start container: %w", err)
	}

	// Wait for container to be ready
	if err := WaitForContainer(containerName, 30*time.Second); err != nil {
		return fmt.Errorf("container failed to start: %w", err)
	}

	return nil
}

// ComposeDown stops and removes the Docker container
func ComposeDown(containerName string) error {
	cfg := config.Get()
	composeCmd := getDockerComposeCommand()

	// Build command arguments
	args := []string{}
	if composeCmd == "docker" {
		args = append(args, "compose")
	}

	args = append(args,
		"-f", cfg.ComposeFile,
		"-p", containerName,
		"down",
	)

	cmd := exec.Command(composeCmd, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to stop container: %w", err)
	}

	return nil
}

// StopContainer stops a running container
func StopContainer(containerName string) error {
	cmd := exec.Command("docker", "stop", containerName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to stop container %s: %w", containerName, err)
	}
	return nil
}

// RemoveContainer removes a container
func RemoveContainer(containerName string) error {
	cmd := exec.Command("docker", "rm", "-f", containerName)
	if err := cmd.Run(); err != nil {
		// Ignore error if container doesn't exist
		if !strings.Contains(err.Error(), "No such container") {
			return fmt.Errorf("failed to remove container %s: %w", containerName, err)
		}
	}
	return nil
}

// IsContainerRunning checks if a container is running
func IsContainerRunning(containerName string) bool {
	cmd := exec.Command("docker", "ps", "--format", "{{.Names}}")
	output, err := cmd.Output()
	if err != nil {
		return false
	}

	containers := strings.Split(string(output), "\n")
	for _, name := range containers {
		if strings.TrimSpace(name) == containerName {
			return true
		}
	}

	return false
}

// WaitForContainer waits for a container to be running and healthy
func WaitForContainer(containerName string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		if IsContainerRunning(containerName) {
			// Container is running, wait a bit for services to start
			time.Sleep(2 * time.Second)
			return nil
		}
		time.Sleep(500 * time.Millisecond)
	}

	return fmt.Errorf("timeout waiting for container %s to start", containerName)
}

// GetContainerLogs retrieves logs from a container
func GetContainerLogs(containerName string, tail int) (string, error) {
	cmd := exec.Command("docker", "logs", "--tail", fmt.Sprintf("%d", tail), containerName)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get logs: %w", err)
	}
	return string(output), nil
}

// FollowContainerLogs follows container logs in real-time
func FollowContainerLogs(containerName string) error {
	cmd := exec.Command("docker", "logs", "-f", containerName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// PullImage pulls the Docker image
func PullImage(imageName string) error {
	platform.PrintInfo(fmt.Sprintf("Pulling image %s...", imageName))

	cmd := exec.Command("docker", "pull", imageName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to pull image: %w", err)
	}

	return nil
}

// ListContainers lists all CodeKiwi containers
func ListContainers() ([]string, error) {
	cfg := config.Get()

	cmd := exec.Command("docker", "ps", "--format", "{{.Names}}")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list containers: %w", err)
	}

	var codeKiwiContainers []string
	containers := strings.Split(string(output), "\n")
	for _, name := range containers {
		name = strings.TrimSpace(name)
		if strings.HasPrefix(name, cfg.ContainerPrefix) {
			codeKiwiContainers = append(codeKiwiContainers, name)
		}
	}

	return codeKiwiContainers, nil
}

// Helper functions

func getDockerComposeCommand() string {
	// Try docker compose v2 first
	cmd := exec.Command("docker", "compose", "version")
	if err := cmd.Run(); err == nil {
		return "docker"
	}

	// Fall back to docker-compose v1
	if _, err := exec.LookPath("docker-compose"); err == nil {
		return "docker-compose"
	}

	// Default to docker compose (will error if not available)
	return "docker"
}