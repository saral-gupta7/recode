package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/saral-gupta7/recode/backend/internal/job"
	"github.com/saral-gupta7/recode/backend/internal/task"
)

var ErrNilTask = errors.New("task and task job must not be nil")

type TaskRepository struct {
	pool *pgxpool.Pool
}

var _ task.Repository = (*TaskRepository)(nil)

func NewTaskRepository(pool *pgxpool.Pool) *TaskRepository {
	return &TaskRepository{pool: pool}
}

func (r *TaskRepository) Create(ctx context.Context, record task.Record) error {
	if record.Job == nil {
		return ErrNilTask
	}

	snapshot := record.Job.Snapshot()
	options, err := json.Marshal(record.Options)
	if err != nil {
		return fmt.Errorf("encode task options: %w", err)
	}

	const query = `
		INSERT INTO jobs (
			id, operation, status, progress, attempt,
			expires_at, result_key, failure_code,
			created_at, updated_at,
			input_key, original_filename, mime_type,
			size_bytes, options, owner_token_hash
		)
		VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8,
			$9, $10,
			$11, $12, $13,
			$14, $15, $16
		)
	`

	if _, err := r.pool.Exec(
		ctx,
		query,
		snapshot.ID,
		snapshot.Operation,
		snapshot.Status,
		snapshot.Progress,
		snapshot.Attempt,
		nullableTime(snapshot.ExpiresAt),
		nullableString(snapshot.ResultKey),
		nullableString(snapshot.FailureCode),
		snapshot.CreatedAt,
		snapshot.UpdatedAt,
		record.InputKey,
		record.OriginalFilename,
		record.MIMEType,
		record.SizeBytes,
		options,
		record.OwnerTokenHash,
	); err != nil {
		return fmt.Errorf("insert task: %w", err)
	}
	return nil
}

func (r *TaskRepository) FindByID(ctx context.Context, id string) (task.Record, error) {
	const query = `
		SELECT
			id, operation, status, progress, attempt,
			expires_at, result_key, failure_code,
			created_at, updated_at,
			input_key, original_filename, mime_type,
			size_bytes, options, owner_token_hash
		FROM jobs
		WHERE id = $1
	`

	var snapshot job.Snapshot
	var operation string
	var status string
	var expiresAt pgtype.Timestamptz
	var resultKey pgtype.Text
	var failureCode pgtype.Text
	var inputKey pgtype.Text
	var originalFilename pgtype.Text
	var mimeType pgtype.Text
	var sizeBytes pgtype.Int8
	var optionsData []byte
	var ownerTokenHash []byte

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&snapshot.ID,
		&operation,
		&status,
		&snapshot.Progress,
		&snapshot.Attempt,
		&expiresAt,
		&resultKey,
		&failureCode,
		&snapshot.CreatedAt,
		&snapshot.UpdatedAt,
		&inputKey,
		&originalFilename,
		&mimeType,
		&sizeBytes,
		&optionsData,
		&ownerTokenHash,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return task.Record{}, task.ErrNotFound
	}
	if err != nil {
		return task.Record{}, fmt.Errorf("select task %q: %w", id, err)
	}

	snapshot.Operation = job.Operation(operation)
	snapshot.Status = job.Status(status)
	if expiresAt.Valid {
		snapshot.ExpiresAt = expiresAt.Time
	}
	if resultKey.Valid {
		snapshot.ResultKey = resultKey.String
	}
	if failureCode.Valid {
		snapshot.FailureCode = failureCode.String
	}

	restoredJob, err := job.Restore(snapshot)
	if err != nil {
		return task.Record{}, fmt.Errorf("restore task job %q: %w", id, err)
	}

	options, err := task.DecodeOptions(optionsData)
	if err != nil {
		return task.Record{}, fmt.Errorf("decode stored task options: %w", err)
	}

	return task.Record{
		Job:              restoredJob,
		InputKey:         inputKey.String,
		OriginalFilename: originalFilename.String,
		MIMEType:         mimeType.String,
		SizeBytes:        sizeBytes.Int64,
		Options:          options,
		OwnerTokenHash:   append([]byte(nil), ownerTokenHash...),
	}, nil
}

func (r *TaskRepository) SaveJob(ctx context.Context, storedJob *job.Job) error {
	return NewJobRepository(r.pool).Save(ctx, storedJob)
}

func (r *TaskRepository) FindExpirable(
	ctx context.Context,
	now time.Time,
	limit int,
) ([]task.Record, error) {
	if limit <= 0 {
		return nil, nil
	}

	const query = `
		SELECT id
		FROM jobs
		WHERE status IN ('completed', 'failed', 'cancelled')
			AND expires_at <= $1
		ORDER BY expires_at
		LIMIT $2
	`

	rows, err := r.pool.Query(ctx, query, now, limit)
	if err != nil {
		return nil, fmt.Errorf("select expirable tasks: %w", err)
	}
	defer rows.Close()

	var records []task.Record
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("scan expirable task ID: %w", err)
		}
		record, err := r.FindByID(ctx, id)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate expirable tasks: %w", err)
	}
	return records, nil
}

func (r *TaskRepository) Delete(ctx context.Context, id string) error {
	result, err := r.pool.Exec(ctx, "DELETE FROM jobs WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("delete task %q: %w", id, err)
	}
	if result.RowsAffected() == 0 {
		return task.ErrNotFound
	}
	return nil
}
