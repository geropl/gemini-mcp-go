# Product Context: gemini-mcp-go

## Problem Statement

AI assistants like Claude Desktop need sophisticated ways to interact with Google's Gemini API for complex coding problems. Existing solutions have limitations:

1.  **Stateless Interactions:** Each query starts fresh without conversation context, making complex problem-solving inefficient.
2.  **Limited File Support:** Difficulty attaching and referencing multiple code files in conversations.
3.  **Resource Management:** No automatic cleanup of uploaded files, leading to API quota waste.
4.  **Context Loss:** Unable to maintain conversation history for follow-up questions and iterative problem-solving.

## Solution

The `gemini-mcp-go` project provides a stateful, session-based MCP server that enables:

1.  **Persistent Conversations:** Maintain chat history and context across multiple queries within sessions.
2.  **File Attachment Support:** Read local code files and attach them to conversations with automatic upload to Gemini API.
3.  **Smart Resource Management:** Automatic cleanup of sessions and uploaded files with TTL-based expiration.
4.  **Context Caching:** Code context and file content are cached per session, reducing token usage on follow-up questions.

## Target Users

-   **AI-Assisted Developers:** Developers using Claude Desktop or other MCP clients who need complex, multi-turn conversations about code.
-   **Code Reviewers:** Teams that want to use AI for thorough code analysis with full file context.
-   **Technical Consultants:** Professionals who need to analyze and provide guidance on complex codebases.
-   **Learning Developers:** Students and junior developers who benefit from iterative, contextual coding assistance.

## How It Works

`gemini-mcp-go` functions as a stateful MCP server that:

1.  **Manages Sessions:** Creates and maintains conversation sessions with unique IDs and persistent chat history.
2.  **Processes Files:** Reads local code files, uploads them to Gemini API, and tracks them per session with caching.
3.  **Handles Conversations:** Facilitates multi-turn conversations with maintained context and file references.
4.  **Cleans Up Resources:** Automatically expires sessions after 1 hour and deletes uploaded files from Gemini API.
5.  **Provides Tools:** Offers `consult_gemini`, `list_sessions`, and `end_session` tools for comprehensive session management.
6.  **Ensures Thread Safety:** Handles multiple concurrent sessions with proper locking and resource management.

## User Experience Goals

-   **Contextual Conversations:** Users can have complex, multi-turn discussions about code with maintained context.
-   **Seamless File Integration:** Easy attachment of local code files with automatic upload and referencing.
-   **Efficient Resource Usage:** Automatic cleanup prevents API quota waste and manages memory efficiently.
-   **Simple Session Management:** Clear tools for creating, listing, and ending conversation sessions.
-   **Reliable Performance:** Robust session handling with proper concurrency and error management.
