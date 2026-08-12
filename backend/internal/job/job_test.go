package job

import (
	"errors"
	"strings"
	"testing"
	"time"
)

var testNow = time.Date(2026, time.July, 23, 9, 0, 0, 0, time.UTC)

func TestNew(t *testing.T) {
	location := time.FixedZone("test-zone", 5*60*60+30*60)
	now := time.Date(2026, time.July, 23, 14, 30, 0, 0, location)

	got, err := New("job-123", OperationImageCrop, now)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	if got.ID() != "job-123" {
		t.Errorf("ID() = %q, want job-123", got.ID())
	}
	if got.Operation() != OperationImageCrop {
		t.Errorf("Operation() = %q, want %q", got.Operation(), OperationImageCrop)
	}
	if got.Status() != StatusQueued {
		t.Errorf("Status() = %q, want %q", got.Status(), StatusQueued)
	}
	if got.Progress().Value() != 0 {
		t.Errorf("Progress().Value() = %d, want 0", got.Progress().Value())
	}
	if got.Attempt() != 0 {
		t.Errorf("Attempt() = %d, want 0", got.Attempt())
	}
	if _, ok := got.ExpiresAt(); ok {
		t.Error("ExpiresAt() scheduled = true, want false")
	}
	if _, ok := got.ResultKey(); ok {
		t.Error("ResultKey() present = true, want false")
	}
	if _, ok := got.FailureCode(); ok {
		t.Error("FailureCode() present = true, want false")
	}
	if !got.CreatedAt().Equal(now) || got.CreatedAt().Location() != time.UTC {
		t.Errorf("CreatedAt() = %v, want instant %v in UTC", got.CreatedAt(), now)
	}
	if !got.UpdatedAt().Equal(got.CreatedAt()) {
		t.Errorf("UpdatedAt() = %v, want CreatedAt() %v", got.UpdatedAt(), got.CreatedAt())
	}
}

func TestNewRejectsInvalidInput(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		operation Operation
		now       time.Time
		wantErr   error
	}{
		{
			name:      "empty ID",
			operation: OperationImageResize,
			now:       testNow,
			wantErr:   ErrInvalidJobID,
		},
		{
			name:      "whitespace ID",
			id:        " \t\n ",
			operation: OperationImageResize,
			now:       testNow,
			wantErr:   ErrInvalidJobID,
		},
		{
			name:      "invalid operation",
			id:        "job-123",
			operation: Operation("unknown"),
			now:       testNow,
			wantErr:   ErrInvalidOperation,
		},
		{
			name:      "zero time",
			id:        "job-123",
			operation: OperationImageResize,
			wantErr:   ErrInvalidTime,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := New(test.id, test.operation, test.now)
			if got != nil {
				t.Fatalf("New() job = %#v, want nil", got)
			}
			if !errors.Is(err, test.wantErr) {
				t.Fatalf("New() error = %v, want errors.Is(_, %v)", err, test.wantErr)
			}
		})
	}
}

func TestRestoreAcceptsValidLifecycleSnapshots(t *testing.T) {
	statuses := []Status{
		StatusQueued,
		StatusProcessing,
		StatusCancelling,
		StatusCompleted,
		StatusFailed,
		StatusCancelled,
		StatusExpired,
	}

	for _, status := range statuses {
		t.Run(string(status), func(t *testing.T) {
			snapshot := validSnapshot(status)
			got, err := Restore(snapshot)
			if err != nil {
				t.Fatalf("Restore() error = %v", err)
			}

			if got.ID() != snapshot.ID ||
				got.Operation() != snapshot.Operation ||
				got.Status() != snapshot.Status ||
				got.Progress().Value() != snapshot.Progress ||
				got.Attempt() != snapshot.Attempt {
				t.Fatalf("Restore() did not preserve snapshot: got %#v, snapshot %#v", got, snapshot)
			}
			if !got.CreatedAt().Equal(snapshot.CreatedAt) ||
				!got.UpdatedAt().Equal(snapshot.UpdatedAt) {
				t.Error("Restore() did not preserve timestamp instants")
			}
			if got.CreatedAt().Location() != time.UTC ||
				got.UpdatedAt().Location() != time.UTC {
				t.Error("Restore() did not normalize timestamps to UTC")
			}

			resultKey, hasResult := got.ResultKey()
			assertOptionalValue(t, "result key", resultKey, hasResult, snapshot.ResultKey)
			failureCode, hasFailure := got.FailureCode()
			assertOptionalValue(t, "failure code", failureCode, hasFailure, snapshot.FailureCode)

			expiresAt, scheduled := got.ExpiresAt()
			if scheduled != !snapshot.ExpiresAt.IsZero() {
				t.Errorf("ExpiresAt() scheduled = %t, want %t", scheduled, !snapshot.ExpiresAt.IsZero())
			}
			if scheduled && !expiresAt.Equal(snapshot.ExpiresAt) {
				t.Errorf("ExpiresAt() = %v, want %v", expiresAt, snapshot.ExpiresAt)
			}
		})
	}
}

