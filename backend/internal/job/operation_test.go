package job

import (
	"errors"
	"testing"
)

func TestEveryImageOperationIsValidAndParsable(t *testing.T) {
	operations := []Operation{
		OperationImageGrayscale, OperationImageConvert, OperationImageCompress,
		OperationImageResize, OperationImageCrop, OperationImageRotate,
		OperationImageFlip, OperationImageThumbnail, OperationImageStripMetadata,
		OperationImageAdjust, OperationImageBlur, OperationImageSharpen,
		OperationImagePixelate, OperationImagePadding,
	}
	for _, operation := range operations {
		t.Run(string(operation), func(t *testing.T) {
			if !operation.Valid() {
				t.Fatalf("%q should be valid", operation)
			}
			parsed, err := ParseOperation(string(operation))
			if err != nil || parsed != operation {
				t.Fatalf("ParseOperation() = %q, %v", parsed, err)
			}
		})
	}
}

func TestVideoAndUnknownOperationsAreRejected(t *testing.T) {
	for _, value := range []string{"", "video_clip", "unknown_operation"} {
		if _, err := ParseOperation(value); !errors.Is(err, ErrInvalidOperation) {
			t.Fatalf("ParseOperation(%q) error = %v, want ErrInvalidOperation", value, err)
		}
	}
}
