package gemini

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// Client is a wrapper around the Gemini API client.
type Client struct {
	genaiClient *genai.GenerativeModel
}

// WithHTTPClient returns an option.ClientOption that sets the HTTP client.
func WithHTTPClient(httpClient *http.Client) option.ClientOption {
	return option.WithHTTPClient(httpClient)
}

// NewClient creates a new Gemini client.
func NewClient(ctx context.Context, apiKey string, opts ...option.ClientOption) (*Client, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("Gemini API key is required")
	}

	allOpts := append([]option.ClientOption{option.WithAPIKey(apiKey)}, opts...)

	client, err := genai.NewClient(ctx, allOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create genai client: %w", err)
	}

	model := client.GenerativeModel("gemini-pro")

	return &Client{
		genaiClient: model,
	}, nil
}

// GenerateContent sends a non-streaming request to the Gemini API.
func (c *Client) GenerateContent(ctx context.Context, prompt string) (string, error) {
	resp, err := c.genaiClient.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", fmt.Errorf("failed to generate content: %w", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no content in response")
	}

	content, ok := resp.Candidates[0].Content.Parts[0].(genai.Text)
	if !ok {
		return "", fmt.Errorf("unexpected response format")
	}

	return string(content), nil
}

// GenerateContentStream sends a streaming request to the Gemini API.
func (c *Client) GenerateContentStream(ctx context.Context, prompt string) (<-chan string, error) {
	iter := c.genaiClient.GenerateContentStream(ctx, genai.Text(prompt))
	ch := make(chan string)

	go func() {
		defer close(ch)
		for {
			resp, err := iter.Next()
			if err != nil {
				// In a real implementation, we'd want to handle this error better,
				// perhaps by sending an error on a separate channel.
				// For now, we'll just close the channel.
				return
			}
			if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
				if content, ok := resp.Candidates[0].Content.Parts[0].(genai.Text); ok {
					ch <- string(content)
				}
			}
		}
	}()

	return ch, nil
}
