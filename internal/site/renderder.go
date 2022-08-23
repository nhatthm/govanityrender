package site

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/aymerick/raymond"
)

const indexFile = `index.html`

// Renderder is the interface for rendering.
type Renderder interface {
	Render(s Site) error
}

var _ Renderder = (*HandlebarsRenderder)(nil)

// HandlebarsRenderder renders to the filesystem.
type HandlebarsRenderder struct {
	homepageTpl   *raymond.Template
	repositoryTpl *raymond.Template
	outputDir     string
}

// Render renders the configuration.
func (h *HandlebarsRenderder) Render(s Site) error {
	if err := h.renderHomepage(s); err != nil {
		return fmt.Errorf("could not render homepage: %w", err)
	}

	for _, r := range s.Repositories {
		if err := h.renderRepository(s.Hostname, r); err != nil {
			return err
		}
	}

	return nil
}

func (h *HandlebarsRenderder) renderHomepage(s Site) error {
	homepageFile := filepath.Join(h.outputDir, indexFile)

	repositories := make([]map[string]any, len(s.Repositories))
	for i, r := range s.Repositories {
		repositories[i] = map[string]any{
			"name":           r.Name,
			"path":           r.Path,
			"deprecated":     r.Deprecated,
			"repositoryURL":  r.RepositoryURL,
			"repositoryName": r.RepositoryName,
			"latestVersion":  r.LatestVersion,
		}
	}

	inputs := map[string]any{
		"pageTitle":       s.PageTitle,
		"pageDescription": s.PageDescription,
		"host":            s.Hostname,
		"repositories":    repositories,
	}

	result, err := h.homepageTpl.Exec(inputs)
	if err != nil {
		return err
	}

	return os.WriteFile(homepageFile, []byte(result), 0o644) // nolint: gosec
}

func (h *HandlebarsRenderder) renderRepository(host string, r Repository) error {
	for _, m := range r.Modules {
		if err := h.renderModule(host, m); err != nil {
			return err
		}
	}

	return nil
}

func (h *HandlebarsRenderder) renderModule(host string, m Module) error {
	moduleDir := filepath.Join(h.outputDir, m.Path)

	if err := os.MkdirAll(moduleDir, 0o755); err != nil { // nolint: gosec
		return fmt.Errorf("could not create repository directory %q: %w", moduleDir, err)
	}

	moduleFile := filepath.Join(moduleDir, indexFile)

	ctx := map[string]any{
		"host":          host,
		"path":          m.Path,
		"importPrefix":  m.ImportPrefix,
		"vcs":           m.VCS,
		"repositoryURL": m.RepositoryURL,
		"homeURL":       m.HomeURL,
		"directoryURL":  m.DirectoryURL,
		"fileURL":       m.FileURL,
	}

	result, err := h.repositoryTpl.Exec(ctx)
	if err != nil {
		return fmt.Errorf("could not render repository %q: %w", m.ImportPrefix, err)
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