func TestRestoreRejectsInvalidSnapshot(t *testing.T) {
	longFailureCode := strings.Repeat("a", maxFailureCodeLength+1)
	tests := []struct {
		name     string
		snapshot Snapshot
		wantErr  error
	}{
		{
			name: "empty ID",
			snapshot: withSnapshot(StatusQueued, func(s *Snapshot) {
				s.ID = ""
			}),
			wantErr: ErrInvalidJobID,
		},
		{
			name: "invalid operation",
			snapshot: withSnapshot(StatusQueued, func(s *Snapshot) {
				s.Operation = Operation("unknown")
			}),
			wantErr: ErrInvalidOperation,
		},
		{
			name: "invalid status",
			snapshot: withSnapshot(StatusQueued, func(s *Snapshot) {
				s.Status = Status("unknown")
			}),
			wantErr: ErrInvalidStatus,
		},
		{
			name: "zero creation time",
			snapshot: withSnapshot(StatusQueued, func(s *Snapshot) {
				s.CreatedAt = time.Time{}
			}),
			wantErr: ErrInvalidTime,
		},
		{
			name: "update before creation",
			snapshot: withSnapshot(StatusQueued, func(s *Snapshot) {
				s.UpdatedAt = s.CreatedAt.Add(-time.Nanosecond)
			}),
			wantErr: ErrInvalidTime,
		},
		{
			name: "negative progress",
			snapshot: withSnapshot(StatusProcessing, func(s *Snapshot) {
				s.Progress = -1
			}),
			wantErr: ErrInvalidProgress,
		},
		{
			name: "progress above one hundred",
			snapshot: withSnapshot(StatusProcessing, func(s *Snapshot) {
				s.Progress = 101
			}),
			wantErr: ErrInvalidProgress,
		},
		{
			name: "negative attempt",
			snapshot: withSnapshot(StatusQueued, func(s *Snapshot) {
				s.Attempt = -1
			}),
			wantErr: ErrInvalidAttempt,
		},
		{
			name: "processing without attempt",
			snapshot: withSnapshot(StatusProcessing, func(s *Snapshot) {
				s.Attempt = 0
			}),
			wantErr: ErrInconsistentState,
		},
		{
			name: "queued with progress",
			snapshot: withSnapshot(StatusQueued, func(s *Snapshot) {
				s.Progress = 1
			}),
			wantErr: ErrInconsistentState,
		},
		{
			name: "active job with expiration",
			snapshot: withSnapshot(StatusProcessing, func(s *Snapshot) {
				s.ExpiresAt = s.UpdatedAt.Add(time.Hour)
			}),
			wantErr: ErrInconsistentState,
		},
		{
			name: "completed without result",
			snapshot: withSnapshot(StatusCompleted, func(s *Snapshot) {
				s.ResultKey = ""
			}),
			wantErr: ErrInconsistentState,
		},
		{
			name: "completed with failure",
			snapshot: withSnapshot(StatusCompleted, func(s *Snapshot) {
				s.FailureCode = "processing_failed"
			}),
			wantErr: ErrInconsistentState,
		},
		{
			name: "failed without failure code",
			snapshot: withSnapshot(StatusFailed, func(s *Snapshot) {
				s.FailureCode = ""
			}),
			wantErr: ErrInconsistentState,
		},
		{
			name: "invalid failure characters",
			snapshot: withSnapshot(StatusFailed, func(s *Snapshot) {
				s.FailureCode = "FFmpeg failed!"
			}),
			wantErr: ErrInvalidFailureCode,
		},
		{
			name: "failure code too long",
			snapshot: withSnapshot(StatusFailed, func(s *Snapshot) {
				s.FailureCode = longFailureCode
			}),
			wantErr: ErrInvalidFailureCode,
		},
		{
			name: "terminal expiration not after update",
			snapshot: withSnapshot(StatusCancelled, func(s *Snapshot) {
				s.ExpiresAt = s.UpdatedAt
			}),
			wantErr: ErrInconsistentState,
		},
		{
			name: "expired before scheduled time",
			snapshot: withSnapshot(StatusExpired, func(s *Snapshot) {
				s.ExpiresAt = s.UpdatedAt.Add(time.Nanosecond)
			}),
			wantErr: ErrInconsistentState,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := Restore(test.snapshot)
			if got != nil {
				t.Fatalf("Restore() job = %#v, want nil", got)
			}
			if !errors.Is(err, test.wantErr) {
				t.Fatalf("Restore() error = %v, want errors.Is(_, %v)", err, test.wantErr)
			}
		})
	}
}

