package job

import (
	"errors"
	"fmt"
)

type Operation string

var ErrInvalidOperation = errors.New("invalid operation")

const (
	OperationVideoGrayscale    Operation = "video_grayscale"
	OperationVideoExtractAudio Operation = "video_extract_audio"
	OperationVideoRemoveAudio  Operation = "video_remove_audio"
	OperationVideoConvert      Operation = "video_convert"
	OperationVideoClip         Operation = "video_clip"

	OperationImageGrayscale Operation = "image_grayscale"
	OperationImageConvert   Operation = "image_convert"
	OperationImageCompress  Operation = "image_compress"
	OperationImageResize    Operation = "image_resize"
)

func (o Operation) Valid() bool {
	switch o {
	case OperationVideoGrayscale,
		OperationVideoExtractAudio,
		OperationVideoRemoveAudio,
		OperationVideoConvert,
		OperationVideoClip,
		OperationImageGrayscale,
		OperationImageConvert,
		OperationImageCompress,
		OperationImageResize:
		return true

	default:
		return false
	}
}

func ParseOperation(value string) (Operation, error) {

	operation := Operation(value)

	if !operation.Valid() {
		return "", fmt.Errorf("%w: %q", ErrInvalidOperation, value)
	}

	return operation, nil
}
