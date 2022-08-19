package renderer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aymerick/raymond"

	"go.nhat.io/vanityrender/internal/config"
)

const indexFile = `index.html`

var repositorySanitizer = strings.NewReplacer("https://", "", "http://", "")

// Renderder is the interface for rendering.
type Renderder interface {
	Render(config config.Config) error
}

var _ Renderder = (*HandlebarsRenderder)(nil)

// HandlebarsRenderder renders to the filesystem.
type HandlebarsRenderder struct {
	homepageTpl   *raymond.Template
	repositoryTpl *raymond.Template
	outputDir     string
}

// Render renders the configuration.
func (h *HandlebarsRenderder) Render(cfg config.Config) error {
	for i := range cfg.Repositories {
		cfg.Repositories[i].Repository = repositorySanitizer.Replace(cfg.Repositories[i].Repository)
	}

	if err := h.renderHomepage(cfg); err != nil {
		return fmt.Errorf("could not render homepage: %w", err)
	}

	for _, r := range cfg.Repositories {
		if err := h.renderLibrary(cfg.Host, r); err != nil {
			return err
		}
	}

	return nil
}

func (h *HandlebarsRenderder) renderHomepage(cfg config.Config) error {
	homepageFile := filepath.Join(h.outputDir, indexFile)

	result, err := h.homepageTpl.Exec(cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(homepageFile, []byte(result), 0o644) // nolint: gosec
}

func (h *HandlebarsRenderder) renderLibrary(host string, r config.Repository) error {
	if err := h.renderModule(host, r, ""); err != nil {
		return err
	}

	for _, submodule := range r.Modules {
		if err := h.renderModule(host, r, submodule.Path); err != nil {
			return err
		}
	}

	return nil
}

func (h *HandlebarsRenderder) renderModule(host string, r config.Repository, submodule string) error {
	moduleDir := filepath.Join(h.outputDir, r.Path, submodule)

	if err := os.MkdirAll(moduleDir, 0o755); err != nil { // nolint: gosec
		return fmt.Errorf("could not create repository directory %q: %w", moduleDir, err)
	}

	moduleFile := filepath.Join(moduleDir, indexFile)

	ctx := map[string]any{
		"host":       host,
		"library":    r.Library,
		"path":       r.Path,
		"repository": r.Repository,
		"ref":        r.Ref,
		"version":    r.Version,
		"deprecated": r.Deprecated,
		"submodule":  submodule,
	}

	result, err := h.repositoryTpl.Exec(ctx)
	if err != nil {
		path := r.Path
		if len(submodule) > 0 {
			path = filepath.Join(path, submodule)
		}

		return fmt.Errorf("could not render repository %q: %w", path, err)
	}

	if err := os.WriteFile(moduleFile, []byte(result), 0o644); err != nil { // nolint: gosec
		return fmt.Errorf("could not write repository file %q: %w", moduleFile, err)
	}

	return nil
}

// NewHandlebarsRenderder creates a new HandlebarsRenderder.
func NewHandlebarsRenderder(homepageSrc, repositorySrc, outputDir string) (*HandlebarsRenderder, error) {
	homepageTpl, err := raymond.Parse(homepageSrc)
	if err != nil {
		return nil, fmt.Errorf("could not parse homepage template: %w", err)
	}

	repositoryTpl, err := raymond.Parse(repositorySrc)
	if err != nil {
		return nil, fmt.Errorf("could not parse repository template: %w", err)
	}

	r := &HandlebarsRenderder{
		homepageTpl:   homepageTpl,
		repositoryTpl: repositoryTpl,
		outputDir:     outputDir,
	}

	return r, nil
}
