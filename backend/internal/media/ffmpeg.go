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

type FFmpeg struct{ path string }

var _ ImageProcessor = (*FFmpeg)(nil)

func NewFFmpeg(path string) *FFmpeg { return &FFmpeg{path: path} }

func (p *FFmpeg) Process(ctx context.Context, operation job.Operation, options task.Options, inputPath, outputPath string) error {
	args, err := commandArguments(operation, options, inputPath, outputPath)
	if err != nil {
		return err
	}
	output, err := exec.CommandContext(ctx, p.path, args...).CombinedOutput()
	if err != nil {
		if len(output) > maxCommandErrorBytes {
			output = output[len(output)-maxCommandErrorBytes:]
		}
		return fmt.Errorf("%w: %s", ErrProcessingFailed, strings.TrimSpace(string(output)))
	}
	return nil
}

func commandArguments(operation job.Operation, options task.Options, inputPath, outputPath string) ([]string, error) {
	base := []string{"-y", "-hide_banner", "-loglevel", "error", "-i", inputPath}
	finish := func(filters string, extra ...string) []string {
		args := append([]string{}, base...)
		if filters != "" {
			args = append(args, "-vf", filters)
		}
		args = append(args, "-frames:v", "1")
		args = append(args, extra...)
		return append(args, outputPath)
	}

	switch operation {
	case job.OperationImageGrayscale:
		return finish("format=gray"), nil
	case job.OperationImageConvert:
		return finish(""), nil
	case job.OperationImageCompress:
		qualityScale := int(math.Round(31 - (float64(options.Quality)*29)/100))
		qualityScale = max(2, min(31, qualityScale))
		return finish("", "-q:v", strconv.Itoa(qualityScale)), nil
	case job.OperationImageResize:
		width, height := "-1", "-1"
		if options.Width > 0 {
			width = strconv.Itoa(options.Width)
		}
		if options.Height > 0 {
			height = strconv.Itoa(options.Height)
		}
		return finish("scale=" + width + ":" + height), nil
	case job.OperationImageCrop:
		filter := fmt.Sprintf("crop=%d:%d:%d:%d", options.Width, options.Height, options.X, options.Y)
		return finish(filter), nil
	case job.OperationImageRotate:
		filter := map[int]string{90: "transpose=1", 180: "hflip,vflip", 270: "transpose=2"}[options.Angle]
		return finish(filter), nil
	case job.OperationImageFlip:
		filter := "hflip"
		if options.FlipDirection == task.FlipDirectionVertical {
			filter = "vflip"
		}
		return finish(filter), nil
	case job.OperationImageThumbnail:
		width, height := thumbnailDimensions(options.Preset)
		filter := fmt.Sprintf("scale=%d:%d:force_original_aspect_ratio=increase,crop=%d:%d", width, height, width, height)
		return finish(filter, "-q:v", "3"), nil
	case job.OperationImageStripMetadata:
		return finish("", "-map_metadata", "-1", "-q:v", "2"), nil
	case job.OperationImageAdjust:
		brightness := float64(options.Brightness) / 100
		contrast := 1 + float64(options.Contrast)/100
		saturation := float64(options.Saturation) / 100
		filter := fmt.Sprintf("eq=brightness=%s:contrast=%s:saturation=%s", decimal(brightness), decimal(contrast), decimal(saturation))
		return finish(filter), nil
	case job.OperationImageBlur:
		return finish("gblur=sigma=" + decimal(options.Strength)), nil
	case job.OperationImageSharpen:
		return finish("unsharp=5:5:" + decimal(options.Strength)), nil
	case job.OperationImagePixelate:
		n := strconv.Itoa(options.BlockSize)
		return finish("pixelize=w=" + n + ":h=" + n + ":mode=avg"), nil
	case job.OperationImagePadding:
		filter := fmt.Sprintf("pad=iw+%d:ih+%d:%d:%d:color=%s", options.PaddingLeft+options.PaddingRight, options.PaddingTop+options.PaddingBottom, options.PaddingLeft, options.PaddingTop, options.Background)
		return finish(filter), nil
	default:
		return nil, fmt.Errorf("%w: unsupported operation %q", ErrProcessingFailed, operation)
	}
}

func thumbnailDimensions(preset string) (int, int) {
	switch preset {
	case "preview":
		return 640, 360
	case "social":
		return 1200, 630
	default:
		return 256, 256
	}
}

func decimal(value float64) string { return strconv.FormatFloat(value, 'f', 3, 64) }
