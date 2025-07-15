# Project Brief: gemini-mcp-go

## Objective

Create a powerful, stateful MCP server for Google's Gemini API that enables complex, multi-turn conversations with file attachment support. This is a Go reimplementation and enhancement of the `mcp-gemini-assistant` Python project.

## Key Requirements

1.  **Stateful Architecture:** Session-based conversation management with context persistence across multiple queries.
2.  **File Integration:** Support for reading and attaching local code files to conversations, with automatic upload to Gemini API.
3.  **Session Management:** Automatic cleanup, TTL-based expiration, and resource management for uploaded files.
4.  **Technology:** Built in Go using the `mcp-go` library and `generative-ai-go` for Gemini integration.
5.  **Testing:** Comprehensive unit tests for session management with mocked dependencies.
6.  **Configuration:** Environment variable and command-line configuration, with guided setup command.
7.  **User Experience:** Simple installation and configuration process for AI assistant integration.

## Scope

-   **Core Architecture**: Stateful, session-based MCP server with conversation context management.
-   **Primary Tool**: `consult_gemini` for complex, multi-turn conversations with file attachment support.
-   **Session Management**: `list_sessions` and `end_session` tools for managing active conversations.
-   **Commands**: `serve` command to run the server, `setup` command for easy installation.
-   **Testing**: Comprehensive unit tests for session management, clean foundation for integration tests.
-   **Documentation**: Complete documentation reflecting the new stateful architecture.
