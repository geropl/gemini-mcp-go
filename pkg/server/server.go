package server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/geropl/gemini-mcp-go/pkg/handlers"
	"github.com/geropl/gemini-mcp-go/pkg/session"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Run starts the MCP server.
func Run(ctx context.Context, apiKey string) error {
	// Start health check endpoint in a separate goroutine
	go startHealthCheck()

	fmt.Println("Server running. Waiting for MCP client to connect.")

	sm, err := session.NewManager(ctx, apiKey)
	if err != nil {
		return fmt.Errorf("failed to create session manager: %w", err)
	}
	defer sm.Close()
	sm.StartCleanupRoutine(ctx)

	s := server.NewMCPServer(
		"gemini-mcp-go",
		"0.1.0", // Version bump for new features
		server.WithToolCapabilities(false),
		server.WithRecovery(),
	)

	consultTool := mcp.NewTool("consult_gemini",
		mcp.WithDescription("Start or continue a conversation with Gemini about complex coding problems."),
		mcp.WithString("session_id", mcp.Description("Optional session ID to continue a previous conversation")),
		mcp.WithString("problem_description", mcp.Description("Detailed description of the coding problem (required for new sessions)")),
		mcp.WithString("code_context", mcp.Description("All relevant code - will be cached for the session (required for new sessions)")),
		mcp.WithArray("attached_files", mcp.Description("Array of file paths to upload and attach to the conversation")),
		mcp.WithObject("file_descriptions", mcp.Description("Optional object mapping file paths to descriptions")),
		mcp.WithString("specific_question", mcp.Required(), mcp.Description("The specific question you want answered")),
		mcp.WithString("additional_context", mcp.Description("Additional context, updates, or what changed since last question")),
		mcp.WithString("preferred_approach", mcp.Description("Type of assistance needed (solution, review, debug, optimize, explain, follow-up)")),
	)
	s.AddTool(consultTool, handlers.NewConsultGeminiHandler(sm))

	listSessionsTool := mcp.NewTool("list_sessions",
		mcp.WithDescription("List all active Gemini consultation sessions."),
	)
	s.AddTool(listSessionsTool, handlers.NewListSessionsHandler(sm))

	endSessionTool := mcp.NewTool("end_session",
		mcp.WithDescription("End a specific Gemini consultation session to free up memory."),
		mcp.WithString("session_id", mcp.Required(), mcp.Description("The session ID to end.")),
	)
	s.AddTool(endSessionTool, handlers.NewEndSessionHandler(sm))

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
