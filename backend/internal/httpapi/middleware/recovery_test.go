package middleware

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRecoveryReturnsSafeErrorAndLogsPanic(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var logs bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logs, nil))

	router := gin.New()
	router.Use(RequestID())
	router.Use(RequestLogger(logger))
	router.Use(Recovery(logger))
	router.GET("/panic/:panicID", func(_ *gin.Context) {
		panic("private panic detail")
	})

	request := httptest.NewRequest(http.MethodGet, "/panic/123", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusInternalServerError)
	}

	requestID := response.Header().Get(requestIDHeader)
	if requestID == "" {
		t.Fatal("X-Request-ID header is empty")
	}

	var body struct {
		Error struct {
			Code      string `json:"code"`
			Message   string `json:"message"`
			RequestID string `json:"request_id"`
		} `json:"error"`
	}

	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatalf("decode recovery response: %v", err)
	}

	if body.Error.Code != "internal_error" {
		t.Fatalf("error code = %q, want %q", body.Error.Code, "internal_error")
	}

	if body.Error.Message != "An unexpected error occurred." {
		t.Fatalf(
			"error message = %q, want a safe generic message",
			body.Error.Message,
		)
	}

	if body.Error.RequestID != requestID {
		t.Fatalf(
			"response request ID = %q, header request ID = %q",
			body.Error.RequestID,
			requestID,
		)
	}

	if strings.Contains(response.Body.String(), "private panic detail") {
		t.Fatal("response exposes the private panic value")
	}

	type logRecord struct {
		Level     string `json:"level"`
		Message   string `json:"msg"`
		RequestID string `json:"request_id"`
		Panic     string `json:"panic"`
		Method    string `json:"method"`
		Path      string `json:"path"`
		Route     string `json:"route"`
		Status    int    `json:"status"`
		Stack     string `json:"stack"`
	}

	var records []logRecord
	decoder := json.NewDecoder(&logs)
	for {
		var record logRecord
		err := decoder.Decode(&record)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			t.Fatalf("decode structured log: %v", err)
		}

		records = append(records, record)
	}

	var panicRecord *logRecord
	var completionRecord *logRecord

	for i := range records {
		switch records[i].Message {
		case "panic recovered":
			panicRecord = &records[i]
		case "http request completed":
			completionRecord = &records[i]
		}
	}

	if panicRecord == nil {
		t.Fatal("panic recovery log is missing")
	}
	if panicRecord.Level != "ERROR" {
		t.Fatalf("panic log level = %q, want %q", panicRecord.Level, "ERROR")
	}
	if panicRecord.RequestID != requestID {
		t.Fatalf("panic log request ID = %q, want %q", panicRecord.RequestID, requestID)
	}
	if panicRecord.Panic != "private panic detail" {
		t.Fatalf("panic log value = %q, want the internal panic detail", panicRecord.Panic)
	}
	if panicRecord.Method != http.MethodGet || panicRecord.Path != "/panic/123" {
		t.Fatalf(
			"panic log request = %s %s, want GET /panic/123",
			panicRecord.Method,
			panicRecord.Path,
		)
	}
	if panicRecord.Stack == "" {
		t.Fatal("panic stack trace is empty")
	}

	if completionRecord == nil {
		t.Fatal("request completion log is missing")
	}
	if completionRecord.RequestID != requestID {
		t.Fatalf(
			"completion log request ID = %q, want %q",
			completionRecord.RequestID,
			requestID,
		)
	}
	if completionRecord.Route != "/panic/:panicID" {
		t.Fatalf(
			"completion log route = %q, want %q",
			completionRecord.Route,
			"/panic/:panicID",
		)
	}
	if completionRecord.Status != http.StatusInternalServerError {
		t.Fatalf(
			"completion log status = %d, want %d",
			completionRecord.Status,
			http.StatusInternalServerError,
		)
	}
}
