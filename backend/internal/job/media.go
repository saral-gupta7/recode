package job

import "fmt"

type MediaKind string

const MediaKindImage MediaKind = "image"

func (o Operation) MediaKind() (MediaKind, error) {
	if o.Valid() {
		return MediaKindImage, nil
	}
	return "", fmt.Errorf("%w: %q", ErrInvalidOperation, o)
}