func TestJobTransitionAndAttemptLifecycle(t *testing.T) {
	job := mustNewJob(t)

	if err := job.TransitionTo(StatusProcessing, testNow.Add(time.Minute)); err != nil {
		t.Fatalf("first TransitionTo(processing) error = %v", err)
	}
	if job.Attempt() != 1 {
		t.Fatalf("Attempt() after first claim = %d, want 1", job.Attempt())
	}

	if err := job.UpdateProgress(40, testNow.Add(2*time.Minute)); err != nil {
		t.Fatalf("UpdateProgress() error = %v", err)
	}
	if err := job.TransitionTo(StatusQueued, testNow.Add(3*time.Minute)); err != nil {
		t.Fatalf("TransitionTo(queued) error = %v", err)
	}
	if job.Progress().Value() != 0 || job.Attempt() != 1 {
		t.Fatalf(
			"after retry: progress = %d, attempt = %d; want 0, 1",
			job.Progress().Value(),
			job.Attempt(),
		)
	}

	if err := job.TransitionTo(StatusProcessing, testNow.Add(4*time.Minute)); err != nil {
		t.Fatalf("second TransitionTo(processing) error = %v", err)
	}
	if job.Attempt() != 2 {
		t.Errorf("Attempt() after second claim = %d, want 2", job.Attempt())
	}
}

func TestJobTransitionFailureDoesNotMutate(t *testing.T) {
	job := mustNewJob(t)
	originalUpdatedAt := job.UpdatedAt()

	err := job.TransitionTo(StatusCompleted, testNow.Add(time.Minute))

	if !errors.Is(err, ErrInvalidTransition) {
		t.Fatalf("TransitionTo() error = %v, want ErrInvalidTransition", err)
	}
	if job.Status() != StatusQueued ||
		job.Progress().Value() != 0 ||
		job.Attempt() != 0 ||
		!job.UpdatedAt().Equal(originalUpdatedAt) {
		t.Error("rejected transition mutated the job")
	}

	err = job.TransitionTo(StatusProcessing, testNow.Add(-time.Nanosecond))
	if !errors.Is(err, ErrInvalidTime) {
		t.Fatalf("TransitionTo() error = %v, want ErrInvalidTime", err)
	}
	if job.Attempt() != 0 {
		t.Error("invalid-time transition incremented attempt")
	}
}

func TestJobUpdateProgress(t *testing.T) {
	job := mustProcessingJob(t)
	location := time.FixedZone("test-zone", 5*60*60+30*60)
	updatedAt := time.Date(2026, time.July, 23, 14, 32, 0, 0, location)

	if err := job.UpdateProgress(42, updatedAt); err != nil {
		t.Fatalf("UpdateProgress() error = %v", err)
	}
	if job.Progress().Value() != 42 {
		t.Errorf("Progress().Value() = %d, want 42", job.Progress().Value())
	}
	if !job.UpdatedAt().Equal(updatedAt) || job.UpdatedAt().Location() != time.UTC {
		t.Errorf("UpdatedAt() = %v, want instant %v in UTC", job.UpdatedAt(), updatedAt)
	}
}

