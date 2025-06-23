# Gemini MCP Go

This project is a Go implementation of an MCP server for Google's Gemini API. It acts as a bridge between an MCP client and the Gemini API, allowing AI assistants and other tools to interact with Gemini in a standardized way.

## Features

-   **Gemini Integration:** Provides access to the Gemini API for content generation and streaming.
-   **MCP Compliant:** Implements the Model Context Protocol for seamless integration with MCP clients.
-   **Easy to Use:** Includes a `setup` command to help with initial configuration.
-   **Health Check:** Provides a `/health` endpoint for monitoring the server's status.

## Installation

To install the Gemini MCP server, you need to have Go installed on your system. You can then clone this repository and build the server using the following commands:

```bash
git clone https://github.com/geropl/gemini-mcp-go.git
cd gemini-mcp-go
go build
```

## Usage

### Setup

Before you can use the server, you need to set your Gemini API key as an environment variable. You can use the `setup` command to check if the environment variable is set and to get instructions on how to set it.

```bash
./gemini-mcp-go setup
```

### Serve

To start the server, use the `serve` command. You **must** set the `GEMINI_API_KEY` environment variable before running the server.

```bash
export GEMINI_API_KEY=YOUR_API_KEY
./gemini-mcp-go serve
```

## Testing

To run the tests, use the following command:

```bash
go test ./...
