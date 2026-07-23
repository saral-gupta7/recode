package postgres_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/saral-gupta7/recode/backend/internal/database"
	"github.com/saral-gupta7/recode/backend/internal/job"
	"github.com/saral-gupta7/recode/backend/internal/repository/postgres"
	"github.com/saral-gupta7/recode/backend/internal/task"
)

func TestTaskRepositoryCreateFindAndSave(t *testing.T) {
	pool, repository := openTestTaskRepository(t)
	id := newTestJobID()
	deleteJobAfterTest(t, pool, id)

	now := time.Date(2026, time.July, 23, 9, 0, 0, 0, time.UTC)
	storedJob, err := job.New(id, job.OperationImageResize, now)
	if err != nil {
		t.Fatalf("job.New() error = %v", err)
	}
	record, err := task.NewRecord(
		storedJob,
		"uploads/"+id+"/input",
		"photo.png",
		"image/png",
		128,
		task.Options{Width: 640},
		[]byte("owner-token-hash"),
	)
	if err != nil {
		t.Fatalf("task.NewRecord() error = %v", err)
	}

	if err := repository.Create(context.Background(), record); err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	found, err := repository.FindByID(context.Background(), id)
	if err != nil {
		t.Fatalf("FindByID() error = %v", err)
	}
	assertSnapshotsEqual(t, found.Job.Snapshot(), record.Job.Snapshot())
	if found.InputKey != record.InputKey ||
		found.OriginalFilename != record.OriginalFilename ||
		found.MIMEType != record.MIMEType ||
		found.SizeBytes != record.SizeBytes ||
		found.Options != record.Options ||
		string(found.OwnerTokenHash) != string(record.OwnerTokenHash) {
		t.Fatalf("execution metadata = %#v, want %#v", found, record)
	}

	if err := found.Job.TransitionTo(job.StatusProcessing, now.Add(time.Minute)); err != nil {
		t.Fatalf("TransitionTo() error = %v", err)
	}
	if err := repository.SaveJob(context.Background(), found.Job); err != nil {
		t.Fatalf("SaveJob() error = %v", err)
	}
	saved, err := repository.FindByID(context.Background(), id)
	if err != nil {
		t.Fatalf("FindByID() after save error = %v", err)
	}
	if saved.Job.Status() != job.StatusProcessing {
		t.Fatalf("saved status = %q, want %q", saved.Job.Status(), job.StatusProcessing)
	}
}

func openTestTaskRepository(t *testing.T) (*pgxpool.Pool, *postgres.TaskRepository) {
	t.Helper()

	databaseURL := os.Getenv("RECODE_TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("RECODE_TEST_DATABASE_URL is not set")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := database.Open(ctx, databaseURL)
	if err != nil {
		t.Fatalf("database.Open() error = %v", err)
	}
	t.Cleanup(pool.Close)
	return pool, postgres.NewTaskRepository(pool)
}
