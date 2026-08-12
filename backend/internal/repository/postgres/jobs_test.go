package postgres_test

import (
	"context"
	"crypto/rand"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/saral-gupta7/recode/backend/internal/database"
	"github.com/saral-gupta7/recode/backend/internal/job"
	"github.com/saral-gupta7/recode/backend/internal/repository/postgres"
)

func TestJobRepositoryCreateAndFindQueuedJob(t *testing.T) {
	pool, repository := openTestRepository(t)
	id := newTestJobID()
	deleteJobAfterTest(t, pool, id)

	createdAt := time.Date(2026, time.July, 23, 9, 0, 0, 0, time.UTC)
	createdJob, err := job.New(id, job.OperationImageGrayscale, createdAt)
	if err != nil {
		t.Fatalf("job.New() error = %v", err)
	}

	if err := repository.Create(context.Background(), createdJob); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	foundJob, err := repository.FindByID(context.Background(), id)
	if err != nil {
		t.Fatalf("FindByID() error = %v", err)
	}

	assertSnapshotsEqual(t, foundJob.Snapshot(), createdJob.Snapshot())
}

func TestJobRepositoryCreateAndFindCompletedJob(t *testing.T) {
	pool, repository := openTestRepository(t)
	id := newTestJobID()
	deleteJobAfterTest(t, pool, id)

	now := time.Date(2026, time.July, 23, 9, 0, 0, 0, time.UTC)
	completedJob, err := job.New(id, job.OperationImageCrop, now)
	if err != nil {
		t.Fatalf("job.New() error = %v", err)
	}
	if err := completedJob.TransitionTo(job.StatusProcessing, now.Add(time.Minute)); err != nil {
		t.Fatalf("TransitionTo(processing) error = %v", err)
	}
	if err := completedJob.UpdateProgress(60, now.Add(2*time.Minute)); err != nil {
		t.Fatalf("UpdateProgress() error = %v", err)
	}
	if err := completedJob.TransitionTo(job.StatusCompleted, now.Add(3*time.Minute)); err != nil {
		t.Fatalf("TransitionTo(completed) error = %v", err)
	}
	if err := completedJob.RecordResult("results/" + id + "/output.mp4"); err != nil {
		t.Fatalf("RecordResult() error = %v", err)
	}
	if err := completedJob.ScheduleExpiration(now.Add(2 * time.Hour)); err != nil {
		t.Fatalf("ScheduleExpiration() error = %v", err)
	}

	if err := repository.Create(context.Background(), completedJob); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	foundJob, err := repository.FindByID(context.Background(), id)
	if err != nil {
		t.Fatalf("FindByID() error = %v", err)
	}

	assertSnapshotsEqual(t, foundJob.Snapshot(), completedJob.Snapshot())
}

func TestJobRepositorySave(t *testing.T) {
	pool, repository := openTestRepository(t)
	id := newTestJobID()
	deleteJobAfterTest(t, pool, id)

	now := time.Date(2026, time.July, 23, 9, 0, 0, 0, time.UTC)
	storedJob, err := job.New(id, job.OperationImageResize, now)
	if err != nil {
		t.Fatalf("job.New() error = %v", err)
	}
	if err := repository.Create(context.Background(), storedJob); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if err := storedJob.TransitionTo(job.StatusProcessing, now.Add(time.Minute)); err != nil {
		t.Fatalf("TransitionTo() error = %v", err)
	}
	if err := storedJob.UpdateProgress(35, now.Add(2*time.Minute)); err != nil {
		t.Fatalf("UpdateProgress() error = %v", err)
	}
	if err := repository.Save(context.Background(), storedJob); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	foundJob, err := repository.FindByID(context.Background(), id)
	if err != nil {
		t.Fatalf("FindByID() error = %v", err)
	}

	assertSnapshotsEqual(t, foundJob.Snapshot(), storedJob.Snapshot())
}

func TestJobRepositoryReturnsNotFound(t *testing.T) {
	_, repository := openTestRepository(t)
	id := newTestJobID()

	foundJob, err := repository.FindByID(context.Background(), id)
	if foundJob != nil {
		t.Fatalf("FindByID() job = %#v, want nil", foundJob)
	}
	if !errors.Is(err, job.ErrNotFound) {
		t.Fatalf("FindByID() error = %v, want job.ErrNotFound", err)
	}

	missingJob, err := job.New(
		id,
		job.OperationImageCompress,
		time.Date(2026, time.July, 23, 9, 0, 0, 0, time.UTC),
	)
	if err != nil {
		t.Fatalf("job.New() error = %v", err)
	}
	if err := repository.Save(context.Background(), missingJob); !errors.Is(err, job.ErrNotFound) {
		t.Fatalf("Save() error = %v, want job.ErrNotFound", err)
	}
}

func TestJobRepositoryRejectsNilJob(t *testing.T) {
	_, repository := openTestRepository(t)

	if err := repository.Create(context.Background(), nil); !errors.Is(err, postgres.ErrNilJob) {
		t.Errorf("Create(nil) error = %v, want postgres.ErrNilJob", err)
	}
	if err := repository.Save(context.Background(), nil); !errors.Is(err, postgres.ErrNilJob) {
		t.Errorf("Save(nil) error = %v, want postgres.ErrNilJob", err)
	}
}

func openTestRepository(t *testing.T) (*pgxpool.Pool, *postgres.JobRepository) {
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

	return pool, postgres.NewJobRepository(pool)
}

func deleteJobAfterTest(t *testing.T, pool *pgxpool.Pool, id string) {
	t.Helper()

	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if _, err := pool.Exec(ctx, "DELETE FROM jobs WHERE id = $1", id); err != nil {
			t.Errorf("delete test job %q: %v", id, err)
		}
	})
}

func newTestJobID() string {
	return "test-" + rand.Text()
}

func assertSnapshotsEqual(t *testing.T, got job.Snapshot, want job.Snapshot) {
	t.Helper()

	if got.ID != want.ID ||
		got.Operation != want.Operation ||
		got.Status != want.Status ||
		got.Progress != want.Progress ||
		got.Attempt != want.Attempt ||
		got.ResultKey != want.ResultKey ||
		got.FailureCode != want.FailureCode {
		t.Errorf("snapshot = %#v, want %#v", got, want)
	}

	if !got.ExpiresAt.Equal(want.ExpiresAt) ||
		!got.CreatedAt.Equal(want.CreatedAt) ||
		!got.UpdatedAt.Equal(want.UpdatedAt) {
		t.Errorf("snapshot times = %#v, want %#v", got, want)
	}
}
