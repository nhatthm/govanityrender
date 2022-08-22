package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	xerrors "go.nhat.io/vanityrender/internal/errors"
)

const (
	// ErrCouldNotReadConfigFile indicates that the configuration file could not be read.
	ErrCouldNotReadConfigFile = xerrors.Error("could not read config file")
	// ErrInvalidConfig indicates that the configuration is invalid.
	ErrInvalidConfig = xerrors.Error("invalid config")
	// ErrMissingHost indicates that the host is missing.
	ErrMissingHost = xerrors.Error("missing host")
)

// Config is the configuration for the application.
type Config struct {
	PageTitle       string       `json:"page_title"`
	PageDescription string       `json:"page_description"`
	Host            string       `json:"host"`
	Repositories    []Repository `json:"repositories"`
}

// Repository is the configuration for a repository.
type Repository struct {
	Name       string `json:"name"`
	Path       string `json:"path"`
	Repository string `json:"repository"`
	Ref        string `json:"ref"`
	Deprecated string `json:"deprecated"`
}

// FromFile reads the configuration from a file.
func FromFile(file string) (Config, error) {
	f, err := os.Open(filepath.Clean(file))
	if err != nil {
		return Config{}, fmt.Errorf("%w: %s", ErrCouldNotReadConfigFile, err.Error())
	}

	defer func() {
		_ = f.Close() // nolint: errcheck
	}()

	var cfg Config

	dec := json.NewDecoder(f)

	if err := dec.Decode(&cfg); err != nil {
		return Config{}, fmt.Errorf("%w: %s", ErrInvalidConfig, err.Error())
	}

	if err := Validate(cfg); err != nil {
		return Config{}, err
	}

	hydrateConfig(&cfg)

	return cfg, nil
}

// Validate validates the configuration.
func Validate(config Config) error {
	if config.Host == "" {
		return ErrMissingHost
	}

	return nil
}

func hydrateConfig(cfg *Config) {
	if len(cfg.PageTitle) == 0 {
		cfg.PageTitle = cfg.Host
	}
}
