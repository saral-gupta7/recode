package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"

	"github.com/saral-gupta7/recode/backend/internal/httpapi/response"
	"github.com/saral-gupta7/recode/backend/internal/observability"
)

func Recovery(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			recovered := recover()
			if recovered == nil {
				return
			}

			requestID, ok := observability.RequestIDFromContext(
				c.Request.Context(),
			)
			if !ok {
				requestID = "unknown"
			}

			logger.ErrorContext(
				c.Request.Context(),
				"panic recovered",
				slog.String("request_id", requestID),
				slog.Any("panic", recovered),
				slog.String("method", c.Request.Method),
				slog.String("path", c.Request.URL.Path),
				slog.String("stack", string(debug.Stack())),
			)

			if c.Writer.Written() {
				c.Abort()
				return
			}

			response.AbortWithError(
				c,
				http.StatusInternalServerError,
				"internal_error",
				"An unexpected error occurred.",
			)
		}()

		c.Next()
	}
}
