package platform

import (
	"fmt"
	"net"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// IsPortAvailable checks if a port is available for binding
func IsPortAvailable(port int) bool {
	switch runtime.GOOS {
	case "windows":
		return isPortAvailableWindows(port)
	default:
		return isPortAvailableUnix(port)
	}
}

func isPortAvailableUnix(port int) bool {
	// Try to listen on the port
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return false
	}
	listener.Close()

	// Also check if Docker is using the port
	cmd := exec.Command("docker", "ps", "--format", "{{.Ports}}")
	output, err := cmd.Output()
	if err == nil {
		portStr := fmt.Sprintf(":%d->", port)
		if strings.Contains(string(output), portStr) {
			return false
		}
	}

	return true
}

func isPortAvailableWindows(port int) bool {
	// First try to bind to the port
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return false
	}
	listener.Close()

	// Use netstat to double-check (Windows cmd compatible)
	cmd := exec.Command("cmd", "/c", fmt.Sprintf("netstat -an | findstr :%d", port))
	output, err := cmd.Output()
	if err != nil {
		// If findstr returns error, it means no match was found (port is free)
		return true
	}

	// Check if port is in LISTENING state
	outputStr := string(output)
	portStr := fmt.Sprintf(":%d", port)
	lines := strings.Split(outputStr, "\n")
	for _, line := range lines {
		if strings.Contains(line, portStr) && strings.Contains(line, "LISTENING") {
			return false
		}
	}

	// Also check Docker
	dockerCmd := exec.Command("docker", "ps", "--format", "{{.Ports}}")
	dockerOutput, err := dockerCmd.Output()
	if err == nil {
		dockerPortStr := fmt.Sprintf(":%d->", port)
		if strings.Contains(string(dockerOutput), dockerPortStr) {
			return false
		}
	}

	return true
}

// FindAvailablePort finds an available port starting from the given port
func FindAvailablePort(startPort int, maxAttempts int) (int, error) {
	for i := 0; i < maxAttempts; i++ {
		port := startPort + i
		if IsPortAvailable(port) {
			// Double-check by trying to bind
			listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
			if err == nil {
				listener.Close()
				// Small delay to ensure port is released
				time.Sleep(100 * time.Millisecond)
				return port, nil
			}
		}
	}
	return 0, fmt.Errorf("no available port found in range %d-%d", startPort, startPort+maxAttempts-1)
}

// WaitForPort waits for a port to become available (for container startup)
func WaitForPort(host string, port int, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), time.Second)
		if err == nil {
			conn.Close()
			return nil
		}
		time.Sleep(500 * time.Millisecond)
	}

	return fmt.Errorf("timeout waiting for port %d", port)
}

// GetPortFromString extracts port number from a string
func GetPortFromString(s string) (int, error) {
	// Extract port from strings like "8080" or "localhost:8080" or "0.0.0.0:8080"
	parts := strings.Split(s, ":")
	portStr := parts[len(parts)-1]

	port, err := strconv.Atoi(strings.TrimSpace(portStr))
	if err != nil {
		return 0, fmt.Errorf("invalid port: %s", portStr)
	}

	if port < 1 || port > 65535 {
		return 0, fmt.Errorf("port out of range: %d", port)
	}

	return port, nil
}