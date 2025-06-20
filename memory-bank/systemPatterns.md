# System Patterns: gemini-mcp-go

## High-Level Architecture

The `gemini-mcp-go` server will follow a modular, layered architecture inspired by `linear-mcp-go`. This design promotes separation of concerns, making the application easier to understand, test, and maintain.

```mermaid
graph TD
    subgraph "gemini-mcp-go"
        A[main.go] --> B[cmd/serve.go]
        B --> C{Server (pkg/server)}
        C -- MCP Requests --> D{Handlers (pkg/handlers)}
        D -- Gemini API Calls --> E{Gemini Client (pkg/gemini)}
        C -- Manages --> F[mcp-go library]
    end

    E -- HTTP Requests --> G[Google Gemini API]
    F -- Stdin/Stdout --> H[AI Assistant]

    subgraph "Testing"
        I[Tests] -- Uses --> J[go-vcr]
        J -- Records/Replays --> G
        I -- Compares against --> K[Golden Files]
    end
```

## Key Architectural Patterns

### 1. Command-Line Interface (CLI) with Cobra

-   **`cmd` Package:** All CLI-related code will reside in the `cmd/` directory.
-   **`root.go`:** This file will define the root command.
-   **`serve.go` & `setup.go`:** Each subcommand (`serve`, `setup`) will have its own file, promoting a clean and organized command structure.

### 2. Modular Packages (`pkg`)

-   **`pkg/server`:** This package will contain the core `Server` struct, responsible for:
    -   Initializing the `mcp-go` session.
    -   Managing the server lifecycle (startup, shutdown).
    -   Routing incoming MCP requests to the appropriate handlers.
    -   Implementing the `/health` check endpoint.
-   **`pkg/handlers`:** This package will implement the business logic for each MCP method.
    -   Each handler (e.g., `handleGenerate`, `handleStream`) will be responsible for a single MCP method.
    -   Handlers will be stateless and receive all necessary context from the `Server`.
-   **`pkg/gemini`:** This package will act as a client wrapper for the Google Gemini API.
    -   It will abstract away the details of the `generative-ai-go` library.
    -   It will handle API authentication, request formatting, and response parsing.
    -   This isolation is crucial for mocking the API during testing.

### 3. Test-Driven Development (TDD) with API Mocking

-   **`go-vcr`:** We will use `go-vcr` to record real HTTP interactions with the Gemini API and save them as "cassettes" in the `testdata/fixtures/` directory.
-   **Offline Testing:** Subsequent test runs will use these cassettes to replay the recorded responses, eliminating the need for a live API connection. This makes tests fast, deterministic, and independent of network conditions.
-   **Golden Files:** For each test case, the expected output will be stored in a `.golden` file in `testdata/golden/`. Tests will compare the actual output against these golden files to ensure correctness.
-   **Test Structure:** Tests will be placed alongside the code they are testing (e.g., `pkg/handlers/generate_test.go`).

### 4. Configuration Management

-   **Environment Variables:** Sensitive information like the `GEMINI_API_KEY` will be loaded from the environment.
-   **Command-Line Flags:** Non-sensitive configuration options will be exposed as flags on the `serve` and `setup` commands.

This architecture provides a solid foundation for building a high-quality, production-ready MCP server.
