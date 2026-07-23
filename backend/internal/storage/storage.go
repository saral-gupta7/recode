package storage

import (
	"context"
	"errors"
	"io"
)

var (
	ErrInvalidKey = errors.New("invalid storage key")
	ErrTooLarge   = errors.New("stored object exceeds size limit")
	ErrNotFound   = errors.New("stored object not found")
)

type Store interface {
	Put(
		ctx context.Context,
		key string,
		source io.Reader,
		maxBytes int64,
	) (int64, error)

	Open(
		ctx context.Context,
		key string,
	) (io.ReadCloser, error)

	Delete(
		ctx context.Context,
		key string,
	) error
}
