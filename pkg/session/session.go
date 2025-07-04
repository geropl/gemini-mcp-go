package session

import (
	"sync"
	"time"

	"github.com/google/generative-ai-go/genai"
)

// ProcessedFile holds information about a file that has been uploaded to Gemini.
type ProcessedFile struct {
	Name     string `json:"name"`
	Path     string `json:"path"`
	MIMEType string `json:"mime_type"`
	URI      string `json:"uri"`
}

// Session represents a single, stateful conversation with the Gemini API.
type Session struct {
	ID                 string                    `json:"id"`
	Chat               *genai.ChatSession        `json:"-"` // Ignore in JSON responses
	Created            time.Time                 `json:"created"`
	LastUsed           time.Time                 `json:"last_used"`
	MessageCount       int                       `json:"message_count"`
	ProblemDescription string                    `json:"problem_description"`
	CodeContext        string                    `json:"code_context"`
	ProcessedFiles     map[string]*ProcessedFile `json:"processed_files"`
	mu                 sync.RWMutex
}

// NewSession creates a new session object.
func NewSession(id string, chat *genai.ChatSession) *Session {
	return &Session{
		ID:             id,
		Chat:           chat,
		Created:        time.Now(),
		LastUsed:       time.Now(),
		MessageCount:   0,
		ProcessedFiles: make(map[string]*ProcessedFile),
	}
}

// Touch updates the LastUsed timestamp of the session.
func (s *Session) Touch() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.LastUsed = time.Now()
}

// AddFile adds a processed file to the session's collection.
func (s *Session) AddFile(file *ProcessedFile) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ProcessedFiles[file.Path] = file
}

// GetFile retrieves a processed file from the session.
func (s *Session) GetFile(filePath string) (*ProcessedFile, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	file, ok := s.ProcessedFiles[filePath]
	return file, ok
}
