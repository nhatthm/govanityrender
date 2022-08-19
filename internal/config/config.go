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

const defaultRef = "master"

// Config is the configuration for the application.
type Config struct {
	PageTitle    string       `yaml:"page_title" toml:"page_title" json:"page_title"`
	Host         string       `yaml:"host" toml:"host" json:"host"`
	Repositories []Repository `yaml:"repositories" toml:"repositories" json:"repositories"`
}

// Repository is the configuration for a repository.
type Repository struct {
	Library    string   `yaml:"library" toml:"library" json:"library"`
	Path       string   `yaml:"path" toml:"path" json:"path"`
	Repository string   `yaml:"repository" toml:"repository" json:"repository"`
	Ref        string   `yaml:"ref" toml:"ref" json:"ref"`
	Version    string   `yaml:"-" toml:"-" json:"-"`
	Deprecated string   `yaml:"deprecated" toml:"deprecated" json:"deprecated"`
	Modules    []Module `yaml:"-" toml:"-" json:"-"`
}

// Module is the configuration for a module.
type Module struct {
	Path string
}

// HydrateFile reads the configuration from a file and hydrates it.
func HydrateFile(file string, hydrators ...Hydrator) (Config, error) {
	cfg, err := FromFile(file)
	if err != nil {
		return Config{}, err
	}

	for _, h := range hydrators {
		if err := h.Hydrate(&cfg); err != nil {
			return Config{}, err
		}
	}

	return cfg, nil
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

	for i := range cfg.Repositories {
		hydrateRepository(&cfg.Repositories[i])
	}
}

func hydrateRepository(r *Repository) {
	if len(r.Ref) == 0 {
		r.Ref = defaultRef
	}
}
