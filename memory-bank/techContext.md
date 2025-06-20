# Tech Context: gemini-mcp-go

## Core Technologies

-   **Go:** The primary programming language for the project. We will use a recent version (e.g., 1.23 or later) to take advantage of the latest language features and standard library improvements.
-   **Model Context Protocol (MCP):** The server will implement the MCP specification for communication with AI assistant clients.

## Key Libraries & Dependencies

-   **`github.com/mark3labs/mcp-go`:** This library will provide the core implementation of the Model Context Protocol, handling the low-level details of JSON-RPC 2.0 messaging and the MCP handshake.
-   **`github.com/spf13/cobra`:** A powerful library for creating command-line applications in Go. We will use it to build the `serve` and `setup` commands, providing a clean and extensible CLI interface.
-   **`github.com/google/generative-ai-go`:** The official Go client library for the Google Gemini API. This will be used to interact with the Gemini service for content generation.
-   **`gopkg.in/dnaeon/go-vcr.v4`:** A library for recording and replaying HTTP interactions. This is a critical component of our testing strategy, allowing us to create deterministic and offline tests by mocking the Gemini API.
-   **`github.com/google/go-cmp`:** A library for comparing Go values in tests. It provides more informative diffs than the standard `==` operator, making it easier to debug test failures.

## Development & Testing

-   **Go Toolchain:** Standard Go tools (`go build`, `go test`, `go mod`) will be used for building, testing, and managing dependencies.
-   **Git & GitHub:** The project will be managed using Git for version control and hosted on GitHub for collaboration and CI/CD.
-   **GitHub Actions:** We will set up GitHub Actions for continuous integration to automatically run tests on every push and pull request.

## Configuration

-   **Environment Variables:** The `GEMINI_API_KEY` will be configured via an environment variable to keep it secure.
-   **Command-Line Flags:** Additional options, such as the server port or logging level, will be configurable via command-line flags using Cobra.
