package media

import (
	"strings"
	"testing"

	"github.com/saral-gupta7/recode/backend/internal/job"
	"github.com/saral-gupta7/recode/backend/internal/task"
)

func TestCommandArgumentsSupportsEveryOperation(t *testing.T) {
	tests := []struct {
		operation job.Operation
		options   task.Options
		want      string
	}{
		{operation: job.OperationImageGrayscale, want: "format=gray"},
		{operation: job.OperationImageConvert},
		{operation: job.OperationImageCompress, options: task.Options{Quality: 80}, want: "-q:v"},
		{operation: job.OperationImageResize, options: task.Options{Width: 320}, want: "scale=320:-1"},
		{operation: job.OperationVideoGrayscale, want: "format=gray"},
		{operation: job.OperationVideoExtractAudio, options: task.Options{Format: "mp3"}, want: "-vn"},
		{operation: job.OperationVideoRemoveAudio, want: "-an"},
		{operation: job.OperationVideoConvert, options: task.Options{Format: "mp4"}, want: "libx264"},
		{operation: job.OperationVideoClip, options: task.Options{StartSeconds: 1, DurationSeconds: 2}, want: "-t"},
	}

	for _, test := range tests {
		t.Run(string(test.operation), func(t *testing.T) {
			args, err := commandArguments(test.operation, test.options, "input", "output")
			if err != nil {
				t.Fatalf("commandArguments() error = %v", err)
			}
			joined := strings.Join(args, " ")
			if test.want != "" && !strings.Contains(joined, test.want) {
				t.Errorf("arguments = %q, want %q", joined, test.want)
			}
			if args[len(args)-1] != "output" {
				t.Errorf("last argument = %q, want output", args[len(args)-1])
			}
		})
	}
}
