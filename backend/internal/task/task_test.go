package task

import (
	"errors"
	"testing"

	"github.com/saral-gupta7/recode/backend/internal/job"
)

func TestNormalizeOptionsDefaults(t *testing.T) {
	tests := []struct {
		operation   job.Operation
		wantFormat  string
		wantQuality int
	}{
		{operation: job.OperationImageConvert, wantFormat: "png"},
		{operation: job.OperationVideoConvert, wantFormat: "mp4"},
		{operation: job.OperationVideoExtractAudio, wantFormat: "mp3"},
		{operation: job.OperationImageCompress, wantQuality: 80},
	}

	for _, test := range tests {
		t.Run(string(test.operation), func(t *testing.T) {
			got, err := NormalizeOptions(test.operation, Options{})
			if err != nil {
				t.Fatalf("NormalizeOptions() error = %v", err)
			}
			if got.Format != test.wantFormat || got.Quality != test.wantQuality {
				t.Errorf("NormalizeOptions() = %#v", got)
			}
		})
	}
}

func TestNormalizeOptionsRejectsInvalidValues(t *testing.T) {
	tests := []struct {
		name      string
		operation job.Operation
		options   Options
	}{
		{name: "image format", operation: job.OperationImageConvert, options: Options{Format: "exe"}},
		{name: "video format", operation: job.OperationVideoConvert, options: Options{Format: "gif"}},
		{name: "audio format", operation: job.OperationVideoExtractAudio, options: Options{Format: "mp4"}},
		{name: "resize dimensions", operation: job.OperationImageResize},
		{name: "negative resize dimension", operation: job.OperationImageResize, options: Options{Width: -1, Height: 100}},
		{name: "oversized resize dimension", operation: job.OperationImageResize, options: Options{Width: maxImageDimension + 1}},
		{name: "compression quality", operation: job.OperationImageCompress, options: Options{Quality: 101}},
		{name: "clip duration", operation: job.OperationVideoClip, options: Options{DurationSeconds: 0}},
		{name: "excessive clip duration", operation: job.OperationVideoClip, options: Options{DurationSeconds: maxClipDuration + 1}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := NormalizeOptions(test.operation, test.options)
			if !errors.Is(err, ErrInvalidOption) {
				t.Fatalf("NormalizeOptions() error = %v, want ErrInvalidOption", err)
			}
		})
	}
}

func TestDecodeOptionsRejectsUnknownField(t *testing.T) {
	_, err := DecodeOptions([]byte(`{"surprise":true}`))
	if !errors.Is(err, ErrInvalidOption) {
		t.Fatalf("DecodeOptions() error = %v, want ErrInvalidOption", err)
	}
}

func TestDecodeOptionsRejectsTrailingJSON(t *testing.T) {
	_, err := DecodeOptions([]byte(`{"width":100} {"height":100}`))
	if !errors.Is(err, ErrInvalidOption) {
		t.Fatalf("DecodeOptions() error = %v, want ErrInvalidOption", err)
	}
}
