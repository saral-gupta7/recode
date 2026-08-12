package task

import (
	"errors"
	"testing"

	"github.com/saral-gupta7/recode/backend/internal/job"
)

func TestNormalizeOptionsDefaults(t *testing.T) {
	tests := []struct {
		operation job.Operation
		input     Options
		check     func(Options) bool
	}{
		{job.OperationImageConvert, Options{}, func(o Options) bool { return o.Format == "png" }},
		{job.OperationImageCompress, Options{}, func(o Options) bool { return o.Quality == 80 }},
		{job.OperationImageThumbnail, Options{}, func(o Options) bool { return o.Preset == "square" }},
		{job.OperationImageBlur, Options{}, func(o Options) bool { return o.Strength == 2 }},
		{job.OperationImagePixelate, Options{}, func(o Options) bool { return o.BlockSize == 12 }},
	}
	for _, test := range tests {
		t.Run(string(test.operation), func(t *testing.T) {
			got, err := NormalizeOptions(test.operation, test.input)
			if err != nil || !test.check(got) {
				t.Fatalf("NormalizeOptions() = %+v, %v", got, err)
			}
		})
	}
}

func TestNormalizeOptionsRejectsInvalidValues(t *testing.T) {
	tests := []struct {
		operation job.Operation
		options   Options
	}{
		{job.OperationImageConvert, Options{Format: "gif"}},
		{job.OperationImageResize, Options{}},
		{job.OperationImageCrop, Options{Width: 10}},
		{job.OperationImageCompress, Options{Quality: 101}},
		{job.OperationImageRotate, Options{Angle: 45}},
		{job.OperationImageFlip, Options{}},
		{job.OperationImageThumbnail, Options{Preset: "giant"}},
		{job.OperationImageAdjust, Options{Brightness: 101}},
		{job.OperationImageBlur, Options{Strength: 21}},
		{job.OperationImagePixelate, Options{BlockSize: 1}},
		{job.OperationImagePadding, Options{PaddingTop: 1, Background: "url(evil)"}},
	}
	for _, test := range tests {
		if _, err := NormalizeOptions(test.operation, test.options); !errors.Is(err, ErrInvalidOption) {
			t.Fatalf("%q error = %v", test.operation, err)
		}
	}
}

func TestDecodeOptionsRejectsUnknownAndTrailingJSON(t *testing.T) {
	for _, data := range [][]byte{[]byte(`{"unknown":1}`), []byte(`{} {}`)} {
		if _, err := DecodeOptions(data); !errors.Is(err, ErrInvalidOption) {
			t.Fatalf("DecodeOptions() error = %v", err)
		}
	}
}

func TestOutputExtensionsAreImages(t *testing.T) {
	if got := OutputExtension(job.OperationImageConvert, Options{Format: "webp"}); got != "webp" {
		t.Fatal(got)
	}
	if got := OutputExtension(job.OperationImageCrop, Options{}); got != "png" {
		t.Fatal(got)
	}
	if got := OutputExtension(job.OperationImageThumbnail, Options{}); got != "jpg" {
		t.Fatal(got)
	}
}
