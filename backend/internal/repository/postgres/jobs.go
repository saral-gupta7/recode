package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/saral-gupta7/recode/backend/internal/job"
)

var ErrNilJob = errors.New("job must not be nil")

type JobRepository struct {
	pool *pgxpool.Pool
}

var _ job.Repository = (*JobRepository)(nil)

func NewJobRepository(pool *pgxpool.Pool) *JobRepository {
	return &JobRepository{pool: pool}
}

func (r *JobRepository) Create(
	ctx context.Context,
	jobToCreate *job.Job,
) error {
	if jobToCreate == nil {
		return ErrNilJob
	}

	snapshot := jobToCreate.Snapshot()

	const query = `
		INSERT INTO jobs (
			id,
			operation,
			status,
			progress,
			attempt,
			expires_at,
			result_key,
			failure_code,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := r.pool.Exec(
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
	)

	if err != nil {
		return fmt.Errorf("insert job: %w", err)
	}

	return nil
}

func nullableTime(value time.Time) any {
	if value.IsZero() {
		return nil
	}

	return value
}

func nullableString(value string) any {
	if value == "" {
		return nil
	}

	return value
}

func (r *JobRepository) FindByID(
	ctx context.Context,
	id string,
) (*job.Job, error) {
	const query = `
		SELECT
			id,
			operation,
			status,
			progress,
			attempt,
			expires_at,
			result_key,
			failure_code,
			created_at,
			updated_at
		FROM jobs
		WHERE id = $1
	`

	var snapshot job.Snapshot
	var operation string
	var status string
	var expiresAt pgtype.Timestamptz
	var resultKey pgtype.Text
	var failureCode pgtype.Text

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
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, job.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("select job %q: %w", id, err)
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
		return nil, fmt.Errorf("restore job %q: %w", id, err)
	}

	return restoredJob, nil
}

func (r *JobRepository) Save(
	ctx context.Context,
	jobToSave *job.Job,
) error {
	if jobToSave == nil {
		return ErrNilJob
	}

	snapshot := jobToSave.Snapshot()

	const query = `
		UPDATE jobs
		SET
			status = $2,
			progress = $3,
			attempt = $4,
			expires_at = $5,
			result_key = $6,
			failure_code = $7,
			updated_at = $8
		WHERE id = $1
	`

	result, err := r.pool.Exec(
		ctx,
		query,
		snapshot.ID,
		snapshot.Status,
		snapshot.Progress,
		snapshot.Attempt,
		nullableTime(snapshot.ExpiresAt),
		nullableString(snapshot.ResultKey),
		nullableString(snapshot.FailureCode),
		snapshot.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("update job %q: %w", snapshot.ID, err)
	}

	if result.RowsAffected() == 0 {
		return job.ErrNotFound
	}

	return nil
}
