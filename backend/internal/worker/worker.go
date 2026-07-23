package worker

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/saral-gupta7/recode/backend/internal/job"
	"github.com/saral-gupta7/recode/backend/internal/media"
	"github.com/saral-gupta7/recode/backend/internal/queue"
	"github.com/saral-gupta7/recode/backend/internal/storage"
	"github.com/saral-gupta7/recode/backend/internal/task"
)

type Worker struct {
	logger         *slog.Logger
	repository     task.Repository
	queue          queue.Queue
	storage        storage.Store
	prober         media.Prober
	processor      media.Processor
	retention      time.Duration
	maxOutputBytes int64
}

func New(
	logger *slog.Logger,
	repository task.Repository,
	jobQueue queue.Queue,
	store storage.Store,
	prober media.Prober,
	processor media.Processor,
	retention time.Duration,
	maxOutputBytes int64,
) *Worker {
	return &Worker{
		logger:         logger,
		repository:     repository,
		queue:          jobQueue,
		storage:        store,
		prober:         prober,
		processor:      processor,
		retention:      retention,
		maxOutputBytes: maxOutputBytes,
	}
}

func (w *Worker) Run(ctx context.Context) error {
	w.logger.Info("worker started")
	nextCleanup := time.Now()
	for {
		if !time.Now().Before(nextCleanup) {
			if err := w.cleanup(ctx); err != nil && ctx.Err() == nil {
				w.logger.Error("cleanup failed", slog.Any("error", err))
			}
			nextCleanup = time.Now().Add(time.Minute)
		}

		jobID, err := w.queue.Dequeue(ctx, 5*time.Second)
		if errors.Is(err, queue.ErrEmpty) {
			continue
		}
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			w.logger.Error("dequeue failed", slog.Any("error", err))
			continue
		}

		if err := w.process(ctx, jobID); err != nil {
			w.logger.Error(
				"job processing failed",
				slog.String("job_id", jobID),
				slog.Any("error", err),
			)
		}
	}
}

func (w *Worker) cleanup(ctx context.Context) error {
	records, err := w.repository.FindExpirable(ctx, time.Now().UTC(), 100)
	if err != nil {
		return err
	}

	for _, record := range records {
		resultKey, hasResult := record.Job.ResultKey()
		now := time.Now().UTC()
		if err := record.Job.TransitionTo(job.StatusExpired, now); err != nil {
			w.logger.Error(
				"expire transition failed",
				slog.String("job_id", record.Job.ID()),
				slog.Any("error", err),
			)
			continue
		}
		if err := w.repository.SaveJob(ctx, record.Job); err != nil {
			w.logger.Error(
				"persist expiration failed",
				slog.String("job_id", record.Job.ID()),
				slog.Any("error", err),
			)
			continue
		}

		_ = w.storage.Delete(context.Background(), record.InputKey)
		if hasResult {
			_ = w.storage.Delete(context.Background(), resultKey)
		}
	}
	return nil
}

