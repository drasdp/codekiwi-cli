package auth

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/drasdp/codekiwi-cli/internal/platform"
)

// AuthConfig represents the structure of auth.json
type AuthConfig struct {
	OpenRouter struct {
		Type string `json:"type"`
		Key  string `json:"key"`
	} `json:"openrouter"`
}

// CheckAndSetupAPIKey checks if API key exists in auth.json and prompts if missing
func CheckAndSetupAPIKey(installDir string) error {
	authFile := filepath.Join(installDir, "opencode", "auth.json")
	authDir := filepath.Dir(authFile)

	// Create opencode directory if it doesn't exist
	if err := os.MkdirAll(authDir, 0755); err != nil {
		return fmt.Errorf("failed to create auth directory: %w", err)
	}

	// Check if auth.json exists and has valid API key
	if hasValidAPIKey(authFile) {
		return nil
	}

	// Prompt user for API key
	return promptAndSaveAPIKey(authFile)
}

// hasValidAPIKey checks if auth.json exists and contains a valid openrouter key
func hasValidAPIKey(authFile string) bool {
	// Check if file exists
	data, err := os.ReadFile(authFile)
	if err != nil {
		return false
	}

	// Parse JSON
	var config AuthConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return false
	}

	// Check if openrouter key is present and not empty
	return config.OpenRouter.Key != "" && config.OpenRouter.Key != "null"
}

// promptAndSaveAPIKey prompts the user for an API key and saves it to auth.json
func promptAndSaveAPIKey(authFile string) error {
	fmt.Println()
	platform.PrintWarning("OpenRouter API 키가 설정되지 않았습니다.")
	platform.PrintInfo("codekiwi.ai에서 발급받으신 API 키를 입력해주세요.")
	fmt.Println()

	// Read API key from user
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("API 키: ")
	apiKey, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}

	// Trim whitespace and newline
	apiKey = strings.TrimSpace(apiKey)

	// Validate input
	if apiKey == "" {
		platform.PrintError("API 키가 입력되지 않았습니다.")
		platform.PrintInfo("API 키는 codekiwi.ai에서 무료로 발급받을 수 있습니다.")
		return fmt.Errorf("API key is required")
	}

	// Create auth config
	config := AuthConfig{}
	config.OpenRouter.Type = "api"
	config.OpenRouter.Key = apiKey

	// Marshal to JSON with indentation
	jsonData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to create auth config: %w", err)
	}

	// Write to file
	if err := os.WriteFile(authFile, jsonData, 0600); err != nil {
		return fmt.Errorf("failed to write auth file: %w", err)
	}

	platform.PrintSuccess(fmt.Sprintf("API 키가 저장되었습니다: %s", authFile))
	fmt.Println()

	return nil
}
