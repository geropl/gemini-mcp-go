package handlers

import (
	"context"
	"fmt"
	"strings"

	"github.com/geropl/gemini-mcp-go/pkg/session"
	"github.com/google/generative-ai-go/genai"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type ConsultGeminiParams struct {
	SessionID          string            `json:"session_id"`
	ProblemDescription string            `json:"problem_description"`
	CodeContext        string            `json:"code_context"`
	AttachedFiles      []string          `json:"attached_files"`
	FileDescriptions   map[string]string `json:"file_descriptions"`
	SpecificQuestion   string            `json:"specific_question"`
	AdditionalContext  string            `json:"additional_context"`
	PreferredApproach  string            `json:"preferred_approach"`
}

func NewConsultGeminiHandler(sm *session.Manager) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var params ConsultGeminiParams
		if err := request.BindArguments(&params); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("invalid parameters: %v", err)), nil
		}

		if params.SpecificQuestion == "" {
			return mcp.NewToolResultError("specific_question is a required parameter"), nil
		}

		result, err := handleConsultGeminiLogic(ctx, sm, &params)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(result), nil
	}
}

func handleConsultGeminiLogic(ctx context.Context, sm *session.Manager, params *ConsultGeminiParams) (string, error) {
	s, err := sm.GetOrCreateSession(params.SessionID)
	if err != nil {
		return "", fmt.Errorf("failed to get or create session: %w", err)
	}

	if s.MessageCount == 0 {
		if params.ProblemDescription == "" {
			return "", fmt.Errorf("problem_description is required for new sessions")
		}
		if params.CodeContext == "" && len(params.AttachedFiles) == 0 {
			return "", fmt.Errorf("either code_context or attached_files is required for new sessions")
		}

		s.ProblemDescription = params.ProblemDescription
		s.CodeContext = params.CodeContext

		var contextParts []string
		contextParts = append(contextParts, fmt.Sprintf("I'm Claude, an AI assistant, and I need your help with a complex coding problem. Here's the context:\n\n**Problem Description:**\n%s", params.ProblemDescription))

		if params.CodeContext != "" {
			contextParts = append(contextParts, fmt.Sprintf("\n**Code Context:**\n%s", params.CodeContext))
		}

		if len(params.AttachedFiles) > 0 {
			contextParts = append(contextParts, "\n**Attached Files:**")
			for _, filePath := range params.AttachedFiles {
				fileInfo, err := sm.ProcessFile(ctx, s, filePath)
				if err != nil {
					contextParts = append(contextParts, fmt.Sprintf("\n- %s (failed to upload: %v)", filePath, err))
					continue
				}
				description := ""
				if desc, ok := params.FileDescriptions[filePath]; ok {
					description = " - " + desc
				}
				contextParts = append(contextParts, fmt.Sprintf("\n- %s%s", fileInfo.Name, description))
			}
		}
		contextParts = append(contextParts, "\n\nPlease help me solve this problem. I may have follow-up questions, so please maintain context throughout our conversation.")

		initialPrompt := strings.Join(contextParts, "")
		var contentParts []genai.Part
		contentParts = append(contentParts, genai.Text(initialPrompt))

		for _, file := range s.ProcessedFiles {
			contentParts = append(contentParts, genai.FileData{MIMEType: file.MIMEType, URI: file.URI})
		}

		_, err := s.Chat.SendMessage(ctx, contentParts...)
		if err != nil {
			return "", fmt.Errorf("failed to send initial message to Gemini: %w", err)
		}
		s.MessageCount++
	}

	var questionParts []string
	questionParts = append(questionParts, fmt.Sprintf("**Question:** %s", params.SpecificQuestion))
	if params.AdditionalContext != "" {
		questionParts = append(questionParts, fmt.Sprintf("\n\n**Additional Context/Updates:**\n%s", params.AdditionalContext))
	}
	if params.PreferredApproach != "" && params.PreferredApproach != "follow-up" {
		questionParts = append(questionParts, fmt.Sprintf("\n\n**Type of Help Needed:** %s", params.PreferredApproach))
	}
	questionPrompt := strings.Join(questionParts, "")

	resp, err := s.Chat.SendMessage(ctx, genai.Text(questionPrompt))
	if err != nil {
		return "", fmt.Errorf("failed to send message to Gemini: %w", err)
	}
	s.MessageCount++

	var responseText strings.Builder
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				if txt, ok := part.(genai.Text); ok {
					responseText.WriteString(string(txt))
				}
			}
		}
	}

	return fmt.Sprintf("**Session ID:** %s\n**Message #%d**\n\n%s\n\n---\n*Use session_id: \"%s\" for follow-up questions*", s.ID, s.MessageCount, responseText.String(), s.ID), nil
}
