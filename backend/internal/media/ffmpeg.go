package media

import (
	"context"
	"fmt"
	"math"
	"os/exec"
	"strconv"
	"strings"

	"github.com/saral-gupta7/recode/backend/internal/job"
	"github.com/saral-gupta7/recode/backend/internal/task"
)

const maxCommandErrorBytes = 4096

type FFmpeg struct {
	path string
}

var _ Processor = (*FFmpeg)(nil)

func NewFFmpeg(path string) *FFmpeg {
	return &FFmpeg{path: path}
}

func (p *FFmpeg) Process(
	ctx context.Context,
	operation job.Operation,
	options task.Options,
	inputPath string,
	outputPath string,
) error {
	args, err := commandArguments(operation, options, inputPath, outputPath)
	if err != nil {
		return err
	}

	command := exec.CommandContext(ctx, p.path, args...)
	output, err := command.CombinedOutput()
	if err != nil {
		if len(output) > maxCommandErrorBytes {
			output = output[len(output)-maxCommandErrorBytes:]
		}
		return fmt.Errorf("%w: %s", ErrProcessingFailed, strings.TrimSpace(string(output)))
	}
	return nil
}

func commandArguments(
	operation job.Operation,
	options task.Options,
	inputPath string,
	outputPath string,
) ([]string, error) {
	base := []string{"-y", "-hide_banner", "-loglevel", "error"}

	switch operation {
	case job.OperationImageGrayscale:
		return append(base, "-i", inputPath, "-vf", "format=gray", "-frames:v", "1", outputPath), nil

	case job.OperationImageConvert:
		return append(base, "-i", inputPath, "-frames:v", "1", outputPath), nil

	case job.OperationImageCompress:
		qualityScale := int(math.Round(31 - (float64(options.Quality)*29)/100))
		qualityScale = max(2, min(31, qualityScale))
		return append(
			base,
			"-i", inputPath,
			"-frames:v", "1",
			"-q:v", strconv.Itoa(qualityScale),
			outputPath,
		), nil

	case job.OperationImageResize:
		width := "-1"
		height := "-1"
		if options.Width > 0 {
			width = strconv.Itoa(options.Width)
		}
		if options.Height > 0 {
			height = strconv.Itoa(options.Height)
		}
		return append(
			base,
			"-i", inputPath,
			"-vf", "scale="+width+":"+height,
			"-frames:v", "1",
			outputPath,
		), nil

	case job.OperationVideoGrayscale:
		return append(
			base,
			"-i", inputPath,
			"-vf", "format=gray",
			"-c:v", "libx264",
			"-preset", "veryfast",
			"-crf", "23",
			"-c:a", "aac",
			outputPath,
		), nil

	case job.OperationVideoExtractAudio:
		args := append(base, "-i", inputPath, "-vn")
		switch options.Format {
		case "wav":
			args = append(args, "-c:a", "pcm_s16le")
		case "m4a":
			args = append(args, "-c:a", "aac")
		default:
			args = append(args, "-c:a", "libmp3lame")
		}
		return append(args, outputPath), nil

	case job.OperationVideoRemoveAudio:
		return append(
			base,
			"-i", inputPath,
			"-an",
			"-c:v", "libx264",
			"-preset", "veryfast",
			"-crf", "23",
			outputPath,
		), nil

	case job.OperationVideoConvert:
		args := append(base, "-i", inputPath)
		switch options.Format {
		case "webm":
			args = append(args, "-c:v", "libvpx-vp9", "-c:a", "libopus")
		default:
			args = append(args, "-c:v", "libx264", "-preset", "veryfast", "-crf", "23", "-c:a", "aac")
		}
		return append(args, outputPath), nil

	case job.OperationVideoClip:
		return append(
			base,
			"-ss", formatSeconds(options.StartSeconds),
			"-i", inputPath,
			"-t", formatSeconds(options.DurationSeconds),
			"-c:v", "libx264",
			"-preset", "veryfast",
			"-crf", "23",
			"-c:a", "aac",
			outputPath,
		), nil

	default:
		return nil, fmt.Errorf("%w: unsupported operation %q", ErrProcessingFailed, operation)
	}
}

func formatSeconds(value float64) string {
	return strconv.FormatFloat(value, 'f', 3, 64)
}
