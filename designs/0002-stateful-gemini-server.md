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
- [ ] **4.1: Unit Tests for Session Manager**
    - [ ] Write tests for `GetOrCreateSession`, `EndSession`, and `ListSessions`.
    - [ ] Mock the Gemini API client.
    - [ ] Write tests for `ProcessFile` and `cleanupSessionFiles` using the mock.
- [ ] **4.2: Integration Tests for Handlers**
    - [ ] Create new golden file tests for `consult_gemini`.
    - [ ] Create golden file tests for `list_sessions`.
    - [ ] Create golden file tests for `end_session`.
- [ ] **4.3: Update Documentation**
    - [ ] Update `README.md` with documentation for the new tools.
    - [ ] Update `memory-bank/projectbrief.md` to reflect the new stateful architecture.
    - [ ] Update `memory-bank/systemPatterns.md` with the new session management components.
