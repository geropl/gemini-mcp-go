package session

import (
	"context"
	"fmt"
	"log"
	"mime"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/generative-ai-go/genai"
)

var mimeTypeMap = map[string]string{
	".go":    "text/x-go",
	".jsx":   "text/javascript",
	".tsx":   "text/typescript",
	".ts":    "text/typescript",
	".vue":   "text/html",
	".svelte": "text/html",
	".md":    "text/markdown",
	".json":  "application/json",
	".py":    "text/x-python",
	".js":    "text/javascript",
	".css":   "text/css",
	".html":  "text/html",
	".xml":   "text/xml",
	".yaml":  "text/yaml",
	".yml":   "text/yaml",
	".toml":  "text/plain",
	".ini":   "text/plain",
	".cfg":   "text/plain",
	".conf":  "text/plain",
	".sh":    "text/x-shellscript",
	".bat":   "text/plain",
	".sql":   "text/x-sql",
}

func guessMIMEType(filePath string) string {
	ext := filepath.Ext(filePath)
	if mimeType, ok := mimeTypeMap[ext]; ok {
		return mimeType
	}
	// Fallback to system's mime package
	if mimeType := mime.TypeByExtension(ext); mimeType != "" {
		return mimeType
	}
	return "text/plain"
}

// ProcessFile uploads a file to the Gemini API and returns its metadata.
func (m *Manager) ProcessFile(ctx context.Context, session *Session, filePath string) (*ProcessedFile, error) {
	if file, ok := session.GetFile(filePath); ok {
		return file, nil // Already processed
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("file not found: %s", filePath)
	}

	fileName := filepath.Base(filePath)
	mimeType := guessMIMEType(filePath)

	fileReader, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer fileReader.Close()

	// Upload to Gemini
	f, err := m.genaiClient.UploadFile(ctx, "", fileReader, &genai.UploadFileOptions{
		MIMEType:    mimeType,
		DisplayName: fileName,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upload file %s: %w", fileName, err)
	}

	// Wait for processing
	timeout := time.After(30 * time.Second)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			return nil, fmt.Errorf("file processing timeout for %s", fileName)
		case <-ticker.C:
			f, err = m.genaiClient.GetFile(ctx, f.Name)
			if err != nil {
				return nil, fmt.Errorf("failed to get file status for %s: %w", fileName, err)
			}
			if f.State == genai.FileStateActive {
				processedFile := &ProcessedFile{
					Name:     fileName,
					Path:     filePath,
					MIMEType: f.MIMEType,
					URI:      f.URI,
				}
				session.AddFile(processedFile)
				return processedFile, nil
			}
			if f.State == genai.FileStateFailed {
				return nil, fmt.Errorf("file upload failed for %s", fileName)
			}
		}
	}
}

// cleanupSessionFiles deletes all files associated with a session from the Gemini API.
func (m *Manager) cleanupSessionFiles(ctx context.Context, session *Session) {
	session.mu.RLock()
	defer session.mu.RUnlock()

	for _, file := range session.ProcessedFiles {
		// The file name is the last part of the URI
		parts := strings.Split(file.URI, "/")
		fileName := parts[len(parts)-1]
		err := m.genaiClient.DeleteFile(ctx, fileName)
		if err != nil {
			log.Printf("Failed to delete file %s for session %s: %v", file.Name, session.ID, err)
		}
	}
}
