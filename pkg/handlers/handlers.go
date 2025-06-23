package handlers

import (
	"context"
	"fmt"
	"strings"

	"github.com/geropl/gemini-mcp-go/pkg/gemini"
	"github.com/geropl/gemini-mcp-go/pkg/interfaces"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type GenerateParams struct {
	Prompt string `json:"prompt"`
}

func HandleGenerate(geminiClient interfaces.GeminiClient) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		prompt, err := request.RequireString("prompt")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		result, err := geminiClient.GenerateContent(ctx, prompt)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		return mcp.NewToolResultText(result), nil
	}
}

func HandleStream(geminiClient *gemini.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		prompt, err := request.RequireString("prompt")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		stream, err := geminiClient.GenerateContentStream(ctx, prompt)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		progressToken := request.Params.Meta.ProgressToken
		mcpServer := server.ServerFromContext(ctx)
		var fullResponse strings.Builder

		for {
			select {
			case <-ctx.Done():
				return mcp.NewToolResultError("request cancelled"), nil
			case chunk, ok := <-stream:
				if !ok {
					return mcp.NewToolResultText(fullResponse.String()), nil
				}
				fullResponse.WriteString(chunk)
				if progressToken != nil {
					err := mcpServer.SendNotificationToClient(
						ctx,
						"stream/update",
						map[string]any{
							"content":       chunk,
							"progressToken": progressToken,
						},
					)
					if err != nil {
						// Log the error but don't stop the stream
						fmt.Printf("failed to send notification: %v\n", err)
					}
				}
			}
		}
	}
}
