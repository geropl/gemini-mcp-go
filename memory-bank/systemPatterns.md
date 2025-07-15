# System Patterns: gemini-mcp-go

## High-Level Architecture

The `gemini-mcp-go` server implements a stateful, session-based architecture for complex multi-turn conversations with Google's Gemini API.

```mermaid
graph TD
    subgraph "gemini-mcp-go"
        A[main.go] --> B[cmd/serve.go]
        B --> C{Server (pkg/server)}
        C -- Manages --> SM{Session Manager (pkg/session)}
        C -- MCP Requests --> D{Handlers (pkg/handlers)}
        D -- Uses --> SM
        SM -- Manages --> S[Sessions with Chat History]
        SM -- File Processing --> F[File Upload & Management]
        SM -- Gemini API --> G[generative-ai-go client]
        C -- Protocol --> H[mcp-go library]
    end

    subgraph "External Dependencies"
        I[Google Gemini API]
        J[MCP Client (e.g., Claude Desktop)]
        K[Local File System]
    end

    G --> I
    J --> C
    F --> K
```

## Key Architectural Patterns

### 1. Session-Based State Management

-   **Persistent Conversations:** Each session maintains chat history and context across multiple queries.
-   **File Context Caching:** Uploaded files are cached per session to reduce token usage on follow-up questions.
-   **Automatic Cleanup:** Sessions expire after 1 hour of inactivity, with automatic file deletion from Gemini API.

### 2. Modular Package Structure

-   **`cmd`:** Command-line interface using Cobra framework.
    -   `root.go` defines the root command and global flags.
    -   `serve.go` starts the MCP server with session management.
    -   `setup.go` provides guided configuration for AI assistant integration.

-   **`pkg/server`:** Core MCP server implementation.
    -   Integrates with `mcp-go` library for protocol handling.
    -   Manages server lifecycle and tool registration.
    -   Provides health check endpoint.

-   **`pkg/session`:** Stateful conversation management.
    -   `Manager` handles session lifecycle, cleanup, and file processing.
    -   `Session` maintains conversation state and chat history.
    -   Thread-safe operations with proper locking.

-   **`pkg/handlers`:** Business logic for MCP tools.
    -   `consult_gemini.go` implements the primary conversation tool.
    -   `session_tools.go` implements session management tools.
    -   Stateless handlers that use session manager for persistence.

### 3. Resource Management

-   **File Upload Lifecycle:** Files are uploaded to Gemini API, tracked per session, and automatically deleted on cleanup.
-   **Memory Management:** Sessions are stored in memory with TTL-based expiration.
-   **Concurrent Safety:** Thread-safe session operations for handling multiple simultaneous conversations.

### 4. Testing Strategy

-   **Unit Tests:** Comprehensive tests for session manager with mocked Gemini client implemented.
-   **Mock Infrastructure:** Complete mock implementation for Gemini API client operations.
-   **Dependency Injection:** Session manager is injected into handlers for testability.
-   **Concurrency Testing:** Thread safety verified with concurrent access patterns.
-   **Clean Test Foundation:** Existing integration tests removed, new foundation prepared.

### 5. Configuration Management

-   **Environment Variables:** `GEMINI_API_KEY` loaded from environment.
-   **Guided Setup:** Interactive setup command for AI assistant configuration.
-   **Multiple AI Assistant Support:** Configuration for Cline, Claude Code, and Roo Code.

This architecture provides a robust foundation for stateful AI conversations with proper resource management and scalability.
