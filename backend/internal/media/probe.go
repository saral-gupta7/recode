package media

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/saral-gupta7/recode/backend/internal/job"
)

var (
	ErrInvalidMedia = errors.New("invalid media")
	ErrWrongMedia   = errors.New("media does not support the requested operation")
)

type MediaInfo struct {
	HasVideo bool
	HasAudio bool
	Duration float64
}

type Prober interface {
	Probe(ctx context.Context, inputPath string) (MediaInfo, error)
}

type FFprobe struct {
	path string
}

func NewFFprobe(path string) *FFprobe {
	return &FFprobe{path: path}
}

func (p *FFprobe) Probe(ctx context.Context, inputPath string) (MediaInfo, error) {
	command := exec.CommandContext(
		ctx,
		p.path,
		"-v", "error",
		"-show_entries", "stream=codec_type",
		"-show_entries", "format=duration",
		"-of", "json",
		inputPath,
	)
	output, err := command.Output()
	if err != nil {
		return MediaInfo{}, fmt.Errorf("%w: ffprobe could not read the upload", ErrInvalidMedia)
	}

	var result struct {
		Streams []struct {
			CodecType string `json:"codec_type"`
		} `json:"streams"`
		Format struct {
			Duration string `json:"duration"`
		} `json:"format"`
	}
	if err := json.Unmarshal(output, &result); err != nil {
		return MediaInfo{}, fmt.Errorf("%w: decode ffprobe output: %v", ErrInvalidMedia, err)
	}

	info := MediaInfo{}
	for _, stream := range result.Streams {
		switch stream.CodecType {
		case "video":
			info.HasVideo = true
		case "audio":
			info.HasAudio = true
		}
	}
	if duration := strings.TrimSpace(result.Format.Duration); duration != "" &&
		duration != "N/A" {
		parsed, err := strconv.ParseFloat(duration, 64)
		if err != nil {
			return MediaInfo{}, fmt.Errorf("%w: invalid media duration", ErrInvalidMedia)
		}
		info.Duration = parsed
	}
	if !info.HasVideo && !info.HasAudio {
		return MediaInfo{}, fmt.Errorf("%w: no supported streams", ErrInvalidMedia)
	}
	return info, nil
}

func ValidateForOperation(operation job.Operation, info MediaInfo) error {
	switch operation {
	case job.OperationImageGrayscale,
		job.OperationImageConvert,
		job.OperationImageCompress,
		job.OperationImageResize:
		if !info.HasVideo || info.Duration > 0 {
			return ErrWrongMedia
		}

	case job.OperationVideoExtractAudio:
		if !info.HasVideo || !info.HasAudio || info.Duration <= 0 {
			return ErrWrongMedia
		}

	case job.OperationVideoGrayscale,
		job.OperationVideoRemoveAudio,
		job.OperationVideoConvert,
		job.OperationVideoClip:
		if !info.HasVideo || info.Duration <= 0 {
			return ErrWrongMedia
		}

	default:
		return ErrWrongMedia
	}
	return nil
}
