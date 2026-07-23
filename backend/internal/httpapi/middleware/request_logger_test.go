package middleware

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/saral-gupta7/recode/backend/internal/observability"
)

func TestRequestIDAddsUniqueResponseHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(RequestID())
	router.GET("/", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	firstID := performRequest(t, router)
	secondID := performRequest(t, router)

	if firstID == secondID {
		t.Fatalf("request IDs are equal: %q", firstID)
	}
}

func TestRequestIDPropagatesToRequestContext(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var contextRequestID string

	router := gin.New()
	router.Use(RequestID())
	router.GET("/", func(c *gin.Context) {
		requestID, ok := observability.RequestIDFromContext(c.Request.Context())
		if !ok {
			c.Status(http.StatusInternalServerError)
			return
		}

		contextRequestID = requestID
		c.Status(http.StatusNoContent)
	})

	responseRequestID := performRequest(t, router)

	if contextRequestID != responseRequestID {
		t.Fatalf(
			"context request ID = %q, response request ID = %q",
			contextRequestID,
			responseRequestID,
		)
	}
}

func TestRequestLoggerWritesStructuredCompletionRecord(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var logs bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logs, nil))

	router := gin.New()
	router.Use(RequestID())
	router.Use(RequestLogger(logger))
	router.GET("/items/:itemID", func(c *gin.Context) {
		c.String(http.StatusCreated, "created")
	})

	request := httptest.NewRequest(http.MethodGet, "/items/123", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	var record struct {
		Message       string `json:"msg"`
		RequestID     string `json:"request_id"`
		Method        string `json:"method"`
		Route         string `json:"route"`
		Status        int    `json:"status"`
		ResponseBytes int    `json:"response_bytes"`
		Duration      *int64 `json:"duration"`
	}

	if err := json.NewDecoder(&logs).Decode(&record); err != nil {
		t.Fatalf("decode structured request log: %v", err)
	}

	if record.Message != "http request completed" {
		t.Fatalf("log message = %q, want %q", record.Message, "http request completed")
	}

	responseRequestID := response.Header().Get(requestIDHeader)
	if record.RequestID == "" || record.RequestID != responseRequestID {
		t.Fatalf(
			"logged request ID = %q, response request ID = %q",
			record.RequestID,
			responseRequestID,
		)
	}

	if record.Method != http.MethodGet {
		t.Fatalf("logged method = %q, want %q", record.Method, http.MethodGet)
	}

	if record.Route != "/items/:itemID" {
		t.Fatalf("logged route = %q, want %q", record.Route, "/items/:itemID")
	}

	if record.Status != http.StatusCreated {
		t.Fatalf("logged status = %d, want %d", record.Status, http.StatusCreated)
	}

	if record.ResponseBytes != len("created") {
		t.Fatalf(
			"logged response bytes = %d, want %d",
			record.ResponseBytes,
			len("created"),
		)
	}

	if record.Duration == nil || *record.Duration < 0 {
		t.Fatalf("logged duration = %v, want a non-negative duration", record.Duration)
	}
}

func performRequest(t *testing.T, router http.Handler) string {
	t.Helper()

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	requestID := response.Header().Get(requestIDHeader)
	if requestID == "" {
		t.Fatal("X-Request-ID header is empty")
	}

	return requestID
}
