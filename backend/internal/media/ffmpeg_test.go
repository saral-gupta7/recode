package media

import (
	"strings"
	"testing"

	"github.com/saral-gupta7/recode/backend/internal/job"
	"github.com/saral-gupta7/recode/backend/internal/task"
)

func TestCommandArgumentsSupportsEveryImageOperation(t *testing.T) {
	tests := []struct {
		operation job.Operation
		options   task.Options
		want      string
	}{
		{job.OperationImageGrayscale, task.Options{}, "format=gray"},
		{job.OperationImageConvert, task.Options{Format: "png"}, "-frames:v"},
		{job.OperationImageCompress, task.Options{Quality: 80}, "-q:v"},
		{job.OperationImageResize, task.Options{Width: 320}, "scale=320:-1"},
		{job.OperationImageCrop, task.Options{Width: 100, Height: 80, X: 2, Y: 3}, "crop=100:80:2:3"},
		{job.OperationImageRotate, task.Options{Angle: 90}, "transpose=1"},
		{job.OperationImageFlip, task.Options{FlipDirection: task.FlipDirectionVertical}, "vflip"},
		{job.OperationImageThumbnail, task.Options{Preset: "social"}, "crop=1200:630"},
		{job.OperationImageStripMetadata, task.Options{}, "-map_metadata"},
		{job.OperationImageAdjust, task.Options{Saturation: 100}, "eq=brightness"},
		{job.OperationImageBlur, task.Options{Strength: 2}, "gblur"},
		{job.OperationImageSharpen, task.Options{Strength: 2}, "unsharp"},
		{job.OperationImagePixelate, task.Options{BlockSize: 12}, "pixelize=w=12:h=12"},
		{job.OperationImagePadding, task.Options{PaddingTop: 2, PaddingRight: 3, PaddingBottom: 4, PaddingLeft: 5, Background: "white"}, "pad=iw+8:ih+6:5:2"},
	}
	for _, test := range tests {
		t.Run(string(test.operation), func(t *testing.T) {
			args, err := commandArguments(test.operation, test.options, "input", "output")
			if err != nil {
				t.Fatal(err)
			}
			joined := strings.Join(args, " ")
			if !strings.Contains(joined, test.want) {
				t.Fatalf("arguments = %q, want %q", joined, test.want)
			}
			if args[len(args)-1] != "output" {
				t.Fatalf("last argument = %q", args[len(args)-1])
			}
		})
	}
}
