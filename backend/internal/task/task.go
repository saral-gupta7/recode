package task

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"

	"github.com/saral-gupta7/recode/backend/internal/job"
)

var (
	ErrNotFound      = errors.New("task not found")
	ErrInvalidInput  = errors.New("invalid task input")
	ErrInvalidOption = errors.New("invalid task option")
)

const maxImageDimension = 16384

var colourPattern = regexp.MustCompile(`^(#[0-9a-fA-F]{6}|[a-zA-Z]{3,20})$`)

type Options struct {
	Format        string        `json:"format,omitempty"`
	Width         int           `json:"width,omitempty"`
	Height        int           `json:"height,omitempty"`
	Quality       int           `json:"quality,omitempty"`
	X             int           `json:"x,omitempty"`
	Y             int           `json:"y,omitempty"`
	Angle         int           `json:"angle,omitempty"`
	FlipDirection FlipDirection `json:"flip_direction,omitempty"`
	Preset        string        `json:"preset,omitempty"`
	Brightness    int           `json:"brightness,omitempty"`
	Contrast      int           `json:"contrast,omitempty"`
	Saturation    int           `json:"saturation,omitempty"`
	Strength      float64       `json:"strength,omitempty"`
	BlockSize     int           `json:"block_size,omitempty"`
	PaddingTop    int           `json:"padding_top,omitempty"`
	PaddingRight  int           `json:"padding_right,omitempty"`
	PaddingBottom int           `json:"padding_bottom,omitempty"`
	PaddingLeft   int           `json:"padding_left,omitempty"`
	Background    string        `json:"background,omitempty"`
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
	options.FlipDirection = FlipDirection(strings.ToLower(strings.TrimSpace(string(options.FlipDirection))))
	options.Preset = strings.ToLower(strings.TrimSpace(options.Preset))
	options.Background = strings.ToLower(strings.TrimSpace(options.Background))

	switch operation {
	case job.OperationImageConvert:
		if options.Format == "" {
			options.Format = "png"
		}
		if !oneOf(options.Format, "jpg", "jpeg", "png", "webp") {
			return Options{}, fmt.Errorf("%w: unsupported image format %q", ErrInvalidOption, options.Format)
		}
	case job.OperationImageResize:
		if !validDimensions(options.Width, options.Height, false) {
			return Options{}, fmt.Errorf("%w: width or height is required", ErrInvalidOption)
		}
	case job.OperationImageCrop:
		if options.X < 0 || options.Y < 0 || !validDimensions(options.Width, options.Height, true) {
			return Options{}, fmt.Errorf("%w: crop coordinates and dimensions are invalid", ErrInvalidOption)
		}
	case job.OperationImageCompress:
		if options.Quality == 0 {
			options.Quality = 80
		}
		if options.Quality < 1 || options.Quality > 100 {
			return Options{}, fmt.Errorf("%w: quality must be between 1 and 100", ErrInvalidOption)
		}
	case job.OperationImageRotate:
		if !oneOfInt(options.Angle, 90, 180, 270) {
			return Options{}, fmt.Errorf("%w: angle must be 90, 180, or 270", ErrInvalidOption)
		}
	case job.OperationImageFlip:
		if !options.FlipDirection.Valid() {
			return Options{}, fmt.Errorf("%w: flip direction must be horizontal or vertical", ErrInvalidOption)
		}
	case job.OperationImageThumbnail:
		if options.Preset == "" {
			options.Preset = "square"
		}
		if !oneOf(options.Preset, "square", "preview", "social") {
			return Options{}, fmt.Errorf("%w: unsupported thumbnail preset", ErrInvalidOption)
		}
	case job.OperationImageAdjust:
		if options.Brightness < -100 || options.Brightness > 100 || options.Contrast < -100 || options.Contrast > 100 || options.Saturation < 0 || options.Saturation > 200 {
			return Options{}, fmt.Errorf("%w: adjustment is outside its allowed range", ErrInvalidOption)
		}
	case job.OperationImageBlur:
		if options.Strength == 0 {
			options.Strength = 2
		}
		if options.Strength < 0.1 || options.Strength > 20 {
			return Options{}, fmt.Errorf("%w: strength must be between 0.1 and 20", ErrInvalidOption)
		}
	case job.OperationImageSharpen:
		if options.Strength == 0 {
			options.Strength = 2
		}
		if options.Strength < 0.1 || options.Strength > 5 {
			return Options{}, fmt.Errorf("%w: strength must be between 0.1 and 5", ErrInvalidOption)
		}
	case job.OperationImagePixelate:
		if options.BlockSize == 0 {
			options.BlockSize = 12
		}
		if options.BlockSize < 2 || options.BlockSize > 100 {
			return Options{}, fmt.Errorf("%w: block size must be between 2 and 100", ErrInvalidOption)
		}
	case job.OperationImagePadding:
		if options.Background == "" {
			options.Background = "white"
		}
		if options.PaddingTop < 0 || options.PaddingRight < 0 || options.PaddingBottom < 0 || options.PaddingLeft < 0 ||
			options.PaddingTop > maxImageDimension || options.PaddingRight > maxImageDimension || options.PaddingBottom > maxImageDimension || options.PaddingLeft > maxImageDimension ||
			(options.PaddingTop+options.PaddingRight+options.PaddingBottom+options.PaddingLeft == 0) || !colourPattern.MatchString(options.Background) {
			return Options{}, fmt.Errorf("%w: padding or background colour is invalid", ErrInvalidOption)
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
	if operation == job.OperationImageConvert {
		return options.Format
	}
	switch operation {
	case job.OperationImageCompress, job.OperationImageResize, job.OperationImageThumbnail, job.OperationImageStripMetadata:
		return "jpg"
	default:
		return "png"
	}
}

func NewRecord(storedJob *job.Job, inputKey, originalFilename, mimeType string, sizeBytes int64, options Options, ownerTokenHash []byte) (Record, error) {
	if storedJob == nil || strings.TrimSpace(inputKey) == "" || sizeBytes <= 0 || len(ownerTokenHash) == 0 {
		return Record{}, ErrInvalidInput
	}
	return Record{Job: storedJob, InputKey: inputKey, OriginalFilename: originalFilename, MIMEType: mimeType, SizeBytes: sizeBytes, Options: options, OwnerTokenHash: append([]byte(nil), ownerTokenHash...)}, nil
}

func ResultExpiration(now time.Time, retention time.Duration) time.Time {
	return now.UTC().Add(retention)
}

func validDimensions(width, height int, both bool) bool {
	if width < 0 || height < 0 || width > maxImageDimension || height > maxImageDimension {
		return false
	}
	if both {
		return width > 0 && height > 0
	}
	return width > 0 || height > 0
}
func oneOf(value string, allowed ...string) bool {
	for _, candidate := range allowed {
		if value == candidate {
			return true
		}
	}
	return false
}
func oneOfInt(value int, allowed ...int) bool {
	for _, candidate := range allowed {
		if value == candidate {
			return true
		}
	}
	return false
}
