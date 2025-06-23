package handlers_test

import (
	"context"
	"flag"
	"testing"

	"github.com/geropl/gemini-mcp-go/pkg/handlers"
	"github.com/geropl/gemini-mcp-go/pkg/interfaces"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/mcptest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var record = flag.Bool("record", false, "Record new cassettes for tests")

type MockGeminiClient struct {
	responses map[string]string
}

func NewMockGeminiClient() *MockGeminiClient {
	return &MockGeminiClient{
		responses: make(map[string]string),
	}
}

func (m *MockGeminiClient) SetResponse(prompt string, response string) {
	m.responses[prompt] = response
}

func (m *MockGeminiClient) GenerateContent(ctx context.Context, prompt string) (string, error) {
	if response, ok := m.responses[prompt]; ok {
		return response, nil
	}
	return "", nil
}

func TestHandleGenerate(t *testing.T) {
	testCases := []struct {
		name        string
		prompt      string
		expectation string
	}{
		{
			name:        "Capital of France",
			prompt:      "What is the capital of France?",
			expectation: "Paris",
		},
		{
			name:        "Capital of Germany",
			prompt:      "What is the capital of Germany?",
			expectation: "Berlin",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var geminiClient interfaces.GeminiClient
			cassetteName := handlers.CassetteName(tc.name)

			mockClient := NewMockGeminiClient()
			mockClient.SetResponse(tc.prompt, tc.expectation)
			geminiClient = mockClient

			server := mcptest.NewUnstartedServer(t)
			server.AddTool(mcp.NewTool("generate", mcp.WithString("prompt", mcp.Required())), handlers.HandleGenerate(geminiClient))
			require.NoError(t, server.Start(context.Background()))

			client := server.Client()

			ctx := context.Background()

			result, err := client.CallTool(ctx, mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name: "generate",
					Arguments: map[string]any{
						"prompt": tc.prompt,
					},
				},
			})
			require.NoError(t, err)

			if *record {
				err = handlers.WriteCassetteData(cassetteName, tc.prompt, result.Content[0].(mcp.TextContent).Text)
				require.NoError(t, err)
			}

			assert.Equal(t, tc.expectation, result.Content[0].(mcp.TextContent).Text)
		})
	}
}
