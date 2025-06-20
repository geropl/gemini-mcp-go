# Product Context: gemini-mcp-go

## Problem Statement

AI assistants and development environments that use the Model Context Protocol (MCP) need a standardized way to interact with various backend services. While `mcp-server-gemini` provided a solution for connecting to Google's Gemini API, it was implemented in TypeScript/Node.js. A Go implementation offers potential benefits in performance, static typing, and consistency with other Go-based tools in the ecosystem, such as `linear-mcp-go`.

## Target Audience

The primary users of this server are developers and AI assistants who need to interact with the Gemini API through an MCP-compatible client.

## How It Works

`gemini-mcp-go` will function as a standalone server that:

1.  **Listens for MCP requests:** It will run as a process, listening on stdin for JSON-RPC 2.0 messages that conform to the MCP specification.
2.  **Translates requests:** It will parse incoming MCP requests (e.g., `generate`, `stream`) and translate them into the appropriate calls to the Google Gemini API.
3.  **Handles API communication:** It will manage the HTTP requests to the Gemini API, including authentication with the `GEMINI_API_KEY`.
4.  **Returns responses:** It will format the responses from the Gemini API into MCP-compliant messages and send them back to the client via stdout.
5.  **Manages streaming:** For `stream` requests, it will handle the real-time streaming of response chunks from the Gemini API back to the client.

## User Experience Goals

-   **Seamless Integration:** The server should be easy to install and configure with any MCP-compatible client.
-   **Reliability:** The server should be robust and handle errors gracefully, providing clear feedback to the client.
-   **Performance:** The Go implementation should provide a responsive experience, especially for streaming requests.
-   **Testability:** The server should be thoroughly tested to ensure correctness and prevent regressions.
