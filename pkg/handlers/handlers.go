package handlers

// This package contains MCP tool handlers for the gemini-mcp-go server.
// 
// Current tools:
// - consult_gemini: Primary tool for stateful conversations with file attachment support
// - list_sessions: List all active consultation sessions
// - end_session: End a specific session and clean up resources
//
// Legacy handlers (generate, stream) have been removed in favor of the new session-based approach.
