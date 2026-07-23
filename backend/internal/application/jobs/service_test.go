package jobs

import (
	"bytes"
	"context"
	"errors"
	"io"
	"sync"
	"testing"
	"time"

	"github.com/saral-gupta7/recode/backend/internal/job"
	"github.com/saral-gupta7/recode/backend/internal/queue"
	"github.com/saral-gupta7/recode/backend/internal/storage"
	"github.com/saral-gupta7/recode/backend/internal/task"
)

func TestServiceCreateAndAuthorize(t *testing.T) {
	repository := newFakeRepository()
	jobQueue := &fakeQueue{}
	store := newFakeStore()
	service := New(repository, jobQueue, store, 100, time.Hour)

	created, err := service.Create(context.Background(), CreateInput{
		Operation:        job.OperationImageGrayscale,
		OriginalFilename: "../photo.png",
		MIMEType:         "image/png",
		Content:          bytes.NewBufferString("image"),
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if created.OwnerToken == "" {
		t.Fatal("OwnerToken is empty")
	}
	if created.Record.OriginalFilename != "photo.png" {
		t.Errorf("OriginalFilename = %q, want photo.png", created.Record.OriginalFilename)
	}
	if len(jobQueue.ids) != 1 || jobQueue.ids[0] != created.Record.Job.ID() {
		t.Errorf("queued IDs = %v", jobQueue.ids)
	}

	if _, err := service.Get(context.Background(), created.Record.Job.ID(), "wrong"); !errors.Is(err, ErrUnauthorized) {
		t.Fatalf("Get() wrong-token error = %v, want ErrUnauthorized", err)
	}
	if _, err := service.Get(context.Background(), created.Record.Job.ID(), created.OwnerToken); err != nil {
		t.Fatalf("Get() correct-token error = %v", err)
	}
}

func TestServiceCancelQueuedJob(t *testing.T) {
	repository := newFakeRepository()
	service := New(repository, &fakeQueue{}, newFakeStore(), 100, time.Hour)

	created, err := service.Create(context.Background(), CreateInput{
		Operation: job.OperationImageGrayscale,
		Content:   bytes.NewBufferString("image"),
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	record, err := service.Cancel(
		context.Background(),
		created.Record.Job.ID(),
		created.OwnerToken,
	)
	if err != nil {
		t.Fatalf("Cancel() error = %v", err)
	}
	if record.Job.Status() != job.StatusCancelled {
		t.Errorf("Status() = %q, want cancelled", record.Job.Status())
	}
	if _, scheduled := record.Job.ExpiresAt(); !scheduled {
		t.Error("cancelled job has no expiration")
	}
}

type fakeRepository struct {
	mu      sync.Mutex
	records map[string]task.Record
}

func newFakeRepository() *fakeRepository {
	return &fakeRepository{records: make(map[string]task.Record)}
}

func (r *fakeRepository) Create(_ context.Context, record task.Record) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.records[record.Job.ID()] = record
	return nil
}

func (r *fakeRepository) FindByID(_ context.Context, id string) (task.Record, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	record, ok := r.records[id]
	if !ok {
		return task.Record{}, task.ErrNotFound
	}
	return record, nil
}

func (r *fakeRepository) FindExpirable(_ context.Context, _ time.Time, _ int) ([]task.Record, error) {
	return nil, nil
}

func (r *fakeRepository) SaveJob(_ context.Context, storedJob *job.Job) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	record, ok := r.records[storedJob.ID()]
	if !ok {
		return task.ErrNotFound
	}
	record.Job = storedJob
	r.records[storedJob.ID()] = record
	return nil
}

func (r *fakeRepository) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.records, id)
	return nil
}

type fakeQueue struct {
	ids []string
}

func (q *fakeQueue) Enqueue(_ context.Context, id string) error {
	q.ids = append(q.ids, id)
	return nil
}

func (q *fakeQueue) Dequeue(_ context.Context, _ time.Duration) (string, error) {
	if len(q.ids) == 0 {
		return "", queue.ErrEmpty
	}
	id := q.ids[0]
	q.ids = q.ids[1:]
	return id, nil
}

type fakeStore struct {
	objects map[string][]byte
}

func newFakeStore() *fakeStore {
	return &fakeStore{objects: make(map[string][]byte)}
}

func (s *fakeStore) Put(_ context.Context, key string, source io.Reader, max int64) (int64, error) {
	data, err := io.ReadAll(io.LimitReader(source, max+1))
	if err != nil {
		return 0, err
	}
	if int64(len(data)) > max {
		return 0, storage.ErrTooLarge
	}
	s.objects[key] = data
	return int64(len(data)), nil
}

func (s *fakeStore) Open(_ context.Context, key string) (io.ReadCloser, error) {
	data, ok := s.objects[key]
	if !ok {
		return nil, storage.ErrNotFound
	}
	return io.NopCloser(bytes.NewReader(data)), nil
}

func (s *fakeStore) Delete(_ context.Context, key string) error {
	delete(s.objects, key)
	return nil
}
