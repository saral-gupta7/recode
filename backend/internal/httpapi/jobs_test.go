package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"testing"
	"time"

	applicationjobs "github.com/saral-gupta7/recode/backend/internal/application/jobs"
	"github.com/saral-gupta7/recode/backend/internal/job"
	"github.com/saral-gupta7/recode/backend/internal/queue"
	"github.com/saral-gupta7/recode/backend/internal/storage"
	"github.com/saral-gupta7/recode/backend/internal/task"
)

type httpRepositoryStub struct {
	records map[string]task.Record
}

func (r *httpRepositoryStub) Create(_ context.Context, record task.Record) error {
	r.records[record.Job.ID()] = record
	return nil
}

func (r *httpRepositoryStub) FindByID(_ context.Context, id string) (task.Record, error) {
	record, ok := r.records[id]
	if !ok {
		return task.Record{}, task.ErrNotFound
	}
	return record, nil
}

func (r *httpRepositoryStub) FindExpirable(context.Context, time.Time, int) ([]task.Record, error) {
	return nil, nil
}

func (r *httpRepositoryStub) SaveJob(_ context.Context, storedJob *job.Job) error {
	record, ok := r.records[storedJob.ID()]
	if !ok {
		return task.ErrNotFound
	}
	record.Job = storedJob
	r.records[storedJob.ID()] = record
	return nil
}

func (r *httpRepositoryStub) Delete(_ context.Context, id string) error {
	delete(r.records, id)
	return nil
}

type httpQueueStub struct {
	ids []string
}

func (q *httpQueueStub) Enqueue(_ context.Context, id string) error {
	q.ids = append(q.ids, id)
	return nil
}

func (q *httpQueueStub) Dequeue(context.Context, time.Duration) (string, error) {
	return "", queue.ErrEmpty
}

func TestCreateAndGetJob(t *testing.T) {
	repository := &httpRepositoryStub{records: make(map[string]task.Record)}
	jobQueue := &httpQueueStub{}
	store, err := storage.NewLocal(t.TempDir())
	if err != nil {
		t.Fatalf("NewLocal() error = %v", err)
	}
	service := applicationjobs.New(repository, jobQueue, store, 1024, time.Hour)
	router := NewRouter(testLogger(), Dependencies{
		Jobs:           service,
		MaxUploadBytes: 1024,
	})

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)
	if err := writer.WriteField("operation", string(job.OperationImageResize)); err != nil {
		t.Fatalf("write operation: %v", err)
	}
	if err := writer.WriteField("options", `{"width":320}`); err != nil {
		t.Fatalf("write options: %v", err)
	}
	header := make(textproto.MIMEHeader)
	header.Set("Content-Disposition", `form-data; name="file"; filename="photo.png"`)
	header.Set("Content-Type", "image/png")
	part, err := writer.CreatePart(header)
	if err != nil {
		t.Fatalf("create file part: %v", err)
	}
	if _, err := io.WriteString(part, "image bytes"); err != nil {
		t.Fatalf("write file: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("close multipart writer: %v", err)
	}

	request := httptest.NewRequest(http.MethodPost, "/api/v1/jobs", &requestBody)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusAccepted {
		t.Fatalf("POST status = %d, want %d; body=%s", response.Code, http.StatusAccepted, response.Body)
	}
	var created jobResponse
	if err := json.NewDecoder(response.Body).Decode(&created); err != nil {
		t.Fatalf("decode create response: %v", err)
	}
	if created.ID == "" || created.OwnerToken == "" {
		t.Fatalf("create response missing credentials: %+v", created)
	}
	if len(jobQueue.ids) != 1 || jobQueue.ids[0] != created.ID {
		t.Fatalf("queued IDs = %v, want [%s]", jobQueue.ids, created.ID)
	}

	getRequest := httptest.NewRequest(http.MethodGet, "/api/v1/jobs/"+created.ID, nil)
	getRequest.Header.Set("Authorization", "Bearer "+created.OwnerToken)
	getResponse := httptest.NewRecorder()
	router.ServeHTTP(getResponse, getRequest)
	if getResponse.Code != http.StatusOK {
		t.Fatalf("GET status = %d, want %d; body=%s", getResponse.Code, http.StatusOK, getResponse.Body)
	}
}

func TestGetJobRejectsMissingOwnershipToken(t *testing.T) {
	repository := &httpRepositoryStub{records: make(map[string]task.Record)}
	store, err := storage.NewLocal(t.TempDir())
	if err != nil {
		t.Fatalf("NewLocal() error = %v", err)
	}
	service := applicationjobs.New(repository, &httpQueueStub{}, store, 1024, time.Hour)
	created, err := service.Create(context.Background(), applicationjobs.CreateInput{
		Operation: job.OperationImageGrayscale,
		Content:   bytes.NewBufferString("image"),
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	router := NewRouter(testLogger(), Dependencies{Jobs: service, MaxUploadBytes: 1024})

	request := httptest.NewRequest(http.MethodGet, "/api/v1/jobs/"+created.Record.Job.ID(), nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("GET status = %d, want %d", response.Code, http.StatusUnauthorized)
	}
}
