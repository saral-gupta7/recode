package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/saral-gupta7/recode/backend/internal/observability"
)

func TestAbortWithErrorWritesStandardEnvelope(t *testing.T) {
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)

	request := httptest.NewRequest(http.MethodGet, "/jobs/missing", nil)
	requestContext := observability.WithRequestID(request.Context(), "request-123")
	c.Request = request.WithContext(requestContext)

	AbortWithError(
		c,
		http.StatusNotFound,
		"job_not_found",
		"The requested job was not found.",
	)

	if !c.IsAborted() {
		t.Fatal("Gin context is not aborted")
	}

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusNotFound)
	}

	wantContentType := "application/json; charset=utf-8"
	if got := recorder.Header().Get("Content-Type"); got != wantContentType {
		t.Fatalf("Content-Type = %q, want %q", got, wantContentType)
	}

	var got errorEnvelope
	if err := json.NewDecoder(recorder.Body).Decode(&got); err != nil {
		t.Fatalf("decode error response: %v", err)
	}

	want := errorBody{
		Code:      "job_not_found",
		Message:   "The requested job was not found.",
		RequestID: "request-123",
	}

	if got.Error != want {
		t.Fatalf("error body = %#v, want %#v", got.Error, want)
	}
}
