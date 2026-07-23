package queue

import (
	"context"
	"errors"
	"time"
)

var ErrEmpty = errors.New("queue is empty")

type Queue interface {
	Enqueue(ctx context.Context, jobID string) error
	Dequeue(ctx context.Context, timeout time.Duration) (string, error)
}
