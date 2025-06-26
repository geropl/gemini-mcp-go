# Design Doc: 0001 - Setup Command

## Objective

The primary goal is to create a `setup` command for the `gemini-mcp-go` server. This command will automate the process of installing the server binary and configuring it for use with various AI assistants, such as Cline, Roo Code, and Claude Code. The functionality will be modeled on the setup command from `linear-mcp-go`, but adapted to the specific needs of the `gemini-mcp-go` server, which requires an API key.

## Implementation Plan

The implementation will be structured across the following steps, ensuring a modular and maintainable design.

### 1. Command and Flag Definition

The `setup` command and its associated flags will be defined in the existing `cmd/setup.go` file.

*   **Command:** A new `cobra.Command` named `setupCmd`.
*   **Flags:**
    *   `--tool` (string): Specifies the target AI assistant(s) (e.g., "cline", "claude-code"). Defaults to "cline".
    *   `--auto-approve` (string): A comma-separated list of tool names that the user wants to auto-approve for execution.
    *   `--project-path` (string): For `claude-code`, specifies project-scoped configuration, allowing the server to be active only for certain projects.

### 2. Core `Run` Logic

The main function of the command will orchestrate the setup process:

1.  **Input Validation:** The command will check if the `GEMINI_API_KEY` environment variable is set. If not, it will guide the user on how to set it before proceeding.
2.  **Binary Installation:** It will reuse the installation logic from `linear-mcp-go`:
    *   Check if the `gemini-mcp-go` binary is already in the system's `PATH` or in the standard `~/mcp-servers/` directory.
    *   If the binary is not found, it will copy the currently running executable to `~/mcp-servers/gemini-mcp-go` and ensure it has execute permissions.
3.  **Configuration Dispatch:** Based on the `--tool` flag, it will call the appropriate setup function (e.g., `setupCline`, `setupClaudeCode`).

### 3. Helper Functions

To keep the code clean and organized, several helper functions will be implemented, adapted from the `linear-mcp-go` project:

*   `checkBinary()` & `copySelfToBinaryPath()`: These functions will manage the installation of the `gemini-mcp-go` executable.
*   `setupTool()`: A generic function that handles the logic of reading, modifying, and writing the JSON configuration file for a given AI assistant.
*   `setupCline()`, `setupRooCode()`, `setupClaudeCode()`: These will be thin wrappers that determine the correct configuration file path for each specific tool and OS, then pass the necessary details to the generic `setupTool` function.

### 4. Configuration File Modification

The core of the setup process is to correctly modify the AI assistant's configuration file. The `setupTool` function will generate a JSON object for the `gemini` server and merge it into the `mcpServers` section of the target configuration file.

An example of the configuration that will be generated for `cline_mcp_settings.json` is:

```json
{
  "mcpServers": {
    "gemini": {
      "command": "/home/vscode/mcp-servers/gemini-mcp-go",
      "args": [
        "serve"
      ],
      "env": {
        "GEMINI_API_KEY": "your_gemini_api_key_from_env"
      },
      "disabled": false,
      "autoApprove": []
    }
  }
}
