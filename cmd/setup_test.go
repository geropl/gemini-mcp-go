package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSetupCmd(t *testing.T) {
	clineConfigPath := "__TMP_DIR__/.vscode-server/data/User/globalStorage/saoudrizwan.claude-dev/settings/cline_mcp_settings.json"
	claudeConfigPath := "__TMP_DIR__/.claude.json"

	testCases := []struct {
		name           string
		args           []string
		env            map[string]string
		initialFiles   map[string]string
		expectedStdout string
		expectedFiles  map[string]string
		expectErr      bool
	}{
		{
			name:      "No API Key",
			args:      []string{"setup"},
			expectErr: true,
		},
		{
			name: "Cline Setup, New Config",
			args: []string{"setup", "--tool=cline"},
			env: map[string]string{
				"GEMINI_API_KEY": "test-api-key",
			},
			expectedStdout: `gemini-mcp-go not found. Installing...
Successfully installed gemini-mcp-go to __TMP_DIR__/mcp-servers/gemini-mcp-go
Configuring Cline...
Successfully configured Cline.
`,
			expectedFiles: map[string]string{
				clineConfigPath: `{
  "mcpServers": {
    "gemini": {
      "args": [
        "serve"
      ],
      "autoApprove": [],
      "command": "__TMP_DIR__/mcp-servers/gemini-mcp-go",
      "disabled": false,
      "env": {
        "GEMINI_API_KEY": "test-api-key"
      }
    }
  }
}`,
			},
		},
		{
			name: "Claude Code Setup, Existing Config",
			args: []string{"setup", "--tool=claude-code"},
			env: map[string]string{
				"GEMINI_API_KEY": "test-api-key",
			},
			initialFiles: map[string]string{
				claudeConfigPath: `{"mcpServers":{"another-server":{}}}`,
			},
			expectedStdout: `gemini-mcp-go not found. Installing...
Successfully installed gemini-mcp-go to __TMP_DIR__/mcp-servers/gemini-mcp-go
Configuring claude-code...
Successfully configured claude-code.
`,
			expectedFiles: map[string]string{
				claudeConfigPath: `{
  "mcpServers": {
    "another-server": {},
    "gemini": {
      "args": [
        "serve"
      ],
      "autoApprove": [],
      "command": "__TMP_DIR__/mcp-servers/gemini-mcp-go",
      "disabled": false,
      "env": {
        "GEMINI_API_KEY": "test-api-key"
      }
    }
  }
}`,
			},
		},
		{
			name: "Auto Approve",
			args: []string{"setup", "--tool=cline", "--auto-approve=tool1,tool2"},
			env: map[string]string{
				"GEMINI_API_KEY": "test-api-key",
			},
			expectedStdout: `gemini-mcp-go not found. Installing...
Successfully installed gemini-mcp-go to __TMP_DIR__/mcp-servers/gemini-mcp-go
Configuring Cline...
Successfully configured Cline.
`,
			expectedFiles: map[string]string{
				clineConfigPath: `{
  "mcpServers": {
    "gemini": {
      "args": [
        "serve"
      ],
      "autoApprove": [
        "tool1",
        "tool2"
      ],
      "command": "__TMP_DIR__/mcp-servers/gemini-mcp-go",
      "disabled": false,
      "env": {
        "GEMINI_API_KEY": "test-api-key"
      }
    }
  }
}`,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a temporary directory for the test
			tmpDir, err := os.MkdirTemp("", "gemini-mcp-go-test")
			if err != nil {
				t.Fatalf("Failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tmpDir)

			// Set userHomeDir to the temp dir for testing
			userHomeDir = tmpDir
			defer func() { userHomeDir = "" }()

			// Set up the environment
			originalEnv := make(map[string]string)
			for k, v := range tc.env {
				originalEnv[k] = os.Getenv(k)
				os.Setenv(k, v)
			}
			defer func() {
				for k, v := range originalEnv {
					os.Setenv(k, v)
				}
			}()

			// Create initial files
			for path, content := range tc.initialFiles {
				fullPath := strings.ReplaceAll(path, "__TMP_DIR__", tmpDir)
				if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
					t.Fatalf("Failed to create dir for initial file: %v", err)
				}
				if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
					t.Fatalf("Failed to write initial file: %v", err)
				}
			}

			// Capture stdout and stderr
			oldStdout := os.Stdout
			oldStderr := os.Stderr
			r, w, _ := os.Pipe()
			os.Stdout = w
			os.Stderr = w

			// Execute the command
			rootCmd.SetArgs(tc.args)
			err = rootCmd.Execute()
			rootCmd.SetArgs([]string{}) // Reset args

			// Restore stdout and stderr
			w.Close()
			os.Stdout = oldStdout
			os.Stderr = oldStderr
			var output bytes.Buffer
			output.ReadFrom(r)

			// Check for expected error
			if tc.expectErr {
				if err == nil {
					t.Errorf("Expected an error, but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			// Compare output with expected stdout
			if tc.expectedStdout != "" {
				actualOutput := strings.ReplaceAll(output.String(), tmpDir, "__TMP_DIR__")
				if diff := cmp.Diff(tc.expectedStdout, actualOutput); diff != "" {
					t.Errorf("Output mismatch (-want +got):\n%s", diff)
				}
			}

			// Compare created files with expected files
			for path, expectedContent := range tc.expectedFiles {
				fullPath := strings.ReplaceAll(path, "__TMP_DIR__", tmpDir)
				actual, err := os.ReadFile(fullPath)
				if err != nil {
					t.Fatalf("Failed to read actual file %s: %v", fullPath, err)
				}

				actualContent := strings.ReplaceAll(string(actual), tmpDir, "__TMP_DIR__")
				if diff := cmp.Diff(strings.TrimSpace(expectedContent), strings.TrimSpace(actualContent)); diff != "" {
					t.Errorf("File content mismatch for %s (-want +got):\n%s", path, diff)
				}
			}
		})
	}
}
