package job

import (
	"context"
	"errors"
)

var ErrNotFound = errors.New("job not found")

type Repository interface {
	Create(ctx context.Context, job *Job) error
	FindByID(ctx context.Context, id string) (*Job, error)
	Save(ctx context.Context, job *Job) error
}