func (w *Worker) process(ctx context.Context, jobID string) error {
	record, err := w.repository.FindByID(ctx, jobID)
	if errors.Is(err, task.ErrNotFound) {
		return nil
	}
	if err != nil {
		return err
	}
	if record.Job.Status() != job.StatusQueued {
		return nil
	}

	if err := record.Job.TransitionTo(job.StatusProcessing, time.Now().UTC()); err != nil {
		return err
	}
	if err := w.repository.SaveJob(ctx, record.Job); err != nil {
		return err
	}

	workDirectory, err := os.MkdirTemp("", "recode-"+jobID+"-*")
	if err != nil {
		return w.fail(ctx, jobID, "work_directory_failed", err)
	}
	defer os.RemoveAll(workDirectory)

	inputPath := filepath.Join(workDirectory, "input")
	if err := w.materializeInput(ctx, record.InputKey, inputPath); err != nil {
		return w.fail(ctx, jobID, "input_read_failed", err)
	}

	info, err := w.prober.Probe(ctx, inputPath)
	if err != nil {
		return w.fail(ctx, jobID, "invalid_media", err)
	}
	if err := media.ValidateForOperation(record.Job.Operation(), info); err != nil {
		return w.fail(ctx, jobID, "unsupported_media", err)
	}

	extension := task.OutputExtension(record.Job.Operation(), record.Options)
	outputPath := filepath.Join(workDirectory, "output."+extension)
	if err := w.processor.Process(
		ctx,
		record.Job.Operation(),
		record.Options,
		inputPath,
		outputPath,
	); err != nil {
		return w.fail(ctx, jobID, "processing_failed", err)
	}

	latest, err := w.repository.FindByID(ctx, jobID)
	if err != nil {
		return err
	}
	if latest.Job.Status() == job.StatusCancelling {
		return w.finishCancelled(ctx, latest)
	}
	if latest.Job.Status() != job.StatusProcessing {
		return nil
	}

	output, err := os.Open(outputPath)
	if err != nil {
		return w.fail(ctx, jobID, "output_read_failed", err)
	}
	defer output.Close()

	resultKey := "results/" + jobID + "/output." + extension
	if _, err := w.storage.Put(ctx, resultKey, output, w.maxOutputBytes); err != nil {
		return w.fail(ctx, jobID, "output_store_failed", err)
	}

	now := time.Now().UTC()
	if err := latest.Job.TransitionTo(job.StatusCompleted, now); err != nil {
		return err
	}
	if err := latest.Job.RecordResult(resultKey); err != nil {
		return err
	}
	if err := latest.Job.ScheduleExpiration(task.ResultExpiration(now, w.retention)); err != nil {
		return err
	}
	if err := w.repository.SaveJob(ctx, latest.Job); err != nil {
		_ = w.storage.Delete(context.Background(), resultKey)
		return err
	}

	_ = w.storage.Delete(context.Background(), latest.InputKey)
	w.logger.Info(
		"job completed",
		slog.String("job_id", jobID),
		slog.String("operation", string(latest.Job.Operation())),
	)
	return nil
}

func (w *Worker) materializeInput(ctx context.Context, key string, target string) error {
	source, err := w.storage.Open(ctx, key)
	if err != nil {
		return err
	}
	defer source.Close()

	file, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0o600)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, source)
	return err
}

func (w *Worker) fail(ctx context.Context, jobID string, code string, cause error) error {
	record, err := w.repository.FindByID(ctx, jobID)
	if err != nil {
		return errors.Join(cause, err)
	}
	if record.Job.Status() == job.StatusCancelling {
		return errors.Join(cause, w.finishCancelled(ctx, record))
	}
	if record.Job.Status() != job.StatusProcessing {
		return cause
	}

	now := time.Now().UTC()
	if err := record.Job.TransitionTo(job.StatusFailed, now); err != nil {
		return errors.Join(cause, err)
	}
	if err := record.Job.RecordFailure(code); err != nil {
		return errors.Join(cause, err)
	}
	if err := record.Job.ScheduleExpiration(task.ResultExpiration(now, w.retention)); err != nil {
		return errors.Join(cause, err)
	}
	if err := w.repository.SaveJob(ctx, record.Job); err != nil {
		return errors.Join(cause, err)
	}
	return fmt.Errorf("%s: %w", code, cause)
}

func (w *Worker) finishCancelled(ctx context.Context, record task.Record) error {
	now := time.Now().UTC()
	if err := record.Job.TransitionTo(job.StatusCancelled, now); err != nil {
		return err
	}
	if err := record.Job.ScheduleExpiration(task.ResultExpiration(now, w.retention)); err != nil {
		return err
	}
	if err := w.repository.SaveJob(ctx, record.Job); err != nil {
		return err
	}
	_ = w.storage.Delete(context.Background(), record.InputKey)
	return nil
}
