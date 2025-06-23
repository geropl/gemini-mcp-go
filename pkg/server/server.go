package server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/geropl/gemini-mcp-go/pkg/gemini"
	"github.com/geropl/gemini-mcp-go/pkg/handlers"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Run starts the MCP server.
func Run(ctx context.Context, apiKey string) error {
	// Start health check endpoint in a separate goroutine
	go startHealthCheck()

	fmt.Println("Server running. Waiting for MCP client to connect.")

	geminiClient, err := gemini.NewClient(ctx, apiKey)
	if err != nil {
		return fmt.Errorf("failed to create gemini client: %w", err)
	}

	s := server.NewMCPServer(
		"gemini-mcp-go",
		"0.0.1",
		server.WithToolCapabilities(false),
		server.WithRecovery(),
	)

	generateTool := mcp.NewTool("gemini_generate",
		mcp.WithDescription("Generate content using the Gemini API."),
		mcp.WithString("prompt", mcp.Required(), mcp.Description("The prompt to send to Gemini.")),
	)
	s.AddTool(generateTool, handlers.HandleGenerate(geminiClient))

	streamTool := mcp.NewTool("gemini_stream",
		mcp.WithDescription("Stream content from the Gemini API."),
		mcp.WithString("prompt", mcp.Required(), mcp.Description("The prompt to send to Gemini.")),
	)
	s.AddTool(streamTool, handlers.HandleStream(geminiClient))

	return server.ServeStdio(s)
}

// startHealthCheck starts a simple HTTP server for health checks.
func startHealthCheck() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})
	if err := http.ListenAndServe(":3005", nil); err != nil {
		log.Fatalf("Health check server failed: %v", err)
	}
}
