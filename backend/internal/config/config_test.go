package config

import (
	"strings"
	"testing"
	"time"
)

func TestLoadUsesDefaultHTTPAddress(t *testing.T) {

	wantAddress := ":8080"

	cfg, err := Load()

	if err != nil {
		t.Fatalf("Load() returned an unexpected error: %v", err)
	}

	if cfg.HTTPAddress != wantAddress {
		t.Fatalf("Load() = %q, want %q", cfg.HTTPAddress, wantAddress)
	}

}

func TestLoadUsesDefaultDatabaseURL(t *testing.T) {
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned an unexpected error: %v", err)
	}

	if cfg.DatabaseURL != defaultDatabaseURL {
		t.Fatalf(
			"Load().DatabaseURL = %q, want %q",
			cfg.DatabaseURL,
			defaultDatabaseURL,
		)
	}
}

func TestValidateRejectsEmptyHTTPAddress(t *testing.T) {
	cfg := Config{HTTPAddress: ""}

	if err := validate(cfg); err == nil {
		t.Fatal("validate() error = nil, want an error")
	}

}

func TestLoadOverridesDefaultsFromEnvironment(t *testing.T) {

	t.Setenv("RECODE_HTTP_ADDRESS", "127.0.0.1:9000")
	t.Setenv("RECODE_SHUTDOWN_TIMEOUT", "25s")
	t.Setenv(
		"RECODE_DATABASE_URL",
		"postgres://app:secret@database:5432/recode?sslmode=require",
	)

	cfg, err := Load()

	if err != nil {
		t.Fatalf("Load() returned an unexpected error: %v", err)
	}

	if cfg.HTTPAddress != "127.0.0.1:9000" {
		t.Fatalf(
			"Load().HTTPAddress = %q, want %q",
			cfg.HTTPAddress,
			"127.0.0.1:9000",
		)
	}

	if cfg.ShutdownTimeout != 25*time.Second {
		t.Fatalf(
			"Load().ShutdownTimeout = %s, want %s",
			cfg.ShutdownTimeout,
			25*time.Second,
		)
	}

	if cfg.DatabaseURL !=
		"postgres://app:secret@database:5432/recode?sslmode=require" {
		t.Fatalf(
			"Load().DatabaseURL = %q, want environment value",
			cfg.DatabaseURL,
		)
	}

}

func TestValidateRejectsInvalidDatabaseURL(t *testing.T) {
	tests := []struct {
		name        string
		databaseURL string
	}{
		{name: "empty", databaseURL: ""},
		{name: "malformed", databaseURL: "://bad"},
		{name: "wrong scheme", databaseURL: "mysql://localhost/recode"},
		{name: "missing host", databaseURL: "postgres:///recode"},
		{name: "missing database", databaseURL: "postgres://localhost"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cfg := Config{
				HTTPAddress:     ":8080",
				ShutdownTimeout: 10 * time.Second,
				DatabaseURL:     test.databaseURL,
			}

			err := validate(cfg)
			if err == nil {
				t.Fatal("validate() error = nil, want an error")
			}
			if !strings.Contains(err.Error(), "database.url") {
				t.Fatalf(
					"validate() error = %q, want it to contain database.url",
					err,
				)
			}
		})
	}
}

func TestLoadUsesDefaultStorageConfiguration(t *testing.T) {
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned an unexpected error: %v", err)
	}

	if cfg.StorageRoot != defaultStorageRoot {
		t.Errorf(
			"Load().StorageRoot = %q, want %q",
			cfg.StorageRoot,
			defaultStorageRoot,
		)
	}
	if cfg.MaxUploadBytes != defaultMaxUploadBytes {
		t.Errorf(
			"Load().MaxUploadBytes = %d, want %d",
			cfg.MaxUploadBytes,
			defaultMaxUploadBytes,
		)
	}
}

func TestLoadUsesDefaultWorkerConfiguration(t *testing.T) {
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned an unexpected error: %v", err)
	}

	if cfg.RedisAddress != defaultRedisAddress {
		t.Errorf("RedisAddress = %q, want %q", cfg.RedisAddress, defaultRedisAddress)
	}
	if cfg.QueueName != defaultQueueName {
		t.Errorf("QueueName = %q, want %q", cfg.QueueName, defaultQueueName)
	}
	if cfg.ResultRetention != 2*time.Hour {
		t.Errorf("ResultRetention = %s, want 2h", cfg.ResultRetention)
	}
	if cfg.FFmpegPath != defaultFFmpegPath || cfg.FFprobePath != defaultFFprobePath {
		t.Errorf("processor paths = %q/%q, want defaults", cfg.FFmpegPath, cfg.FFprobePath)
	}
}

