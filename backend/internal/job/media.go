package job

import "fmt"

type MediaKind string

const (
	MediaKindImage MediaKind = "image"
	MediaKindVideo MediaKind = "video"
)

func (o Operation) MediaKind() (MediaKind, error) {
	switch o {
	case OperationVideoGrayscale,
		OperationVideoExtractAudio,
		OperationVideoRemoveAudio,
		OperationVideoConvert,
		OperationVideoClip:
		return MediaKindVideo, nil

	case OperationImageGrayscale,
		OperationImageConvert,
		OperationImageCompress,
		OperationImageResize:
		return MediaKindImage, nil

	default:
		return "", fmt.Errorf("%w: %q", ErrInvalidOperation, o)
	}
}
