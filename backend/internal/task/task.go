package task

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/saral-gupta7/recode/backend/internal/job"
)

var (
	ErrNotFound      = errors.New("task not found")
	ErrInvalidInput  = errors.New("invalid task input")
	ErrInvalidOption = errors.New("invalid task option")
)

const (
	maxImageDimension = 16384
	maxClipDuration   = 6 * 60 * 60
)

type Options struct {
	Format          string  `json:"format,omitempty"`
	Width           int     `json:"width,omitempty"`
	Height          int     `json:"height,omitempty"`
	Quality         int     `json:"quality,omitempty"`
	StartSeconds    float64 `json:"start_seconds,omitempty"`
	DurationSeconds float64 `json:"duration_seconds,omitempty"`
}

type Record struct {
	Job              *job.Job
	InputKey         string
	OriginalFilename string
	MIMEType         string
	SizeBytes        int64
	Options          Options
	OwnerTokenHash   []byte
}

type Repository interface {
	Create(ctx context.Context, record Record) error
	FindByID(ctx context.Context, id string) (Record, error)
	FindExpirable(ctx context.Context, now time.Time, limit int) ([]Record, error)
	SaveJob(ctx context.Context, storedJob *job.Job) error
	Delete(ctx context.Context, id string) error
}

func NormalizeOptions(operation job.Operation, options Options) (Options, error) {
	options.Format = strings.ToLower(strings.TrimPrefix(strings.TrimSpace(options.Format), "."))

	switch operation {
	case job.OperationImageConvert:
		if options.Format == "" {
			options.Format = "png"
		}
		if !oneOf(options.Format, "jpg", "jpeg", "png", "webp") {
			return Options{}, fmt.Errorf("%w: unsupported image format %q", ErrInvalidOption, options.Format)
		}

	case job.OperationVideoConvert:
		if options.Format == "" {
			options.Format = "mp4"
		}
		if !oneOf(options.Format, "mp4", "webm", "mov") {
			return Options{}, fmt.Errorf("%w: unsupported video format %q", ErrInvalidOption, options.Format)
		}

	case job.OperationVideoExtractAudio:
		if options.Format == "" {
			options.Format = "mp3"
		}
		if !oneOf(options.Format, "mp3", "wav", "m4a") {
			return Options{}, fmt.Errorf("%w: unsupported audio format %q", ErrInvalidOption, options.Format)
		}

	case job.OperationImageResize:
		if options.Width < 0 ||
			options.Height < 0 ||
			options.Width > maxImageDimension ||
			options.Height > maxImageDimension ||
			(options.Width == 0 && options.Height == 0) {
			return Options{}, fmt.Errorf("%w: width or height is required", ErrInvalidOption)
		}

	case job.OperationImageCompress:
		if options.Quality == 0 {
			options.Quality = 80
		}
		if options.Quality < 1 || options.Quality > 100 {
			return Options{}, fmt.Errorf("%w: quality must be between 1 and 100", ErrInvalidOption)
		}

	case job.OperationVideoClip:
		if options.StartSeconds < 0 ||
			options.DurationSeconds <= 0 ||
			options.DurationSeconds > maxClipDuration {
			return Options{}, fmt.Errorf("%w: valid start and positive duration are required", ErrInvalidOption)
		}
	}

	return options, nil
}

func DecodeOptions(data []byte) (Options, error) {
	if len(data) == 0 {
		return Options{}, nil
	}

	var options Options
	decoder := json.NewDecoder(strings.NewReader(string(data)))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&options); err != nil {
		return Options{}, fmt.Errorf("%w: %v", ErrInvalidOption, err)
	}
	var trailing any
	if err := decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
		return Options{}, fmt.Errorf("%w: options must contain one JSON object", ErrInvalidOption)
	}
	return options, nil
}

func OutputExtension(operation job.Operation, options Options) string {
	switch operation {
	case job.OperationImageConvert, job.OperationVideoConvert, job.OperationVideoExtractAudio:
		return options.Format
	case job.OperationImageGrayscale:
		return "png"
	case job.OperationImageCompress, job.OperationImageResize:
		return "jpg"
	default:
		return "mp4"
	}
}

func NewRecord(
	storedJob *job.Job,
	inputKey string,
	originalFilename string,
	mimeType string,
	sizeBytes int64,
	options Options,
	ownerTokenHash []byte,
) (Record, error) {
	if storedJob == nil ||
		strings.TrimSpace(inputKey) == "" ||
		sizeBytes <= 0 ||
		len(ownerTokenHash) == 0 {
		return Record{}, ErrInvalidInput
	}

	return Record{
		Job:              storedJob,
		InputKey:         inputKey,
		OriginalFilename: originalFilename,
		MIMEType:         mimeType,
		SizeBytes:        sizeBytes,
		Options:          options,
		OwnerTokenHash:   append([]byte(nil), ownerTokenHash...),
	}, nil
}

func ResultExpiration(now time.Time, retention time.Duration) time.Time {
	return now.UTC().Add(retention)
}

func oneOf(value string, allowed ...string) bool {
	for _, candidate := range allowed {
		if value == candidate {
			return true
		}
	}
	return false
}
