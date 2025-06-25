This document provides a detailed summary of the `mcp-server-gemini` project, intended to guide its reimplementation in another programming language.

### **1. Project Overview**

The `mcp-server-gemini` is a WebSocket-based server that implements the Model Context Protocol (MCP). It acts as a bridge between an MCP-compatible client (like Claude Desktop) and Google's Gemini Pro API. The server listens for JSON-RPC 2.0 requests over WebSockets, translates them into calls to the Gemini API, and returns the results to the client.

### **2. Core Components & Architecture**

The application is built around a few key components:

*   **`MCPServer` (`server.ts`):** This is the main class that sets up and manages the server. Its responsibilities include:
    *   Creating an HTTP server with a `/health` endpoint for basic monitoring.
    *   Attaching a WebSocket server to the HTTP server to handle client connections.
    *   Managing the lifecycle of WebSocket connections: establishing, monitoring for stale connections, and cleaning up on disconnection.
    *   Instantiating the `ProtocolManager` and `MCPHandlers`.
    *   Handling graceful shutdown on `SIGINT` and `SIGTERM` signals.

*   **`MCPHandlers` (`handlers.ts`):** This class is responsible for processing the business logic of incoming MCP requests. It receives requests from the `MCPServer`, validates them, and then interacts with the Gemini API. Key methods include:
    *   `handleInitialize`: Responds with the server's capabilities.
    *   `handleGenerate`: Handles single, non-streaming requests to the Gemini API.
    *   `handleStream`: Manages streaming requests, sending back chunks of the response as they are received from the Gemini API.
    *   `handleCancel`: Cancels in-flight requests.
    *   `handleConfigure`: (Placeholder for) handling configuration updates.
    *   It uses an `EventEmitter` to send streaming responses back to the main server class.

*   **`ProtocolManager` (`protocol.ts`):** This class manages the state of the MCP protocol. It ensures that requests are handled in the correct order (e.g., `initialize` must be called before other methods). It also defines constants for the protocol version, server capabilities, and error codes.

*   **`index.ts`:** This is the entry point of the application. It reads the `GEMINI_API_KEY` from the environment, instantiates the `MCPServer`, and starts it. It contains an initial, simpler implementation of the server, while `server.ts` contains a more robust and feature-complete version. A reimplementation should focus on the `server.ts` architecture.

### **3. Configuration**

The server is configured primarily through environment variables:

*   `GEMINI_API_KEY` (required): The API key for accessing the Google Gemini API.
*   `PORT` (optional, defaults to `3005`): The port on which the server will listen.
*   `DEBUG` (optional, defaults to `false`): If set to `'true'`, enables debug logging to the console.

### **4. Communication Protocol**

*   **Transport:** The server uses WebSockets for communication.
*   **Payload Format:** All messages are JSON-RPC 2.0 objects.
*   **MCP Methods:** The server implements the following standard MCP methods:
    *   `initialize`: To set up the connection and exchange capabilities.
    *   `generate`: For a single request-response cycle.
    *   `stream`: For streaming responses.
    *   `cancel`: To cancel an ongoing request.
    *   `configure`: To adjust server settings.
*   **Health Check:** A standard HTTP `GET` request to `/health` on the server's port will return a JSON object with the server's status, uptime, and number of active connections.

### **5. Request Handling Flow**

1.  A client connects to the WebSocket server.
2.  The server creates a `ConnectionState` object for the client to track its status.
3.  The client sends an `initialize` request. The `ProtocolManager` marks the connection as initialized.
4.  The client sends a `generate` or `stream` request.
5.  The `MCPServer` forwards the request to `MCPHandlers`.
6.  `MCPHandlers` validates the request parameters.
7.  `MCPHandlers` calls the `generateContent` or `generateContentStream` method of the Gemini API.
8.  For `generate`, the full response is returned in a single JSON-RPC response.
9.  For `stream`, multiple `stream` notifications are sent back to the client as chunks are received, followed by a final notification with `done: true`.
10. `MCPHandlers` uses `AbortController` to manage cancellation of requests.

### **6. State Management**

*   **Connection State:** The `MCPServer` maintains a `Map` of connected clients, where each key is a `WebSocket` object and the value is a `ConnectionState` object containing information like IP address, connection time, and active request IDs.
*   **Protocol State:** The `ProtocolManager` tracks whether a connection has been initialized and whether a shutdown has been requested.
*   **Active Requests:** The `MCPHandlers` class maintains a `Map` of active requests (`requestId` -> `AbortController`) to allow for cancellation.

### **7. Error Handling**

*   The server defines a set of standard JSON-RPC and custom error codes in `protocol.ts`.
*   Errors are sent back to the client as a standard JSON-RPC error object.
*   The server logs detailed error information to the console, including the connection state at the time of the error.

### **8. Dependencies**

A reimplementation would need equivalents for the following key Node.js packages:

*   `ws`: For WebSocket server functionality.
*   `@google/generative-ai`: The official Google Gemini API client library.
*   A standard HTTP library for the health check endpoint.

This summary should provide a solid foundation for reimplementing the `mcp-server-gemini` in any language that supports WebSockets and has a library for interacting with the Gemini API.
