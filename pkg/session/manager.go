package session

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/generative-ai-go/genai"
	"github.com/google/uuid"
	"google.golang.org/api/option"
)

const (
	// DefaultSystemPrompt is the default system prompt for Gemini.
	DefaultSystemPrompt = `You are an expert coding assistant helping Claude (another AI) solve complex programming problems.

Your role:
- Provide clear, practical solutions with working code examples
- Explain your reasoning concisely but thoroughly
- Focus on best practices, security, and maintainability
- Suggest optimizations when relevant
- Point out potential issues or edge cases
- Use the specific technologies and frameworks shown in the provided code context

Response guidelines:
- Start with a brief summary of your approach
- Provide complete, runnable code examples when possible
- Explain key concepts or non-obvious implementations
- Suggest testing strategies when appropriate
- Be direct and actionable - Claude needs specific guidance to help the user
- If you need additional context to provide a solid answer, ask Claude specific clarifying questions about:
  - Requirements or constraints not mentioned
  - Preferred approaches or technologies
  - Error messages or specific behaviors
  - Environment details or deployment context
  - Performance requirements or scale considerations

Remember: You're consulting with another AI to help a human developer, so be precise and comprehensive in your technical advice.`

	// SessionTTL is the time-to-live for a session.
	SessionTTL = 1 * time.Hour
)

// Manager handles the lifecycle of sessions.
type Manager struct {
	sessions   map[string]*Session
	mu         sync.RWMutex
	genaiClient *genai.Client
}

// NewManager creates a new session manager.
func NewManager(ctx context.Context, apiKey string) (*Manager, error) {
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create genai client: %w", err)
	}

	return &Manager{
		sessions:   make(map[string]*Session),
		genaiClient: client,
	}, nil
}

// GetOrCreateSession gets an existing session or creates a new one.
func (m *Manager) GetOrCreateSession(sessionID string) (*Session, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if sessionID != "" {
		if session, ok := m.sessions[sessionID]; ok {
			session.Touch()
			return session, nil
		}
	} else {
		sessionID = uuid.New().String()
	}

	// Create new session
	model := m.genaiClient.GenerativeModel("gemini-1.5-pro-latest")
	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{genai.Text(DefaultSystemPrompt)},
	}
	chat := model.StartChat()
	chat.History = []*genai.Content{}

	session := NewSession(sessionID, chat)
	m.sessions[sessionID] = session

	return session, nil
}

// ListSessions lists all active sessions.
func (m *Manager) ListSessions() []*Session {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var sessions []*Session
	for _, session := range m.sessions {
		sessions = append(sessions, session)
	}
	return sessions
}

// EndSession ends a specific session.
func (m *Manager) EndSession(sessionID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	session, ok := m.sessions[sessionID]
	if !ok {
		return fmt.Errorf("session not found: %s", sessionID)
	}

	m.cleanupSessionFiles(context.Background(), session)
	delete(m.sessions, sessionID)
	return nil
}

// Close cleans up the manager's resources.
func (m *Manager) Close() {
	m.genaiClient.Close()
}

// StartCleanupRoutine starts a goroutine to periodically clean up expired sessions.
func (m *Manager) StartCleanupRoutine(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Minute)
	go func() {
		for {
			select {
			case <-ticker.C:
				m.cleanupExpiredSessions()
			case <-ctx.Done():
				ticker.Stop()
				return
			}
		}
	}()
}

func (m *Manager) cleanupExpiredSessions() {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	for id, session := range m.sessions {
		if now.Sub(session.LastUsed) > SessionTTL {
			m.cleanupSessionFiles(context.Background(), session)
			delete(m.sessions, id)
		}
	}
}
