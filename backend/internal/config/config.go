package config

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/knadh/koanf/providers/env/v2"
	"github.com/knadh/koanf/v2"
)

const (
	defaultHTTPAddress = ":8080"
	delimiter = "."
	environmentPrefix = "RECODE_"
	defaultShutdownTimeout = "10s"
)

type Config struct {
	HTTPAddress string
	ShutdownTimeout time.Duration
}


func enviromentEntry(key,value string) (string, any) {
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
		
	return Config{HTTPAddress: k.String("http.address"), ShutdownTimeout: timeout,}, nil
}

func Load() (Config, error) {
	k := koanf.New(delimiter)

	if err := k.Set("http.address", defaultHTTPAddress); err != nil {
		return Config{}, fmt.Errorf("set default http.address default: %w", err)
	}

	if err := k.Set("shutdown.timeout", defaultShutdownTimeout); err != nil {
		return Config{}, fmt.Errorf("set default shutdown.timeout default: %w", err)
	}
	
	provider := env.Provider(delimiter, env.Opt{
		Prefix: environmentPrefix,
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

	return nil
}
