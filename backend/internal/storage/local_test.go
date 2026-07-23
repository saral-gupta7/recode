package storage

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLocalLifecycle(t *testing.T) {
	store, err := NewLocal(t.TempDir())
	if err != nil {
		t.Fatalf("NewLocal() error = %v", err)
	}

	written, err := store.Put(
		context.Background(),
		"uploads/job/input",
		strings.NewReader("media"),
		10,
	)
	if err != nil {
		t.Fatalf("Put() error = %v", err)
	}
	if written != 5 {
		t.Errorf("Put() bytes = %d, want 5", written)
	}

	reader, err := store.Open(context.Background(), "uploads/job/input")
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	content, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("ReadAll() error = %v", err)
	}
	_ = reader.Close()
	if string(content) != "media" {
		t.Errorf("content = %q, want media", content)
	}

	if err := store.Delete(context.Background(), "uploads/job/input"); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
	if _, err := store.Open(context.Background(), "uploads/job/input"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("Open() after delete error = %v, want ErrNotFound", err)
	}
}

func TestLocalRejectsOversizedObjectWithoutPublishing(t *testing.T) {
	root := t.TempDir()
	store, err := NewLocal(root)
	if err != nil {
		t.Fatalf("NewLocal() error = %v", err)
	}

	_, err = store.Put(
		context.Background(),
		"uploads/job/input",
		strings.NewReader("too large"),
		3,
	)
	if !errors.Is(err, ErrTooLarge) {
		t.Fatalf("Put() error = %v, want ErrTooLarge", err)
	}
	if _, err := os.Stat(filepath.Join(root, "uploads/job/input")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("published file error = %v, want os.ErrNotExist", err)
	}
}

func TestLocalRejectsUnsafeKeys(t *testing.T) {
	store, err := NewLocal(t.TempDir())
	if err != nil {
		t.Fatalf("NewLocal() error = %v", err)
	}

	for _, key := range []string{"", " ", ".", "..", "../outside", "/absolute"} {
		t.Run(key, func(t *testing.T) {
			_, err := store.Put(context.Background(), key, strings.NewReader("x"), 1)
			if !errors.Is(err, ErrInvalidKey) {
				t.Fatalf("Put(%q) error = %v, want ErrInvalidKey", key, err)
			}
		})
	}
}

func TestLocalHonorsCancelledContext(t *testing.T) {
	store, err := NewLocal(t.TempDir())
	if err != nil {
		t.Fatalf("NewLocal() error = %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if _, err := store.Put(ctx, "uploads/job/input", strings.NewReader("x"), 1); !errors.Is(err, context.Canceled) {
		t.Fatalf("Put() error = %v, want context.Canceled", err)
	}
}
