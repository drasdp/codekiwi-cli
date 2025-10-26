package platform

import (
	"fmt"
	"os/exec"
	"runtime"
)

// OpenBrowser opens the default browser with the given URL
func OpenBrowser(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		// Use cmd.exe /c start for Windows
		cmd = exec.Command("cmd", "/c", "start", "", url)
	case "darwin":
		// Use open command for macOS
		cmd = exec.Command("open", url)
	case "linux":
		// Try xdg-open first (most common)
		if err := exec.Command("which", "xdg-open").Run(); err == nil {
			cmd = exec.Command("xdg-open", url)
		} else if err := exec.Command("which", "gnome-open").Run(); err == nil {
			cmd = exec.Command("gnome-open", url)
		} else if err := exec.Command("which", "kde-open").Run(); err == nil {
			cmd = exec.Command("kde-open", url)
		} else {
			return fmt.Errorf("no suitable browser opener found")
		}
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	// Start the command without waiting
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to open browser: %w", err)
	}

	// Detach from the process
	go cmd.Wait()

	return nil
}

// CanOpenBrowser checks if the system can open a browser
func CanOpenBrowser() bool {
	switch runtime.GOOS {
	case "windows", "darwin":
		return true
	case "linux":
		// Check if we have a display (not SSH session)
		if exec.Command("which", "xdg-open").Run() == nil {
			// Check if DISPLAY is set (GUI available)
			cmd := exec.Command("sh", "-c", "echo $DISPLAY")
			output, err := cmd.Output()
			if err == nil && len(output) > 1 {
				return true
			}
		}
		return false
	default:
		return false
	}
}