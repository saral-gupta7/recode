package config

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/knadh/koanf/providers/env/v2"
	"github.com/knadh/koanf/v2"
)

const (
	defaultHTTPAddress           = ":8080"
	delimiter                    = "."
	environmentPrefix            = "RECODE_"
	defaultShutdownTimeout       = "10s"
	defaultDatabaseURL           = "postgres://recode:recode_dev@localhost:5433/recode?sslmode=disable"
	defaultStorageRoot           = "./data"
	defaultMaxUploadBytes  int64 = 1 << 30
	defaultRedisAddress          = "localhost:6379"
	defaultQueueName             = "recode:jobs"
	defaultResultRetention       = "2h"
	defaultFFmpegPath            = "ffmpeg"
	defaultFFprobePath           = "ffprobe"
)

type Config struct {
	HTTPAddress     string
	ShutdownTimeout time.Duration
	DatabaseURL     string
	StorageRoot     string
	MaxUploadBytes  int64
	RedisAddress    string
	QueueName       string
	ResultRetention time.Duration
	FFmpegPath      string
	FFprobePath     string
}

func enviromentEntry(key, value string) (string, any) {
	key = strings.TrimPrefix(key, environmentPrefix)
	key = strings.ToLower(key)

	key = strings.ReplaceAll(key, "_", delimiter)
	return key, value
}

func decode(k *koanf.Koanf) (Config, error) {
	timeoutText := k.String("shutdown.timeout")

	timeout, err := time.ParseDuration(timeoutText)
	if err != nil {
		return Config{}, fmt.Errorf(
			"shutdown.timeout: parse duration %q: %w",
			timeoutText,
			err,
		)

	}

	retentionText := k.String("result.retention")
	retention, err := time.ParseDuration(retentionText)
	if err != nil {
		return Config{}, fmt.Errorf(
			"result.retention: parse duration %q: %w",
			retentionText,
			err,
		)
	}

	return Config{
		HTTPAddress:     k.String("http.address"),
		ShutdownTimeout: timeout,
		DatabaseURL:     k.String("database.url"),
		StorageRoot:     k.String("storage.root"),
		MaxUploadBytes:  k.Int64("upload.max.bytes"),
		RedisAddress:    k.String("redis.address"),
		QueueName:       k.String("queue.name"),
		ResultRetention: retention,
		FFmpegPath:      k.String("ffmpeg.path"),
		FFprobePath:     k.String("ffprobe.path"),
	}, nil
}

func Load() (Config, error) {
	k := koanf.New(delimiter)

	if err := k.Set("http.address", defaultHTTPAddress); err != nil {
		return Config{}, fmt.Errorf("set default http.address default: %w", err)
	}

	if err := k.Set("shutdown.timeout", defaultShutdownTimeout); err != nil {
		return Config{}, fmt.Errorf("set default shutdown.timeout default: %w", err)
	}

	if err := k.Set("database.url", defaultDatabaseURL); err != nil {
		return Config{}, fmt.Errorf("set default database.url default: %w", err)
	}

	if err := k.Set("storage.root", defaultStorageRoot); err != nil {
		return Config{}, fmt.Errorf("set default storage.root: %w", err)
	}

	if err := k.Set("upload.max.bytes", defaultMaxUploadBytes); err != nil {
		return Config{}, fmt.Errorf("set default upload.max.bytes: %w", err)
	}

	if err := k.Set("redis.address", defaultRedisAddress); err != nil {
		return Config{}, fmt.Errorf("set default redis.address: %w", err)
	}

	if err := k.Set("queue.name", defaultQueueName); err != nil {
		return Config{}, fmt.Errorf("set default queue.name: %w", err)
	}

	if err := k.Set("result.retention", defaultResultRetention); err != nil {
		return Config{}, fmt.Errorf("set default result.retention: %w", err)
	}

	if err := k.Set("ffmpeg.path", defaultFFmpegPath); err != nil {
		return Config{}, fmt.Errorf("set default ffmpeg.path: %w", err)
	}

	if err := k.Set("ffprobe.path", defaultFFprobePath); err != nil {
		return Config{}, fmt.Errorf("set default ffprobe.path: %w", err)
	}

	provider := env.Provider(delimiter, env.Opt{
		Prefix:        environmentPrefix,
		TransformFunc: enviromentEntry,
	})

	if err := k.Load(provider, nil); err != nil {
		return Config{}, fmt.Errorf("load environment configuration: %w", err)
	}

	cfg, err := decode(k)
	if err != nil {
		return Config{}, fmt.Errorf("decode configuration: %w", err)
	}

	if err := validate(cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func validate(cfg Config) error {

	if cfg.HTTPAddress == "" {
		return errors.New("http.address must not be empty")
	}

	if cfg.ShutdownTimeout <= 0 {
		return errors.New("shutdown.timeout must be greater than zero")
	}

	databaseURL, err := url.Parse(cfg.DatabaseURL)

	if err != nil ||
		(databaseURL.Scheme != "postgres" &&
			databaseURL.Scheme != "postgresql") ||
		databaseURL.Host == "" ||
		strings.Trim(databaseURL.Path, "/") == "" {
		return errors.New("database.url must be a valid PostgreSQL URL")
	}

	if strings.TrimSpace(cfg.StorageRoot) == "" {
		return errors.New("storage.root must not be empty")
	}

	if cfg.MaxUploadBytes <= 0 {
		return errors.New("upload.max.bytes must be greater than zero")
	}

	if strings.TrimSpace(cfg.RedisAddress) == "" {
		return errors.New("redis.address must not be empty")
	}

	if strings.TrimSpace(cfg.QueueName) == "" {
		return errors.New("queue.name must not be empty")
	}

	if cfg.ResultRetention <= 0 {
		return errors.New("result.retention must be greater than zero")
	}

	if strings.TrimSpace(cfg.FFmpegPath) == "" {
		return errors.New("ffmpeg.path must not be empty")
	}

	if strings.TrimSpace(cfg.FFprobePath) == "" {
		return errors.New("ffprobe.path must not be empty")
	}

	return nil
}
