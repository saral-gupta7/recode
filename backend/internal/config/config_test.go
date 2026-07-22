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

func TestValidateRejectsEmptyHTTPAddress(t *testing.T) {
	cfg := Config{HTTPAddress: ""}

	if err := validate(cfg); err == nil {
		t.Fatal("validate() error = nil, want an error")
	}

}

func TestLoadOverridesDefaultsFromEnvironment(t *testing.T) {

	t.Setenv("RECODE_HTTP_ADDRESS", "127.0.0.1:9000")
	t.Setenv("RECODE_SHUTDOWN_TIMEOUT", "25s")


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


