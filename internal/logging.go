package middleware

import (
	"context"
	"log/slog"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// WithLogging wraps a tool handler with logging
func WithLogging(logger *slog.Logger, toolName string, handler server.ToolHandlerFunc) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		start := time.Now()

		logger.Info("tool call started",
			"tool", toolName,
			"arguments", request.GetArguments(),
		)

		result, err := handler(ctx, request)

		duration := time.Since(start)

		if err != nil {
			logger.Error("tool call failed",
				"tool", toolName,
				"duration_ms", duration.Milliseconds(),
				"error", err.Error(),
			)
		} else {
			logger.Info("tool call completed",
				"tool", toolName,
				"duration_ms", duration.Milliseconds(),
			)
		}

		return result, err
	}
}

// WrapServerWithLogging wraps all tools in a server with logging middleware
func WrapServerWithLogging(s *server.MCPServer, logger *slog.Logger) {
	tools := s.ListTools()

	for toolName, serverTool := range tools {
		// Capture variables for closure
		name := toolName
		originalHandler := serverTool.Handler

		// Wrap the handler with logging
		wrappedHandler := WithLogging(logger, name, originalHandler)

		// Re-register the tool with the wrapped handler
		s.AddTool(serverTool.Tool, wrappedHandler)
	}
}
