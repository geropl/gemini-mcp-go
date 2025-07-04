package handlers

import (
	"context"
	"fmt"
	"strings"

	"github.com/geropl/gemini-mcp-go/pkg/session"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// NewListSessionsHandler creates a new handler for the list_sessions tool.
func NewListSessionsHandler(sm *session.Manager) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		sessions := sm.ListSessions()
		if len(sessions) == 0 {
			return mcp.NewToolResultText("No active sessions"), nil
		}

		var sessionText strings.Builder
		sessionText.WriteString("Active sessions:\n")
		for _, s := range sessions {
			problemSummary := s.ProblemDescription
			if len(problemSummary) > 100 {
				problemSummary = problemSummary[:100] + "..."
			}
			hasCodeContext := "No"
			if s.CodeContext != "" {
				hasCodeContext = "Yes"
			}

			sessionText.WriteString(fmt.Sprintf(
				"- **%s**\n  Messages: %d\n  Created: %s\n  Last used: %s\n  Files attached: %d\n  Code context: %s\n  Problem: %s\n\n",
				s.ID,
				s.MessageCount,
				s.Created.Format("2006-01-02 15:04:05"),
				s.LastUsed.Format("2006-01-02 15:04:05"),
				len(s.ProcessedFiles),
				hasCodeContext,
				problemSummary,
			))
		}

		return mcp.NewToolResultText(sessionText.String()), nil
	}
}

type EndSessionParams struct {
	SessionID string `json:"session_id"`
}

// NewEndSessionHandler creates a new handler for the end_session tool.
func NewEndSessionHandler(sm *session.Manager) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var params EndSessionParams
		if err := request.BindArguments(&params); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("invalid parameters: %v", err)), nil
		}

		if params.SessionID == "" {
			return mcp.NewToolResultError("session_id is a required parameter"), nil
		}

		err := sm.EndSession(params.SessionID)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Session %s has been ended", params.SessionID)), nil
	}
}
