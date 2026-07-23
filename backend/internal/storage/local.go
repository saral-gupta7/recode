package storage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Local struct {
	root string
}

var _ Store = (*Local)(nil)

func NewLocal(root string) (*Local, error) {
	if strings.TrimSpace(root) == "" {
		return nil, errors.New("storage root must not be empty")
	}

	absoluteRoot, err := filepath.Abs(root)
	if err != nil {
		return nil, fmt.Errorf("resolve storage root: %w", err)
	}

	if err := os.MkdirAll(absoluteRoot, 0o750); err != nil {
		return nil, fmt.Errorf("create storage root: %w", err)
	}

	return &Local{
		root: absoluteRoot,
	}, nil
}

func (l *Local) Put(
	ctx context.Context,
	key string,
	source io.Reader,
	maxBytes int64,
) (int64, error) {
	if maxBytes <= 0 {
		return 0, errors.New("storage size limit must be greater than zero")
	}

	target, err := l.resolve(key)
	if err != nil {
		return 0, err
	}

	if err := os.MkdirAll(filepath.Dir(target), 0o750); err != nil {
		return 0, fmt.Errorf("create object directory: %w", err)
	}

	temporary, err := os.CreateTemp(
		filepath.Dir(target),
		".upload-*",
	)
	if err != nil {
		return 0, fmt.Errorf("create temporary object: %w", err)
	}

	temporaryName := temporary.Name()
	published := false

	defer func() {
		_ = temporary.Close()

		if !published {
			_ = os.Remove(temporaryName)
		}
	}()

	written, err := io.Copy(
		temporary,
		io.LimitReader(
			contextReader{
				ctx:    ctx,
				reader: source,
			},
			maxBytes+1,
		),
	)
	if err != nil {
		return 0, fmt.Errorf("write temporary object: %w", err)
	}

	if written > maxBytes {
		return 0, ErrTooLarge
	}

	if err := temporary.Sync(); err != nil {
		return 0, fmt.Errorf("sync temporary object: %w", err)
	}

	if err := temporary.Close(); err != nil {
		return 0, fmt.Errorf("close temporary object: %w", err)
	}

	if err := os.Rename(temporaryName, target); err != nil {
		return 0, fmt.Errorf("publish object: %w", err)
	}

	published = true

	return written, nil
}

func (l *Local) resolve(key string) (string, error) {
	if strings.TrimSpace(key) == "" || filepath.IsAbs(key) {
		return "", ErrInvalidKey
	}

	cleanKey := filepath.Clean(key)

	if cleanKey == "." ||
		cleanKey == ".." ||
		strings.HasPrefix(
			cleanKey,
			".."+string(filepath.Separator),
		) {
		return "", ErrInvalidKey
	}

	path := filepath.Join(l.root, cleanKey)

	relative, err := filepath.Rel(l.root, path)
	if err != nil ||
		relative == ".." ||
		strings.HasPrefix(
			relative,
			".."+string(filepath.Separator),
		) {
		return "", ErrInvalidKey
	}

	return path, nil
}

type contextReader struct {
	ctx    context.Context
	reader io.Reader
}

func (r contextReader) Read(buffer []byte) (int, error) {
	select {
	case <-r.ctx.Done():
		return 0, r.ctx.Err()

	default:
		return r.reader.Read(buffer)
	}
}

func (l *Local) Open(ctx context.Context, key string) (io.ReadCloser, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	path, err := l.resolve(key)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("open object: %w", err)
	}
	return file, nil
}

func (l *Local) Delete(ctx context.Context, key string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	path, err := l.resolve(key)
	if err != nil {
		return err
	}

	if err := os.Remove(path); errors.Is(err, os.ErrNotExist) {
		return ErrNotFound
	} else if err != nil {
		return fmt.Errorf("delete object: %w", err)
	}
	return nil
}
