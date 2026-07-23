package job

import (
	"errors"
	"testing"
)

func TestOperationMediaKind(t *testing.T) {
	tests := []struct {
		operation Operation
		want      MediaKind
	}{
		{operation: OperationVideoGrayscale, want: MediaKindVideo},
		{operation: OperationVideoExtractAudio, want: MediaKindVideo},
		{operation: OperationVideoRemoveAudio, want: MediaKindVideo},
		{operation: OperationVideoConvert, want: MediaKindVideo},
		{operation: OperationVideoClip, want: MediaKindVideo},
		{operation: OperationImageGrayscale, want: MediaKindImage},
		{operation: OperationImageConvert, want: MediaKindImage},
		{operation: OperationImageCompress, want: MediaKindImage},
		{operation: OperationImageResize, want: MediaKindImage},
	}

	for _, test := range tests {
		t.Run(string(test.operation), func(t *testing.T) {
			got, err := test.operation.MediaKind()
			if err != nil {
				t.Fatalf("MediaKind() error = %v", err)
			}
			if got != test.want {
				t.Errorf("MediaKind() = %q, want %q", got, test.want)
			}
		})
	}
}

func TestOperationMediaKindRejectsInvalidOperation(t *testing.T) {
	operations := []Operation{
		"",
		"video_destroy_server",
		"image_unknown",
	}

	for _, operation := range operations {
		t.Run(string(operation), func(t *testing.T) {
			got, err := operation.MediaKind()

			if got != "" {
				t.Errorf("MediaKind() = %q, want empty media kind", got)
			}
			if !errors.Is(err, ErrInvalidOperation) {
				t.Fatalf("MediaKind() error = %v, want errors.Is(_, ErrInvalidOperation)", err)
			}
		})
	}
}
