package media

import (
	"errors"
	"testing"

	"github.com/saral-gupta7/recode/backend/internal/job"
)

func TestValidateForOperation(t *testing.T) {
	tests := []struct {
		name      string
		operation job.Operation
		info      MediaInfo
		wantError bool
	}{
		{
			name:      "image operation accepts still image",
			operation: job.OperationImageResize,
			info:      MediaInfo{HasVideo: true},
		},
		{
			name:      "image operation rejects timed video",
			operation: job.OperationImageResize,
			info:      MediaInfo{HasVideo: true, Duration: 10},
			wantError: true,
		},
		{
			name:      "video operation accepts video",
			operation: job.OperationVideoConvert,
			info:      MediaInfo{HasVideo: true, Duration: 10},
		},
		{
			name:      "video operation rejects image",
			operation: job.OperationVideoConvert,
			info:      MediaInfo{HasVideo: true},
			wantError: true,
		},
		{
			name:      "audio extraction requires audio stream",
			operation: job.OperationVideoExtractAudio,
			info:      MediaInfo{HasVideo: true, Duration: 10},
			wantError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ValidateForOperation(test.operation, test.info)
			if test.wantError && !errors.Is(err, ErrWrongMedia) {
				t.Fatalf("ValidateForOperation() error = %v, want ErrWrongMedia", err)
			}
			if !test.wantError && err != nil {
				t.Fatalf("ValidateForOperation() error = %v, want nil", err)
			}
		})
	}
}
