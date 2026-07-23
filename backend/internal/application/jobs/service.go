package jobs

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/saral-gupta7/recode/backend/internal/job"
	"github.com/saral-gupta7/recode/backend/internal/queue"
	"github.com/saral-gupta7/recode/backend/internal/storage"
	"github.com/saral-gupta7/recode/backend/internal/task"
)

var (
	ErrUnauthorized = errors.New("job ownership token is invalid")
	ErrNotReady     = errors.New("job result is not ready")
	ErrConflict     = errors.New("job state conflicts with the requested action")
)

type Service struct {
	repository task.Repository
	queue      queue.Queue
	storage    storage.Store
	maxBytes   int64
	retention  time.Duration
	now        func() time.Time
}

type CreateInput struct {
	Operation        job.Operation
	Options          task.Options
	OriginalFilename string
	MIMEType         string
	Content          io.Reader
}

type Created struct {
	Record     task.Record
	OwnerToken string
}

func New(
	repository task.Repository,
	jobQueue queue.Queue,
	store storage.Store,
	maxBytes int64,
	retention time.Duration,
) *Service {
	return &Service{
		repository: repository,
		queue:      jobQueue,
		storage:    store,
		maxBytes:   maxBytes,
		retention:  retention,
		now:        time.Now,
	}
}

func (s *Service) Create(ctx context.Context, input CreateInput) (Created, error) {
	if !input.Operation.Valid() || input.Content == nil {
		return Created{}, task.ErrInvalidInput
	}

	options, err := task.NormalizeOptions(input.Operation, input.Options)
	if err != nil {
		return Created{}, err
	}

	jobID := rand.Text()
	ownerToken := rand.Text()
	ownerHash := sha256.Sum256([]byte(ownerToken))
	inputKey := "uploads/" + jobID + "/input"

	size, err := s.storage.Put(ctx, inputKey, input.Content, s.maxBytes)
	if err != nil {
		return Created{}, fmt.Errorf("store upload: %w", err)
	}
	removeInput := true
	defer func() {
		if removeInput {
			_ = s.storage.Delete(context.Background(), inputKey)
		}
	}()

	storedJob, err := job.New(jobID, input.Operation, s.now().UTC())
	if err != nil {
		return Created{}, fmt.Errorf("create job: %w", err)
	}

	record, err := task.NewRecord(
		storedJob,
		inputKey,
		filepath.Base(input.OriginalFilename),
		input.MIMEType,
		size,
		options,
		ownerHash[:],
	)
	if err != nil {
		return Created{}, err
	}

	if err := s.repository.Create(ctx, record); err != nil {
		return Created{}, fmt.Errorf("persist job: %w", err)
	}

	if err := s.queue.Enqueue(ctx, jobID); err != nil {
		_ = s.repository.Delete(context.Background(), jobID)
		return Created{}, fmt.Errorf("enqueue job: %w", err)
	}

	removeInput = false
	return Created{Record: record, OwnerToken: ownerToken}, nil
}

func (s *Service) Get(ctx context.Context, id string, token string) (task.Record, error) {
	record, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return task.Record{}, err
	}
	if !owns(record, token) {
		return task.Record{}, ErrUnauthorized
	}
	return record, nil
}

func (s *Service) Cancel(ctx context.Context, id string, token string) (task.Record, error) {
	record, err := s.Get(ctx, id, token)
	if err != nil {
		return task.Record{}, err
	}

	now := s.now().UTC()
	switch record.Job.Status() {
	case job.StatusQueued:
		if err := record.Job.TransitionTo(job.StatusCancelled, now); err != nil {
			return task.Record{}, err
		}
		if err := record.Job.ScheduleExpiration(task.ResultExpiration(now, s.retention)); err != nil {
			return task.Record{}, err
		}

	case job.StatusProcessing:
		if err := record.Job.TransitionTo(job.StatusCancelling, now); err != nil {
			return task.Record{}, err
		}

	case job.StatusCancelling, job.StatusCancelled:
		return record, nil

	default:
		return task.Record{}, ErrConflict
	}

	if err := s.repository.SaveJob(ctx, record.Job); err != nil {
		return task.Record{}, err
	}
	return record, nil
}

func (s *Service) OpenResult(
	ctx context.Context,
	id string,
	token string,
) (task.Record, io.ReadCloser, error) {
	record, err := s.Get(ctx, id, token)
	if err != nil {
		return task.Record{}, nil, err
	}
	if record.Job.Status() != job.StatusCompleted {
		return task.Record{}, nil, ErrNotReady
	}

	resultKey, ok := record.Job.ResultKey()
	if !ok {
		return task.Record{}, nil, ErrNotReady
	}
	reader, err := s.storage.Open(ctx, resultKey)
	if err != nil {
		return task.Record{}, nil, err
	}
	return record, reader, nil
}

func (s *Service) Delete(ctx context.Context, id string, token string) error {
	record, err := s.Get(ctx, id, token)
	if err != nil {
		return err
	}

	switch record.Job.Status() {
	case job.StatusProcessing, job.StatusCancelling:
		return ErrConflict
	}

	if err := s.repository.Delete(ctx, id); err != nil {
		return err
	}

	_ = s.storage.Delete(context.Background(), record.InputKey)
	if resultKey, ok := record.Job.ResultKey(); ok {
		_ = s.storage.Delete(context.Background(), resultKey)
	}
	return nil
}

func owns(record task.Record, token string) bool {
	hash := sha256.Sum256([]byte(token))
	return len(record.OwnerTokenHash) == len(hash) &&
		subtle.ConstantTimeCompare(record.OwnerTokenHash, hash[:]) == 1
}
