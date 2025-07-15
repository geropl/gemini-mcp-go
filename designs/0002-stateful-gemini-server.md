# Design Doc: Stateful Gemini Server with Session Management

**Objective:** Transform `gemini-mcp-go` from a stateless server into a stateful, session-based server that mirrors the functionality of the `mcp-gemini-assistant` Python project. This includes adding support for session management, file attachments, and context caching.

---

## Phase 1: Core Data Structures and Session Management

The foundation of the new functionality is stateful session management. We'll create the necessary data structures and the session manager.

### Sub-Tasks
- [x] **1.1: Create `pkg/session/session.go`**
    - [x] Define `ProcessedFile` struct.
    - [x] Define `Session` struct with all required fields (`ID`, `Chat`, `Created`, `LastUsed`, `MessageCount`, `ProblemDescription`, `CodeContext`, `ProcessedFiles`, `sync.RWMutex`).
- [x] **1.2: Create `pkg/session/manager.go`**
    - [x] Define `Manager` struct (`map[string]*Session`, `sync.RWMutex`, `*genai.Client`).
    - [x] Implement `NewManager()` constructor.
    - [x] Implement `GetOrCreateSession(sessionID string) (*Session, error)`.
    - [x] Implement `ListSessions() []*Session`.
    - [x] Implement `EndSession(sessionID string) error`.
- [x] **1.3: Implement Session Cleanup**
    - [x] Add `StartCleanupRoutine()` method to the `Manager`.
    - [x] Implement the goroutine to periodically check for and remove expired sessions.

---

## Phase 2: File Processing and Gemini Integration

This phase focuses on handling file uploads, which is a key feature of the `consult_gemini` tool.

### Sub-Tasks
- [x] **2.1: Implement File Processing in Session Manager**
    - [x] Create `ProcessFile(session *Session, filePath string) (*ProcessedFile, error)` method in `pkg/session/manager.go`.
    - [x] Implement file existence check.
    - [x] Implement MIME type guessing helper.
    - [x] Implement file upload logic using `generative-ai-go`.
    - [x] Implement polling for `ACTIVE` or `FAILED` state.
    - [x] Implement logic to store `ProcessedFile` info in the session.
- [x] **2.2: Implement File Cleanup**
    - [x] Create `cleanupSessionFiles(session *Session)` method in `pkg/session/manager.go`.
    - [x] Implement logic to call Gemini client's `DeleteFile` method.
    - [x] Integrate `cleanupSessionFiles` into `EndSession` and the session cleanup routine.

---

## Phase 3: Implementing the New MCP Tools

Now we'll create the user-facing tools, deprecating the old ones.

### Sub-Tasks
- [x] **3.1: Create `pkg/handlers/consult_gemini.go`**
    - [x] Define `ConsultGeminiParams` struct.
    - [x] Create `HandleConsultGemini` function.
    - [x] Implement logic to get or create a session.
    - [x] Implement logic for handling a new session (prompt building, file processing).
    - [x] Implement logic for sending the question to Gemini.
    - [x] Implement response formatting.
- [x] **3.2: Create `pkg/handlers/session_tools.go`**
    - [x] Implement `HandleListSessions`.
    - [x] Implement `HandleEndSession`.
- [x] **3.3: Update `cmd/serve.go`**
    - [x] Instantiate the `session.NewManager()` in `serveCmd`.
    - [x] Launch the session cleanup routine.
    - [x] Register handlers for `consult_gemini`, `list_sessions`, and `end_session`.
    - [x] Deprecate (remove or comment out) old `generate` and `stream` handlers.
    - **Note:** The server setup logic was implemented in `pkg/server/server.go`, which is called by `cmd/serve.go`.

---

## Phase 4: Testing and Documentation

With the new features in place, we need to ensure they are reliable and well-documented.

### Sub-Tasks
- [x] **4.1: Documentation Updates**
    - [x] Update all relevant documentation files to reflect new stateful architecture.
    - [x] Delete outdated documentation files.
    - [x] Ensure memory bank files are current and accurate.
- [x] **4.2: Unit Tests for Session Manager**
    - [x] Write comprehensive tests for `GetOrCreateSession`, `EndSession`, and `ListSessions`.
    - [x] Create mock Gemini API client for testing.
    - [x] Write tests for `ProcessFile` and `cleanupSessionFiles` using the mock.
    - [x] Test session cleanup routine and TTL behavior.
- [x] **4.3: Clean Up Existing Tests**
    - [x] Delete all existing integration tests for tools to start fresh.
    - [x] Remove test infrastructure hints and outdated test approaches.
    - [x] Prepare clean foundation for new integration tests in future phases.

**Note:** Integration tests for the new handlers (`consult_gemini`, `list_sessions`, `end_session`) will be implemented in a separate phase after the foundation is solid.
