package job

import (
	"errors"
	"fmt"
)

var ErrInvalidProgress = errors.New("invalid job progress")

type Progress uint8

func NewProgress(value int) (Progress, error) {
	if value < 0 || value > 100 {
		return 0, fmt.Errorf("%w: %d", ErrInvalidProgress, value)
	}

	return Progress(value), nil
}

func (p Progress) Value() int {
	return int(p)
}
