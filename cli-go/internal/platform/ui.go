package platform

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/fatih/color"
)

var (
	// Color functions
	GreenColor  = color.New(color.FgGreen).SprintFunc()
	BlueColor   = color.New(color.FgBlue).SprintFunc()
	YellowColor = color.New(color.FgYellow).SprintFunc()
	RedColor    = color.New(color.FgRed).SprintFunc()
	CyanColor   = color.New(color.FgCyan).SprintFunc()
	BoldColor   = color.New(color.Bold).SprintFunc()

	// Emoji/icon mapping
	successIcon = "✅"
	infoIcon    = "ℹ️"
	warningIcon = "⚠️"
	errorIcon   = "❌"
	rocketIcon  = "🚀"
	folderIcon  = "📁"
	globeIcon   = "🌐"

	// Check if we should use ASCII instead of emoji
	useASCII = false
)

func init() {
	// Detect if we should use ASCII characters instead of emoji
	if runtime.GOOS == "windows" {
		// Check if we're in Windows cmd (not PowerShell or Windows Terminal)
		if isWindowsCmd() {
			useASCII = true
			// Disable colors for cmd
			color.NoColor = true
		}
	}

	// Update icons based on environment
	if useASCII {
		successIcon = "[OK]"
		infoIcon = "[INFO]"
		warningIcon = "[WARN]"
		errorIcon = "[ERROR]"
		rocketIcon = "[START]"
		folderIcon = "[DIR]"
		globeIcon = "[WEB]"
	}
}

// isWindowsCmd detects if we're running in Windows cmd (not PowerShell)
func isWindowsCmd() bool {
	// Check for PowerShell or Windows Terminal indicators
	if os.Getenv("WT_SESSION") != "" {
		// Windows Terminal
		return false
	}
	if os.Getenv("POWERSHELL_DISTRIBUTION_CHANNEL") != "" {
		// PowerShell Core
		return false
	}
	if strings.Contains(os.Getenv("PSModulePath"), "PowerShell") {
		// Windows PowerShell
		return false
	}

	// Check COMSPEC (usually points to cmd.exe)
	comspec := os.Getenv("COMSPEC")
	if strings.Contains(strings.ToLower(comspec), "cmd.exe") {
		return true
	}

	// Default to ASCII for safety on Windows
	return true
}

// Print functions with emoji/color support

func PrintSuccess(message string) {
	if useASCII {
		fmt.Printf("%s %s\n", successIcon, message)
	} else {
		fmt.Printf("%s %s\n", GreenColor(successIcon), GreenColor(message))
	}
}

func PrintInfo(message string) {
	if useASCII {
		fmt.Printf("%s %s\n", infoIcon, message)
	} else {
		fmt.Printf("%s  %s\n", BlueColor(infoIcon), BlueColor(message))
	}
}

func PrintWarning(message string) {
	if useASCII {
		fmt.Printf("%s %s\n", warningIcon, message)
	} else {
		fmt.Printf("%s  %s\n", YellowColor(warningIcon), YellowColor(message))
	}
}

func PrintError(message string) {
	if useASCII {
		fmt.Printf("%s %s\n", errorIcon, message)
	} else {
		fmt.Printf("%s %s\n", RedColor(errorIcon), RedColor(message))
	}
}

func PrintStarting(projectName string) {
	if useASCII {
		fmt.Printf("%s Starting CodeKiwi for %s\n", rocketIcon, projectName)
	} else {
		fmt.Printf("%s %s\n", rocketIcon, CyanColor(fmt.Sprintf("Starting CodeKiwi for %s", BoldColor(projectName))))
	}
}

func PrintURL(label, url string) {
	if useASCII {
		fmt.Printf("%s %s: %s\n", globeIcon, label, url)
	} else {
		fmt.Printf("%s %s: %s\n", globeIcon, BoldColor(label), CyanColor(url))
	}
}

func PrintDirectory(label, path string) {
	if useASCII {
		fmt.Printf("%s %s: %s\n", folderIcon, label, path)
	} else {
		fmt.Printf("%s %s: %s\n", folderIcon, BoldColor(label), path)
	}
}

// PrintBanner prints the CodeKiwi banner
func PrintBanner() {
	banner := `
╔═══════════════════════════════════════╗
║           CodeKiwi CLI                ║
║     Docker-based Cloud IDE            ║
╚═══════════════════════════════════════╝
`
	if useASCII {
		fmt.Print(banner)
	} else {
		fmt.Print(CyanColor(banner))
	}
}

// ClearLine clears the current line in terminal
func ClearLine() {
	fmt.Print("\r\033[K")
}

// IsColorSupported returns whether the terminal supports colors
func IsColorSupported() bool {
	return !color.NoColor
}

// EnableColor forces color output
func EnableColor() {
	color.NoColor = false
}

// DisableColor forces no color output
func DisableColor() {
	color.NoColor = true
}