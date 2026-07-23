package worker

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/saral-gupta7/recode/backend/internal/job"
	"github.com/saral-gupta7/recode/backend/internal/media"
	"github.com/saral-gupta7/recode/backend/internal/queue"
	"github.com/saral-gupta7/recode/backend/internal/storage"
	"github.com/saral-gupta7/recode/backend/internal/task"
)

type repositoryStub struct {
	record task.Record
}

func (r *repositoryStub) Create(context.Context, task.Record) error {
	return nil
}

func (r *repositoryStub) FindByID(_ context.Context, id string) (task.Record, error) {
	if r.record.Job == nil || r.record.Job.ID() != id {
		return task.Record{}, task.ErrNotFound
	}
	return r.record, nil
}

func (r *repositoryStub) FindExpirable(context.Context, time.Time, int) ([]task.Record, error) {
	return nil, nil
}

func (r *repositoryStub) SaveJob(_ context.Context, storedJob *job.Job) error {
	r.record.Job = storedJob
	return nil
}

func (r *repositoryStub) Delete(context.Context, string) error {
	return nil
}

type queueStub struct{}

func (queueStub) Enqueue(context.Context, string) error {
	return nil
}

func (queueStub) Dequeue(context.Context, time.Duration) (string, error) {
	return "", queue.ErrEmpty
}

type proberStub struct {
	info media.MediaInfo
	err  error
}

func (p proberStub) Probe(context.Context, string) (media.MediaInfo, error) {
	return p.info, p.err
}

type processorStub struct {
	called bool
}

func (p *processorStub) Process(
	_ context.Context,
	_ job.Operation,
	_ task.Options,
	_ string,
	outputPath string,
) error {
	p.called = true
	return os.WriteFile(outputPath, []byte("processed"), 0o600)
}

func TestProcessCompletesQueuedJob(t *testing.T) {
	store, err := storage.NewLocal(t.TempDir())
	if err != nil {
		t.Fatalf("NewLocal() error = %v", err)
	}
	if _, err := store.Put(
		context.Background(),
		"uploads/job-1/input",
		bytes.NewBufferString("source"),
		1024,
	); err != nil {
		t.Fatalf("store input: %v", err)
	}

	storedJob, err := job.New("job-1", job.OperationImageGrayscale, time.Now().UTC())
	if err != nil {
		t.Fatalf("job.New() error = %v", err)
	}
	repository := &repositoryStub{record: task.Record{
		Job:      storedJob,
		InputKey: "uploads/job-1/input",
	}}
	processor := &processorStub{}
	runner := New(
		slog.New(slog.NewTextHandler(io.Discard, nil)),
		repository,
		queueStub{},
		store,
		proberStub{info: media.MediaInfo{HasVideo: true}},
		processor,
		time.Hour,
		1024,
	)

	if err := runner.process(context.Background(), "job-1"); err != nil {
		t.Fatalf("process() error = %v", err)
	}
	if !processor.called {
		t.Fatal("processor was not called")
	}
	if got := repository.record.Job.Status(); got != job.StatusCompleted {
		t.Fatalf("status = %q, want %q", got, job.StatusCompleted)
	}
	resultKey, ok := repository.record.Job.ResultKey()
	if !ok {
		t.Fatal("completed job has no result key")
	}
	result, err := store.Open(context.Background(), resultKey)
	if err != nil {
		t.Fatalf("open result: %v", err)
	}
	defer result.Close()
	body, err := io.ReadAll(result)
	if err != nil {
		t.Fatalf("read result: %v", err)
	}
	if string(body) != "processed" {
		t.Fatalf("result = %q, want processed", body)
	}
	if _, err := store.Open(context.Background(), "uploads/job-1/input"); !errors.Is(err, storage.ErrNotFound) {
		t.Fatalf("input open error = %v, want ErrNotFound", err)
	}
}

func TestProcessRejectsInvalidMediaBeforeProcessor(t *testing.T) {
	root := t.TempDir()
	store, err := storage.NewLocal(root)
	if err != nil {
		t.Fatalf("NewLocal() error = %v", err)
	}
	if _, err := store.Put(
		context.Background(),
		"uploads/job-2/input",
		bytes.NewBufferString("not media"),
		1024,
	); err != nil {
		t.Fatalf("store input: %v", err)
	}

	storedJob, err := job.New("job-2", job.OperationVideoConvert, time.Now().UTC())
	if err != nil {
		t.Fatalf("job.New() error = %v", err)
	}
	repository := &repositoryStub{record: task.Record{
		Job:      storedJob,
		InputKey: "uploads/job-2/input",
	}}
	processor := &processorStub{}
	runner := New(
		slog.New(slog.NewTextHandler(io.Discard, nil)),
		repository,
		queueStub{},
		store,
		proberStub{err: media.ErrInvalidMedia},
		processor,
		time.Hour,
		1024,
	)

	err = runner.process(context.Background(), "job-2")
	if err == nil {
		t.Fatal("process() error = nil, want error")
	}
	if processor.called {
		t.Fatal("processor was called for invalid media")
	}
	if got := repository.record.Job.Status(); got != job.StatusFailed {
		t.Fatalf("status = %q, want %q", got, job.StatusFailed)
	}
	snapshot := repository.record.Job.Snapshot()
	if snapshot.FailureCode != "invalid_media" {
		t.Fatalf("failure code = %q, want invalid_media", snapshot.FailureCode)
	}
}

func TestMaterializeInputCreatesPrivateFile(t *testing.T) {
	store, err := storage.NewLocal(t.TempDir())
	if err != nil {
		t.Fatalf("NewLocal() error = %v", err)
	}
	if _, err := store.Put(
		context.Background(),
		"uploads/job-3/input",
		bytes.NewBufferString("source"),
		1024,
	); err != nil {
		t.Fatalf("store input: %v", err)
	}

	runner := &Worker{storage: store}
	target := filepath.Join(t.TempDir(), "input")
	if err := runner.materializeInput(context.Background(), "uploads/job-3/input", target); err != nil {
		t.Fatalf("materializeInput() error = %v", err)
	}
	info, err := os.Stat(target)
	if err != nil {
		t.Fatalf("stat input: %v", err)
	}
	if got := info.Mode().Perm(); got != 0o600 {
		t.Fatalf("input mode = %o, want 600", got)
	}
}