func TestJobUpdateProgressRejectsInvalidUpdateWithoutMutation(t *testing.T) {
	tests := []struct {
		name    string
		value   int
		now     time.Time
		wantErr error
	}{
		{name: "negative", value: -1, now: testNow.Add(2 * time.Minute), wantErr: ErrInvalidProgress},
		{name: "above one hundred", value: 101, now: testNow.Add(2 * time.Minute), wantErr: ErrInvalidProgress},
		{name: "one hundred reserved for completion", value: 100, now: testNow.Add(2 * time.Minute), wantErr: ErrProgressUpdateNotAllowed},
		{name: "backward progress", value: 39, now: testNow.Add(2 * time.Minute), wantErr: ErrProgressUpdateNotAllowed},
		{name: "zero time", value: 50, wantErr: ErrInvalidTime},
		{name: "backward time", value: 50, now: testNow, wantErr: ErrInvalidTime},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			job := mustProcessingJob(t)
			if err := job.UpdateProgress(40, testNow.Add(90*time.Second)); err != nil {
				t.Fatalf("setup UpdateProgress() error = %v", err)
			}
			originalUpdatedAt := job.UpdatedAt()

			err := job.UpdateProgress(test.value, test.now)

			if !errors.Is(err, test.wantErr) {
				t.Fatalf("UpdateProgress() error = %v, want errors.Is(_, %v)", err, test.wantErr)
			}
			if job.Progress().Value() != 40 || !job.UpdatedAt().Equal(originalUpdatedAt) {
				t.Error("rejected progress update mutated the job")
			}
		})
	}
}

func TestJobUpdateProgressRejectsWrongStatus(t *testing.T) {
	job := mustNewJob(t)

	err := job.UpdateProgress(20, testNow.Add(time.Minute))

	if !errors.Is(err, ErrProgressUpdateNotAllowed) {
		t.Fatalf("UpdateProgress() error = %v, want ErrProgressUpdateNotAllowed", err)
	}
}

func TestCompletedJobOutcomeAndExpiration(t *testing.T) {
	job := mustProcessingJob(t)
	completedAt := testNow.Add(2 * time.Minute)
	if err := job.TransitionTo(StatusCompleted, completedAt); err != nil {
		t.Fatalf("TransitionTo(completed) error = %v", err)
	}
	if job.Progress().Value() != 100 {
		t.Errorf("completed progress = %d, want 100", job.Progress().Value())
	}

	if err := job.RecordResult("results/job-123/output.mp4"); err != nil {
		t.Fatalf("RecordResult() error = %v", err)
	}
	resultKey, ok := job.ResultKey()
	if !ok || resultKey != "results/job-123/output.mp4" {
		t.Errorf("ResultKey() = %q, %t", resultKey, ok)
	}

	expiresAt := completedAt.Add(time.Hour)
	if err := job.ScheduleExpiration(expiresAt); err != nil {
		t.Fatalf("ScheduleExpiration() error = %v", err)
	}

	if err := job.TransitionTo(StatusExpired, expiresAt.Add(-time.Nanosecond)); !errors.Is(err, ErrExpirationNotDue) {
		t.Fatalf("early TransitionTo(expired) error = %v, want ErrExpirationNotDue", err)
	}
	if err := job.TransitionTo(StatusExpired, expiresAt); err != nil {
		t.Fatalf("TransitionTo(expired) error = %v", err)
	}
	if _, ok := job.ResultKey(); ok {
		t.Error("expired job retained result key")
	}
}

func TestFailedJobOutcomeAndExpiration(t *testing.T) {
	job := mustProcessingJob(t)
	failedAt := testNow.Add(2 * time.Minute)
	if err := job.TransitionTo(StatusFailed, failedAt); err != nil {
		t.Fatalf("TransitionTo(failed) error = %v", err)
	}
	if err := job.RecordFailure("processing_failed"); err != nil {
		t.Fatalf("RecordFailure() error = %v", err)
	}
	if code, ok := job.FailureCode(); !ok || code != "processing_failed" {
		t.Errorf("FailureCode() = %q, %t", code, ok)
	}
	if err := job.ScheduleExpiration(failedAt.Add(time.Hour)); err != nil {
		t.Fatalf("ScheduleExpiration() error = %v", err)
	}
}

