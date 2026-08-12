package media

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os/exec"
	"strconv"
	"strings"

	"github.com/saral-gupta7/recode/backend/internal/job"
	"github.com/saral-gupta7/recode/backend/internal/task"
)

var (
	ErrInvalidMedia = errors.New("invalid media")
	ErrWrongMedia   = errors.New("media does not support the requested operation")
)

const maxDecodedPixels int64 = 100_000_000

type MediaInfo struct {
	HasImage bool
	Width    int
	Height   int
	Frames   int
	Format   string
}

type Prober interface {
	Probe(ctx context.Context, inputPath string) (MediaInfo, error)
}
type FFprobe struct{ path string }

func NewFFprobe(path string) *FFprobe { return &FFprobe{path: path} }

func (p *FFprobe) Probe(ctx context.Context, inputPath string) (MediaInfo, error) {
	output, err := exec.CommandContext(ctx, p.path,
		"-v", "error", "-count_frames", "-select_streams", "v:0",
		"-show_entries", "stream=codec_type,codec_name,width,height,nb_frames,nb_read_frames:format=format_name",
		"-of", "json", inputPath).Output()
	if err != nil {
		return MediaInfo{}, fmt.Errorf("%w: ffprobe could not read the upload", ErrInvalidMedia)
	}
	var result struct {
		Streams []struct {
			CodecType  string `json:"codec_type"`
			CodecName  string `json:"codec_name"`
			Width      int    `json:"width"`
			Height     int    `json:"height"`
			Frames     string `json:"nb_frames"`
			ReadFrames string `json:"nb_read_frames"`
		} `json:"streams"`
		Format struct {
			Name string `json:"format_name"`
		} `json:"format"`
	}
	if err := json.Unmarshal(output, &result); err != nil {
		return MediaInfo{}, fmt.Errorf("%w: decode ffprobe output: %v", ErrInvalidMedia, err)
	}
	if len(result.Streams) == 0 {
		return MediaInfo{}, fmt.Errorf("%w: no image stream", ErrInvalidMedia)
	}
	stream := result.Streams[0]
	frames := 1
	value := strings.TrimSpace(stream.ReadFrames)
	if value == "" || value == "N/A" {
		value = strings.TrimSpace(stream.Frames)
	}
	if value != "" && value != "N/A" {
		if parsed, parseErr := strconv.Atoi(value); parseErr == nil {
			frames = parsed
		}
	}
	return MediaInfo{HasImage: stream.CodecType == "video", Width: stream.Width, Height: stream.Height, Frames: frames, Format: stream.CodecName}, nil
}

func ValidateForOperation(operation job.Operation, info MediaInfo) error {
	if !operation.Valid() || !info.HasImage || !supportedImageCodec(info.Format) || info.Width <= 0 || info.Height <= 0 || info.Width > 16384 || info.Height > 16384 || int64(info.Width)*int64(info.Height) > maxDecodedPixels || info.Frames > 1 {
		return ErrWrongMedia
	}
	return nil
}

func ValidateTransform(operation job.Operation, options task.Options, info MediaInfo) error {
	width, height := int64(info.Width), int64(info.Height)
	switch operation {
	case job.OperationImageCrop:
		if int64(options.X)+int64(options.Width) > width || int64(options.Y)+int64(options.Height) > height {
			return ErrWrongMedia
		}
	case job.OperationImageResize:
		if options.Width > 0 {
			width = int64(options.Width)
		}
		if options.Height > 0 {
			height = int64(options.Height)
		}
		if options.Width == 0 {
			width = int64(math.Ceil(float64(info.Width) * float64(options.Height) / float64(info.Height)))
		}
		if options.Height == 0 {
			height = int64(math.Ceil(float64(info.Height) * float64(options.Width) / float64(info.Width)))
		}
	case job.OperationImagePadding:
		width += int64(options.PaddingLeft + options.PaddingRight)
		height += int64(options.PaddingTop + options.PaddingBottom)
	}
	if width <= 0 || height <= 0 || width > 16384 || height > 16384 || width*height > maxDecodedPixels {
		return ErrWrongMedia
	}
	return nil
}

func supportedImageCodec(codec string) bool {
	switch strings.ToLower(codec) {
	case "mjpeg", "png", "webp":
		return true
	default:
		return false
	}
}
