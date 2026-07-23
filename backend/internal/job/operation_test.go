package job

import (
	"errors"
	"testing"
)

func TestOperationValid(t *testing.T) {
	tests := []struct {
		name      string
		operation Operation
		want      bool
	}{
		{name: "video grayscale", operation: OperationVideoGrayscale, want: true},
		{name: "video extract audio", operation: OperationVideoExtractAudio, want: true},
		{name: "video remove audio", operation: OperationVideoRemoveAudio, want: true},
		{name: "video convert", operation: OperationVideoConvert, want: true},
		{name: "video clip", operation: OperationVideoClip, want: true},
		{name: "image grayscale", operation: OperationImageGrayscale, want: true},
		{name: "image convert", operation: OperationImageConvert, want: true},
		{name: "image compress", operation: OperationImageCompress, want: true},
		{name: "image resize", operation: OperationImageResize, want: true},
		{name: "empty", operation: Operation(""), want: false},
		{name: "unknown", operation: Operation("unknown_operation"), want: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.operation.Valid(); got != test.want {
				t.Fatalf("Operation(%q).Valid() = %t, want %t", test.operation, got, test.want)
			}
		})
	}
}

func TestParseOperation(t *testing.T) {
	operation, err := ParseOperation("video_clip")
	if err != nil {
		t.Fatalf("ParseOperation() unexpected error: %v", err)
	}

	if operation != OperationVideoClip {
		t.Fatalf("ParseOperation() = %q, want %q", operation, OperationVideoClip)
	}
}

func TestParseOperationRejectsInvalidValue(t *testing.T) {
	operation, err := ParseOperation("turn_video_blue")

	if !errors.Is(err, ErrInvalidOperation) {
		t.Fatalf("ParseOperation() error = %v, want ErrInvalidOperation", err)
	}

	if operation != "" {
		t.Fatalf("ParseOperation() = %q, want zero Operation", operation)
	}
}
