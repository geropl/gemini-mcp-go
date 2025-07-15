# Gemini MCP Go

A powerful, stateful MCP server for Google's Gemini API, implemented in Go. It allows an MCP client to have complex, multi-turn conversations with Gemini, including attaching files and maintaining context across multiple requests.

This project is a Go reimplementation and enhancement of the concepts found in the `mcp-gemini-assistant` Python project.

## Key Features

-   **Session Management**: Maintain conversation context across multiple queries using session IDs.
-   **File Attachments**: Read and include local code files in conversations. The server handles uploading files to the Gemini API.
-   **Context Caching**: Code context and file content are cached per session, reducing token usage on follow-up questions.
-   **Automatic Cleanup**: Sessions automatically expire after 1 hour of inactivity, and all associated uploaded files are deleted from the Gemini API to save space.
-   **Parallel Conversations**: The server can handle multiple, independent sessions at once.

## Installation

1.  Ensure you have Go (1.21 or later) installed.
2.  Clone this repository:
    ```bash
    git clone https://github.com/geropl/gemini-mcp-go.git
    cd gemini-mcp-go
    ```
3.  Build the server:
    ```bash
    go build -o gemini-mcp-go .
    ```

## Usage

### 1. Set API Key

The server requires a Gemini API key. You can set it using the `setup` command, which will guide you, or by setting the environment variable directly.

```bash
# Guided setup
./gemini-mcp-go setup

# Or set it directly
export GEMINI_API_KEY="your-gemini-api-key-here"
```

### 2. Start the Server

Run the `serve` command to start the MCP server.

```bash
./gemini-mcp-go serve
```

The server will start and listen for requests from an MCP client on stdin/stdout.

## Tools Available

### 1. `consult_gemini`

Start or continue a conversation with Gemini. This is the primary tool for interacting with the server.

**Parameters:**

-   `session_id` (string, optional): The ID of a previous session to continue the conversation. If omitted, a new session is created.
-   `problem_description` (string, required for new sessions): A detailed description of the coding problem.
-   `code_context` (string, optional): A block of code relevant to the problem. This is cached for the session.
-   `attached_files` (array of strings, optional): A list of absolute file paths to read from the local filesystem and attach to the conversation.
-   `file_descriptions` (object, optional): A map where keys are file paths (matching those in `attached_files`) and values are descriptions of the files.
-   `specific_question` (string, required): The specific question you want to ask Gemini.
-   `additional_context` (string, optional): Any new information, updates, or changes since the last question in the session.
-   `preferred_approach` (string, optional): The type of help needed (e.g., "solution", "review", "debug", "optimize", "explain", "follow-up").

### 2. `list_sessions`

List all active consultation sessions currently managed by the server.

### 3. `end_session`

End a specific session to free up memory and delete any associated files uploaded to the Gemini API.

**Parameters:**

-   `session_id` (string, required): The ID of the session to terminate.

## Example Workflow

**1. Start a new conversation with file attachments:**

An MCP client would send a `call` request for the `consult_gemini` tool with parameters like:

```json
{
  "problem_description": "I need to optimize this React component for performance",
  "attached_files": [
    "/path/to/src/components/Dashboard.jsx",
    "/path/to/src/hooks/useData.js"
  ],
  "file_descriptions": {
    "/path/to/src/components/Dashboard.jsx": "Main dashboard component with performance issues",
    "/path/to/src/hooks/useData.js": "Custom hook for data fetching"
  },
  "specific_question": "How can I improve the rendering performance of this dashboard?",
  "preferred_approach": "optimize"
}
```

The server's response will include a `session_id` for continuing the conversation.

**2. Ask a follow-up question:**

Using the `session_id` from the previous response:

```json
{
  "session_id": "abc-123-def-456",
  "specific_question": "I implemented your suggestion, but now I'm getting stale data issues. How do I handle cache invalidation?",
  "additional_context": "Added the LRU cache as suggested, but users see old data after updates."
}
```

**3. List active sessions:**

To see all current conversations:

```json
{}
```

This returns information about all active sessions including their IDs, creation times, and message counts.

**4. End the conversation:**

Once the problem is solved, end the session to clean up resources:

```json
{
  "session_id": "abc-123-def-456"
}
```

This will delete any uploaded files from the Gemini API and free up memory.

## Session Management

### Automatic Cleanup
- Sessions automatically expire after **1 hour** of inactivity
- Expired sessions and their uploaded files are automatically deleted
- A background cleanup routine runs every 5 minutes to remove expired sessions

### File Handling
- Files are uploaded to the Gemini API when first referenced in a session
- File content is cached per session to reduce token usage on follow-up questions
- Files are automatically deleted from the Gemini API when sessions end or expire
- Supported file types include most common programming languages and text formats

### Concurrent Sessions
- The server can handle multiple independent sessions simultaneously
- Each session maintains its own conversation context and file cache
- Thread-safe operations ensure data consistency across concurrent requests

## Testing

To run the project's tests, use the following command:

```bash
go test ./...
