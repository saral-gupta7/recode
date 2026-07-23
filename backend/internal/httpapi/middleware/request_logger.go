package middleware

import (
	"crypto/rand"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/saral-gupta7/recode/backend/internal/observability"
)

const requestIDHeader = "X-Request-ID"

func newRequestID() string {
	return rand.Text()
}

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := newRequestID()

		c.Header(requestIDHeader, requestID)

		requestContext := observability.WithRequestID(c.Request.Context(), requestID)

		c.Request = c.Request.WithContext(requestContext)

		c.Next()
	}
}

func RequestLogger(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		startedAt := time.Now()

		c.Next()

		duration := time.Since(startedAt)

		requestID, ok := observability.RequestIDFromContext(c.Request.Context())

		if !ok {
			requestID = "unknown"
		}

		route := c.FullPath()

		if route == "" {
			route = "unmatched"
		}

		logger.InfoContext(c.Request.Context(),
			"http request completed",
			slog.String("request_id", requestID),
			slog.String("method", c.Request.Method),
			slog.String("route", route),
			slog.Int("status", c.Writer.Status()),
			slog.Int("response_bytes", c.Writer.Size()),
			slog.Duration("duration", duration))
	}
}