func TestLoadOverridesWorkerConfigurationFromEnvironment(t *testing.T) {
	t.Setenv("RECODE_REDIS_ADDRESS", "cache:6380")
	t.Setenv("RECODE_QUEUE_NAME", "jobs:test")
	t.Setenv("RECODE_RESULT_RETENTION", "45m")
	t.Setenv("RECODE_FFMPEG_PATH", "/opt/bin/ffmpeg")
	t.Setenv("RECODE_FFPROBE_PATH", "/opt/bin/ffprobe")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned an unexpected error: %v", err)
	}

	if cfg.RedisAddress != "cache:6380" ||
		cfg.QueueName != "jobs:test" ||
		cfg.ResultRetention != 45*time.Minute ||
		cfg.FFmpegPath != "/opt/bin/ffmpeg" ||
		cfg.FFprobePath != "/opt/bin/ffprobe" {
		t.Fatalf("Load() worker configuration = %#v", cfg)
	}
}

func TestLoadOverridesStorageConfigurationFromEnvironment(t *testing.T) {
	t.Setenv("RECODE_STORAGE_ROOT", "/var/lib/recode")
	t.Setenv("RECODE_UPLOAD_MAX_BYTES", "524288000")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned an unexpected error: %v", err)
	}

	if cfg.StorageRoot != "/var/lib/recode" {
		t.Errorf(
			"Load().StorageRoot = %q, want /var/lib/recode",
			cfg.StorageRoot,
		)
	}
	if cfg.MaxUploadBytes != 524288000 {
		t.Errorf(
			"Load().MaxUploadBytes = %d, want 524288000",
			cfg.MaxUploadBytes,
		)
	}
}

func TestValidateRejectsInvalidStorageConfiguration(t *testing.T) {
	valid := Config{
		HTTPAddress:     ":8080",
		ShutdownTimeout: 10 * time.Second,
		DatabaseURL:     defaultDatabaseURL,
		StorageRoot:     defaultStorageRoot,
		MaxUploadBytes:  defaultMaxUploadBytes,
	}

	tests := []struct {
		name   string
		mutate func(*Config)
		want   string
	}{
		{
			name: "empty storage root",
			mutate: func(cfg *Config) {
				cfg.StorageRoot = " "
			},
			want: "storage.root",
		},
		{
			name: "zero upload limit",
			mutate: func(cfg *Config) {
				cfg.MaxUploadBytes = 0
			},
			want: "upload.max.bytes",
		},
		{
			name: "negative upload limit",
			mutate: func(cfg *Config) {
				cfg.MaxUploadBytes = -1
			},
			want: "upload.max.bytes",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cfg := valid
			test.mutate(&cfg)

			err := validate(cfg)
			if err == nil {
				t.Fatal("validate() error = nil, want an error")
			}
			if !strings.Contains(err.Error(), test.want) {
				t.Fatalf(
					"validate() error = %q, want it to contain %q",
					err,
					test.want,
				)
			}
		})
	}
}

func TestLoadRejectsInvalidShutdownTimeout(t *testing.T) {
	t.Setenv("RECODE_SHUTDOWN_TIMEOUT", "tomorrow")

	cfg, err := Load()

	if err == nil {
		t.Fatal("Load() error = nil, want an error")
	}

	if cfg != (Config{}) {
		t.Fatalf("Load() config = %#v, want zero Config", cfg)
	}

	if !strings.Contains(err.Error(), "shutdown.timeout") {
		t.Fatalf(
			"Load() error = %q, want it to contain %q",
			err,
			"shutdown.timeout",
		)
	}
}

func TestValidateRejectsNonPositiveShutdownTimeout(t *testing.T) {
	tests := []struct {
		name    string
		timeout time.Duration
	}{
		{
			name:    "zero",
			timeout: 0,
		},
		{
			name:    "negative",
			timeout: -1 * time.Second,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cfg := Config{
				HTTPAddress:     ":8080",
				ShutdownTimeout: test.timeout,
			}

			err := validate(cfg)
			if err == nil {
				t.Fatal("validate() error = nil, want an error")
			}

			if !strings.Contains(err.Error(), "shutdown.timeout") {
				t.Fatalf(
					"validate() error = %q, want it to contain %q",
					err,
					"shutdown.timeout",
				)
			}
		})
	}
}

func TestEnvironmentEntryTransformsRecodeKey(t *testing.T) {
	wantKey := "http.address"
	wantValue := "127.0.0.1:9000"

	gotKey, gotValue := enviromentEntry(
		"RECODE_HTTP_ADDRESS",
		wantValue,
	)

	if gotKey != wantKey {
		t.Fatalf("enviromentEntry() key = %q, want %q", gotKey, wantKey)
	}

	if gotValue != wantValue {
		t.Fatalf("enviromentEntry() value = %#v, want %q", gotValue, wantValue)
	}
}
