# Project Brief: gemini-mcp-go

## Objective

The primary goal of this project is to reimplement the `mcp-server-gemini` project in Go, creating a new MCP server named `gemini-mcp-go`.

This project will leverage the architecture, patterns, and best practices from the existing `linear-mcp-go` project to ensure the new server is robust, maintainable, and testable.

## Key Requirements

1.  **Functionality:** The new server must replicate the core functionality of `mcp-server-gemini`, acting as a bridge between an MCP client and the Google Gemini API.
2.  **Technology:** The server will be built in Go, using the `mcp-go` library for the core MCP implementation.
3.  **Architecture:** The project will adopt the modular structure of `linear-mcp-go`, with distinct packages for command-line handling (`cmd`), core application logic (`pkg`), and test data (`testdata`).
4.  **Testing:** A comprehensive testing strategy will be implemented using golden files for output validation, mirroring the approach in `linear-mcp-go`.
5.  **Configuration:** The server will be configurable via environment variables and command-line flags, particularly for the `GEMINI_API_KEY`.
6.  **User Experience:** The server will include a `setup` command to simplify installation and configuration for end-users.

## Scope

-   Implement `initialize`, `generate`, `stream`, and `cancel` MCP methods.
-   Create a `serve` command to run the server.
-   Create a `setup` command for easy installation.
-   Develop a suite of tests with high coverage.
-   Document the project in a `README.md` and maintain a `memory-bank` for development tracking.
