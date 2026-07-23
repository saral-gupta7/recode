package job

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

const maxFailureCodeLength = 64

var (
	ErrInvalidJobID             = errors.New("invalid job ID")
	ErrInvalidTime              = errors.New("invalid job time")
	ErrInvalidTransition        = errors.New("invalid job status transition")
	ErrInvalidStatus            = errors.New("invalid job status")
	ErrProgressUpdateNotAllowed = errors.New("job progress update not allowed")
	ErrInconsistentState        = errors.New("inconsistent job state")
	ErrInvalidAttempt           = errors.New("invalid job attempt")
	ErrExpirationNotAllowed     = errors.New("job expiration not allowed")
	ErrInvalidExpiration        = errors.New("invalid job expiration")
	ErrExpirationNotDue         = errors.New("job expiration not due")
	ErrOutcomeNotAllowed        = errors.New("job outcome not allowed")
	ErrInvalidResultKey         = errors.New("invalid job result key")
	ErrInvalidFailureCode       = errors.New("invalid job failure code")
)

// Snapshot is the persistence-facing representation used to rehydrate a Job.
// Raw numeric and string values are validated by Restore before entering the
// aggregate.
type Snapshot struct {
	ID          string
	Operation   Operation
	Status      Status
	Progress    int
	Attempt     int
	ExpiresAt   time.Time
	ResultKey   string
	FailureCode string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Job struct {
	id          string
	operation   Operation
	status      Status
	progress    Progress
	attempt     int
	expiresAt   time.Time
	resultKey   string
	failureCode string
	createdAt   time.Time
	updatedAt   time.Time
}

func New(id string, operation Operation, now time.Time) (*Job, error) {
	if strings.TrimSpace(id) == "" {
		return nil, ErrInvalidJobID
	}

	if !operation.Valid() {
		return nil, fmt.Errorf("%w: %q", ErrInvalidOperation, operation)
	}

	if now.IsZero() {
		return nil, ErrInvalidTime
	}

	now = now.UTC()

	return &Job{
		id:        id,
		operation: operation,
		status:    StatusQueued,
		progress:  Progress(0),
		attempt:   0,
		createdAt: now,
		updatedAt: now,
	}, nil
}

func Restore(snapshot Snapshot) (*Job, error) {
	if strings.TrimSpace(snapshot.ID) == "" {
		return nil, ErrInvalidJobID
	}

	if !snapshot.Operation.Valid() {
		return nil, fmt.Errorf("%w: %q", ErrInvalidOperation, snapshot.Operation)
	}

	if !snapshot.Status.Valid() {
		return nil, fmt.Errorf("%w: %q", ErrInvalidStatus, snapshot.Status)
	}

	if snapshot.CreatedAt.IsZero() ||
		snapshot.UpdatedAt.IsZero() ||
		snapshot.UpdatedAt.Before(snapshot.CreatedAt) {
		return nil, ErrInvalidTime
	}

	progress, err := NewProgress(snapshot.Progress)
	if err != nil {
		return nil, err
	}

	if snapshot.Attempt < 0 {
		return nil, fmt.Errorf("%w: %d", ErrInvalidAttempt, snapshot.Attempt)
	}

	if !validProgressForStatus(snapshot.Status, progress) ||
		!validAttemptForStatus(snapshot.Status, snapshot.Attempt) {
		return nil, fmt.Errorf(
			"%w: status %q, progress %d, attempt %d",
			ErrInconsistentState,
			snapshot.Status,
			progress.Value(),
			snapshot.Attempt,
		)
	}

	if err := validateSnapshotMetadata(snapshot); err != nil {
		return nil, err
	}

	return &Job{
		id:          snapshot.ID,
		operation:   snapshot.Operation,
		status:      snapshot.Status,
		progress:    progress,
		attempt:     snapshot.Attempt,
		expiresAt:   snapshot.ExpiresAt.UTC(),
		resultKey:   snapshot.ResultKey,
		failureCode: snapshot.FailureCode,
		createdAt:   snapshot.CreatedAt.UTC(),
		updatedAt:   snapshot.UpdatedAt.UTC(),
	}, nil
}

func (j *Job) TransitionTo(next Status, now time.Time) error {
	if !j.status.CanTransitionTo(next) {
		return fmt.Errorf(
			"%w: %q to %q",
			ErrInvalidTransition,
			j.status,
			next,
		)
	}

	if now.IsZero() || now.Before(j.updatedAt) {
		return ErrInvalidTime
	}

	if next == StatusExpired {
		if j.expiresAt.IsZero() || now.Before(j.expiresAt) {
			return ErrExpirationNotDue
		}
	}

	switch next {
	case StatusProcessing:
		j.attempt++

	case StatusCompleted:
		j.progress = Progress(100)

	case StatusQueued:
		j.progress = Progress(0)

	case StatusExpired:
		j.resultKey = ""
		j.failureCode = ""
	}

	j.status = next
	j.updatedAt = now.UTC()

	return nil
}

func (j *Job) UpdateProgress(value int, now time.Time) error {
	if j.status != StatusProcessing {
		return fmt.Errorf(
			"%w: status %q",
			ErrProgressUpdateNotAllowed,
			j.status,
		)
	}

	progress, err := NewProgress(value)
	if err != nil {
		return err
	}

	if progress.Value() == 100 || progress < j.progress {
		return fmt.Errorf(
			"%w: %d to %d",
			ErrProgressUpdateNotAllowed,
			j.progress.Value(),
			progress.Value(),
		)
	}

	if now.IsZero() || now.Before(j.updatedAt) {
		return ErrInvalidTime
	}

	j.progress = progress
	j.updatedAt = now.UTC()

	return nil
}

func (j *Job) RecordResult(resultKey string) error {
	if j.status != StatusCompleted {
		return fmt.Errorf("%w: status %q", ErrOutcomeNotAllowed, j.status)
	}
	if strings.TrimSpace(resultKey) == "" {
		return ErrInvalidResultKey
	}

	j.resultKey = resultKey
	j.failureCode = ""
	return nil
}

func (j *Job) RecordFailure(code string) error {
	if j.status != StatusFailed {
		return fmt.Errorf("%w: status %q", ErrOutcomeNotAllowed, j.status)
	}
	if !validFailureCode(code) {
		return fmt.Errorf("%w: %q", ErrInvalidFailureCode, code)
	}

	j.failureCode = code
	j.resultKey = ""
	return nil
}

func (j *Job) ScheduleExpiration(expiresAt time.Time) error {
	switch j.status {
	case StatusCompleted, StatusFailed, StatusCancelled:
	default:
		return fmt.Errorf("%w: status %q", ErrExpirationNotAllowed, j.status)
	}

	if expiresAt.IsZero() || !expiresAt.After(j.updatedAt) {
		return ErrInvalidExpiration
	}

	j.expiresAt = expiresAt.UTC()
	return nil
}

func (j *Job) ID() string {
	return j.id
}

func (j *Job) Operation() Operation {
	return j.operation
}

func (j *Job) Status() Status {
	return j.status
}

func (j *Job) Progress() Progress {
	return j.progress
}

func (j *Job) Attempt() int {
	return j.attempt
}

func (j *Job) ExpiresAt() (time.Time, bool) {
	if j.expiresAt.IsZero() {
		return time.Time{}, false
	}
	return j.expiresAt, true
}

func (j *Job) ResultKey() (string, bool) {
	return j.resultKey, j.resultKey != ""
}

func (j *Job) FailureCode() (string, bool) {
	return j.failureCode, j.failureCode != ""
}

func (j *Job) CreatedAt() time.Time {
	return j.createdAt
}

func (j *Job) UpdatedAt() time.Time {
	return j.updatedAt
}

func (j *Job) Snapshot() Snapshot {
	return Snapshot{
		ID:          j.id,
		Operation:   j.operation,
		Status:      j.status,
		Progress:    j.progress.Value(),
		Attempt:     j.attempt,
		ExpiresAt:   j.expiresAt,
		ResultKey:   j.resultKey,
		FailureCode: j.failureCode,
		CreatedAt:   j.createdAt,
		UpdatedAt:   j.updatedAt,
	}
}

func validProgressForStatus(status Status, progress Progress) bool {
	switch status {
	case StatusQueued:
		return progress.Value() == 0
	case StatusProcessing,
		StatusCancelling,
		StatusFailed,
		StatusCancelled:
		return progress.Value() < 100
	case StatusCompleted:
		return progress.Value() == 100
	case StatusExpired:
		return true
	default:
		return false
	}
}

func validAttemptForStatus(status Status, attempt int) bool {
	switch status {
	case StatusProcessing, StatusCancelling, StatusCompleted, StatusFailed:
		return attempt > 0
	default:
		return true
	}
}

func validateSnapshotMetadata(snapshot Snapshot) error {
	hasExpiration := !snapshot.ExpiresAt.IsZero()
	hasResult := snapshot.ResultKey != ""
	hasFailure := snapshot.FailureCode != ""

	switch snapshot.Status {
	case StatusQueued, StatusProcessing, StatusCancelling:
		if hasExpiration || hasResult || hasFailure {
			return ErrInconsistentState
		}

	case StatusCompleted:
		if !hasExpiration ||
			!snapshot.ExpiresAt.After(snapshot.UpdatedAt) ||
			!hasResult ||
			hasFailure {
			return ErrInconsistentState
		}
		if strings.TrimSpace(snapshot.ResultKey) == "" {
			return ErrInvalidResultKey
		}

	case StatusFailed:
		if !hasExpiration ||
			!snapshot.ExpiresAt.After(snapshot.UpdatedAt) ||
			hasResult ||
			!hasFailure {
			return ErrInconsistentState
		}
		if !validFailureCode(snapshot.FailureCode) {
			return fmt.Errorf("%w: %q", ErrInvalidFailureCode, snapshot.FailureCode)
		}

	case StatusCancelled:
		if !hasExpiration ||
			!snapshot.ExpiresAt.After(snapshot.UpdatedAt) ||
			hasResult ||
			hasFailure {
			return ErrInconsistentState
		}

	case StatusExpired:
		if !hasExpiration ||
			snapshot.UpdatedAt.Before(snapshot.ExpiresAt) ||
			hasResult ||
			hasFailure {
			return ErrInconsistentState
		}
	}

	return nil
}

func validFailureCode(code string) bool {
	if code == "" || len(code) > maxFailureCodeLength {
		return false
	}

	for _, character := range code {
		if (character >= 'a' && character <= 'z') ||
			(character >= '0' && character <= '9') ||
			character == '_' {
			continue
		}
		return false
	}

	return true
}
