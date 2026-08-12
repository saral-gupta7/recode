package job

import (
	"errors"
	"fmt"
)

type Operation string

var ErrInvalidOperation = errors.New("invalid operation")

const (
	OperationImageGrayscale     Operation = "image_grayscale"
	OperationImageConvert       Operation = "image_convert"
	OperationImageCompress      Operation = "image_compress"
	OperationImageResize        Operation = "image_resize"
	OperationImageCrop          Operation = "image_crop"
	OperationImageRotate        Operation = "image_rotate"
	OperationImageFlip          Operation = "image_flip"
	OperationImageThumbnail     Operation = "image_thumbnail"
	OperationImageStripMetadata Operation = "image_strip_metadata"
	OperationImageAdjust        Operation = "image_adjust"
	OperationImageBlur          Operation = "image_blur"
	OperationImageSharpen       Operation = "image_sharpen"
	OperationImagePixelate      Operation = "image_pixelate"
	OperationImagePadding       Operation = "image_padding"
)

func (o Operation) Valid() bool {
	switch o {
	case OperationImageGrayscale,
		OperationImageConvert,
		OperationImageCompress,
		OperationImageResize,
		OperationImageCrop,
		OperationImageRotate,
		OperationImageFlip,
		OperationImageThumbnail,
		OperationImageStripMetadata,
		OperationImageAdjust,
		OperationImageBlur,
		OperationImageSharpen,
		OperationImagePixelate,
		OperationImagePadding:
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