func TestJobRejectsInvalidOutcome(t *testing.T) {
	queued := mustNewJob(t)
	if err := queued.RecordResult("result.mp4"); !errors.Is(err, ErrOutcomeNotAllowed) {
		t.Errorf("queued RecordResult() error = %v, want ErrOutcomeNotAllowed", err)
	}
	if err := queued.RecordFailure("processing_failed"); !errors.Is(err, ErrOutcomeNotAllowed) {
		t.Errorf("queued RecordFailure() error = %v, want ErrOutcomeNotAllowed", err)
	}

	completed := mustProcessingJob(t)
	if err := completed.TransitionTo(StatusCompleted, testNow.Add(2*time.Minute)); err != nil {
		t.Fatalf("TransitionTo(completed) error = %v", err)
	}
	if err := completed.RecordResult(" \t "); !errors.Is(err, ErrInvalidResultKey) {
		t.Errorf("RecordResult() error = %v, want ErrInvalidResultKey", err)
	}

	failed := mustProcessingJob(t)
	if err := failed.TransitionTo(StatusFailed, testNow.Add(2*time.Minute)); err != nil {
		t.Fatalf("TransitionTo(failed) error = %v", err)
	}
	for _, code := range []string{"", "FFmpeg failed!", strings.Repeat("a", maxFailureCodeLength+1)} {
		if err := failed.RecordFailure(code); !errors.Is(err, ErrInvalidFailureCode) {
			t.Errorf("RecordFailure(%q) error = %v, want ErrInvalidFailureCode", code, err)
		}
	}
}

func TestJobRejectsInvalidExpiration(t *testing.T) {
	queued := mustNewJob(t)
	if err := queued.ScheduleExpiration(testNow.Add(time.Hour)); !errors.Is(err, ErrExpirationNotAllowed) {
		t.Errorf("queued ScheduleExpiration() error = %v, want ErrExpirationNotAllowed", err)
	}

	completed := mustProcessingJob(t)
	completedAt := testNow.Add(2 * time.Minute)
	if err := completed.TransitionTo(StatusCompleted, completedAt); err != nil {
		t.Fatalf("TransitionTo(completed) error = %v", err)
	}
	for _, expiresAt := range []time.Time{{}, completedAt, completedAt.Add(-time.Nanosecond)} {
		if err := completed.ScheduleExpiration(expiresAt); !errors.Is(err, ErrInvalidExpiration) {
			t.Errorf("ScheduleExpiration(%v) error = %v, want ErrInvalidExpiration", expiresAt, err)
		}
	}
}

func validSnapshot(status Status) Snapshot {
	snapshot := Snapshot{
		ID:        "job-123",
		Operation: OperationImageCrop,
		Status:    status,
		CreatedAt: testNow,
		UpdatedAt: testNow.Add(time.Minute),
	}

	switch status {
	case StatusQueued:
	case StatusProcessing:
		snapshot.Progress = 40
		snapshot.Attempt = 1
	case StatusCancelling:
		snapshot.Progress = 40
		snapshot.Attempt = 1
	case StatusCompleted:
		snapshot.Progress = 100
		snapshot.Attempt = 1
		snapshot.ResultKey = "results/job-123/output.mp4"
		snapshot.ExpiresAt = snapshot.UpdatedAt.Add(time.Hour)
	case StatusFailed:
		snapshot.Progress = 40
		snapshot.Attempt = 1
		snapshot.FailureCode = "processing_failed"
		snapshot.ExpiresAt = snapshot.UpdatedAt.Add(time.Hour)
	case StatusCancelled:
		snapshot.Progress = 40
		snapshot.Attempt = 1
		snapshot.ExpiresAt = snapshot.UpdatedAt.Add(time.Hour)
	case StatusExpired:
		snapshot.Progress = 100
		snapshot.Attempt = 1
		snapshot.ExpiresAt = snapshot.UpdatedAt.Add(-time.Minute)
	}

	return snapshot
}

func withSnapshot(status Status, mutate func(*Snapshot)) Snapshot {
	snapshot := validSnapshot(status)
	mutate(&snapshot)
	return snapshot
}

func mustNewJob(t *testing.T) *Job {
	t.Helper()
	job, err := New("job-123", OperationImageCrop, testNow)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	return job
}

func mustProcessingJob(t *testing.T) *Job {
	t.Helper()
	job := mustNewJob(t)
	if err := job.TransitionTo(StatusProcessing, testNow.Add(time.Minute)); err != nil {
		t.Fatalf("TransitionTo(processing) error = %v", err)
	}
	return job
}

func assertOptionalValue(t *testing.T, name string, got string, present bool, want string) {
	t.Helper()
	if present != (want != "") {
		t.Errorf("%s present = %t, want %t", name, present, want != "")
	}
	if got != want {
		t.Errorf("%s = %q, want %q", name, got, want)
	}
}
