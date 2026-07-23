package media

import (
	"context"
	"errors"

	"github.com/saral-gupta7/recode/backend/internal/job"
	"github.com/saral-gupta7/recode/backend/internal/task"
)

var ErrProcessingFailed = errors.New("media processing failed")

type Processor interface {
	Process(
		ctx context.Context,
		operation job.Operation,
		options task.Options,
		inputPath string,
		outputPath string,
	) error
}
