package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

var (
	tool        string
	autoApprove string
	userHomeDir string
)

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Installs and configures the gemini-mcp-go server",
	Long: `This command installs the gemini-mcp-go binary to a standard location
and configures it for use with various AI assistants like Cline, Roo Code,
and Claude Code.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if os.Getenv("GEMINI_API_KEY") == "" {
			return fmt.Errorf("the GEMINI_API_KEY environment variable is not set. Please get your API key from https://makersuite.google.com/app/apikey and set it")
		}

		if err := checkBinary(); err != nil {
			return fmt.Errorf("failed to check binary: %w", err)
		}

		tools := strings.Split(tool, ",")
		for _, t := range tools {
			var output string
			var err error
			switch strings.TrimSpace(t) {
			case "cline":
				output, err = setupCline()
				if err != nil {
					return fmt.Errorf("failed to setup cline: %w", err)
				}
			case "claude-code":
				output, err = setupClaudeCode()
				if err != nil {
					return fmt.Errorf("failed to setup claude-code: %w", err)
				}
			case "roo-code":
				output, err = setupRooCode()
				if err != nil {
					return fmt.Errorf("failed to setup roo-code: %w", err)
				}
			default:
				return fmt.Errorf("unknown tool: %s", t)
			}
			fmt.Print(output)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
	setupCmd.Flags().StringVar(&tool, "tool", "cline", "The AI assistant to configure (e.g., 'cline', 'claude-code')")
	setupCmd.Flags().StringVar(&autoApprove, "auto-approve", "", "A comma-separated list of tool names to auto-approve for execution")
}

func checkBinary() error {
	// Check if the binary is in the PATH
	if _, err := exec.LookPath("gemini-mcp-go"); err == nil {
		fmt.Println("gemini-mcp-go is already in your PATH.")
		return nil
	}

	// Check if the binary is in the standard mcp-servers directory
	home, err := getHomeDir()
	if err != nil {
		return err
	}
	binaryPath := filepath.Join(home, "mcp-servers", "gemini-mcp-go")
	if _, err := os.Stat(binaryPath); err == nil {
		fmt.Printf("gemini-mcp-go is already installed at %s.\n", binaryPath)
		return nil
	}

	// If not found, copy it
	fmt.Println("gemini-mcp-go not found. Installing...")
	return copySelfToBinaryPath(binaryPath)
}

func copySelfToBinaryPath(destPath string) error {
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get current executable path: %w", err)
	}

	input, err := os.ReadFile(exePath)
	if err != nil {
		return fmt.Errorf("failed to read current executable: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	if err := os.WriteFile(destPath, input, 0755); err != nil {
		return fmt.Errorf("failed to write executable to destination: %w", err)
	}

	fmt.Printf("Successfully installed gemini-mcp-go to %s\n", destPath)
	return nil
}

func setupTool(toolName, configPath string) (string, error) {
	var outputBuffer bytes.Buffer
	outputBuffer.WriteString(fmt.Sprintf("Configuring %s...\n", toolName))

	home, err := getHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}
	binaryPath := filepath.Join(home, "mcp-servers", "gemini-mcp-go")

	serverConfig := map[string]interface{}{
		"command": binaryPath,
		"args":    []string{"serve"},
		"env": map[string]string{
			"GEMINI_API_KEY": os.Getenv("GEMINI_API_KEY"),
		},
		"disabled":    false,
		"autoApprove": []string{},
	}

	if autoApprove != "" {
		serverConfig["autoApprove"] = strings.Split(autoApprove, ",")
	}

	file, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Create a new file with the default structure
			file = []byte(`{"mcpServers":{}}`)
		} else {
			return "", fmt.Errorf("failed to read config file: %w", err)
		}
	}

	var config map[string]interface{}
	if err := json.Unmarshal(file, &config); err != nil {
		return "", fmt.Errorf("failed to unmarshal config JSON: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return "", fmt.Errorf("failed to create config directory: %w", err)
	}

	mcpServers, ok := config["mcpServers"].(map[string]interface{})
	if !ok {
		mcpServers = make(map[string]interface{})
	}

	mcpServers["gemini"] = serverConfig
	config["mcpServers"] = mcpServers

	output, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal config JSON: %w", err)
	}

	if err := os.WriteFile(configPath, output, 0644); err != nil {
		return "", fmt.Errorf("failed to write config file: %w", err)
	}

	outputBuffer.WriteString(fmt.Sprintf("Successfully configured %s.\n", toolName))
	return outputBuffer.String(), nil
}

func setupCline() (string, error) {
	home, err := getHomeDir()
	if err != nil {
		return "", err
	}
	var configDir string
	switch runtime.GOOS {
	case "darwin":
		configDir = filepath.Join(home, "Library", "Application Support", "Code", "User", "globalStorage", "saoudrizwan.claude-dev", "settings")
	case "linux":
		configDir = filepath.Join(home, ".vscode-server", "data", "User", "globalStorage", "saoudrizwan.claude-dev", "settings")
	case "windows":
		configDir = filepath.Join(home, "AppData", "Roaming", "Code", "User", "globalStorage", "saoudrizwan.claude-dev", "settings")
	default:
		return "", fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
	return setupTool("Cline", filepath.Join(configDir, "cline_mcp_settings.json"))
}

func setupRooCode() (string, error) {
	home, err := getHomeDir()
	if err != nil {
		return "", err
	}
	var configDir string
	switch runtime.GOOS {
	case "darwin":
		configDir = filepath.Join(home, "Library", "Application Support", "Code", "User", "globalStorage", "rooveterinaryinc.roo-cline", "settings")
	case "linux":
		configDir = filepath.Join(home, ".vscode-server", "data", "User", "globalStorage", "rooveterinaryinc.roo-cline", "settings")
	case "windows":
		configDir = filepath.Join(home, "AppData", "Roaming", "Code", "User", "globalStorage", "rooveterinaryinc.roo-cline", "settings")
	default:
		return "", fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
	return setupTool("Roo Code", filepath.Join(configDir, "cline_mcp_settings.json"))
}

func setupClaudeCode() (string, error) {
	home, err := getHomeDir()
	if err != nil {
		return "", err
	}
	var configPath string
	switch runtime.GOOS {
	case "darwin":
		configPath = filepath.Join(home, "Library", "Application Support", "com.anthropic.claude-code", "claude_code_mcp_settings.json")
	case "linux":
		configPath = filepath.Join(home, ".claude.json")
	case "windows":
		configPath = filepath.Join(home, "AppData", "Roaming", "claude-code", "claude_code_mcp_settings.json")
	default:
		return "", fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
	return setupTool("claude-code", configPath)
}

func getHomeDir() (string, error) {
	if userHomeDir != "" {
		return userHomeDir, nil
	}
	return os.UserHomeDir()
}
