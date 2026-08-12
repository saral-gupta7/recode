package media

import (
	"errors"
	"testing"

	"github.com/saral-gupta7/recode/backend/internal/job"
)

func TestValidateForOperation(t *testing.T) {
	tests := []struct {
		name      string
		info      MediaInfo
		wantError bool
	}{
		{"valid still image", MediaInfo{HasImage: true, Width: 1200, Height: 800, Frames: 1, Format: "png"}, false},
		{"missing dimensions", MediaInfo{HasImage: true, Format: "png"}, true},
		{"too many pixels", MediaInfo{HasImage: true, Width: 15000, Height: 15000, Frames: 1, Format: "png"}, true},
		{"animated image", MediaInfo{HasImage: true, Width: 320, Height: 240, Frames: 2, Format: "png"}, true},
		{"video codec", MediaInfo{HasImage: true, Width: 320, Height: 240, Frames: 1, Format: "h264"}, true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ValidateForOperation(job.OperationImageCrop, test.info)
			if test.wantError && !errors.Is(err, ErrWrongMedia) {
				t.Fatalf("error = %v", err)
			}
			if !test.wantError && err != nil {
				t.Fatal(err)
			}
		})
	}
}
