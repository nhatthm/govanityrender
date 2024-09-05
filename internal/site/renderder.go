package site

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/aymerick/raymond"
	"github.com/fatih/color"

	"go.nhat.io/vanityrender/internal/version"
)

const (
	indexHTMLFile    = `index.html`
	notFoundHTMLFile = `404.html`
)

// Renderder is the interface for rendering.
type Renderder interface {
	Render(s Site) error
}

var _ Renderder = (*HandlebarsRenderder)(nil)

// HandlebarsRenderder renders to the filesystem.
type HandlebarsRenderder struct {
	homepageTpl   *raymond.Template
	notFoundTpl   *raymond.Template
	repositoryTpl *raymond.Template
	outputDir     string

	output io.Writer
}

// Render renders the configuration.
func (h *HandlebarsRenderder) Render(s Site) error {
	if err := h.renderHomepage(s); err != nil {
		return fmt.Errorf("could not render homepage: %w", err)
	}

	if err := h.render404(s); err != nil {
		return fmt.Errorf("could not render 404: %w", err)
	}

	for _, r := range s.Repositories {
		if err := h.renderRepository(s.Hostname, r); err != nil {
			return err
		}
	}

	return nil
}

func (h *HandlebarsRenderder) renderHomepage(s Site) error {
	homepageFile := filepath.Join(h.outputDir, indexHTMLFile)

	repositories := make([]map[string]any, len(s.Repositories))
	for i, r := range s.Repositories {
		repositories[i] = map[string]any{
			"name":           r.Name,
			"path":           r.Path,
			"deprecated":     r.Deprecated,
			"hidden":         r.Hidden,
			"repositoryURL":  r.RepositoryURL,
			"repositoryName": r.RepositoryName,
			"latestVersion":  r.LatestVersion,
		}
	}

	inputs := map[string]any{
		"pageTitle":       s.PageTitle,
		"pageDescription": s.PageDescription,
		"host":            s.Hostname,
		"sourceURL":       s.SourceURL,
		"repositories":    repositories,
		"renderer":        version.Info(),
	}

	result, err := h.homepageTpl.Exec(inputs)
	if err != nil {
		return err
	}

	if err := os.WriteFile(homepageFile, []byte(result), 0o644); err != nil { // nolint: gosec
		return err
	}

	_, _ = fmt.Fprintln(h.output, color.HiGreenString("Render"), ":", indexHTMLFile) //nolint: errcheck

	return nil
}

func (h *HandlebarsRenderder) render404(s Site) error {
	notFoundFile := filepath.Join(h.outputDir, notFoundHTMLFile)

	inputs := map[string]any{
		"pageTitle":       s.PageTitle,
		"pageDescription": s.PageDescription,
		"host":            s.Hostname,
		"sourceURL":       s.SourceURL,
		"renderer":        version.Info(),
	}

	result, err := h.notFoundTpl.Exec(inputs)
	if err != nil {
		return err
	}

	if err := os.WriteFile(notFoundFile, []byte(result), 0o644); err != nil { // nolint: gosec
		return err
	}

	_, _ = fmt.Fprintln(h.output, color.HiGreenString("Render"), ":", notFoundHTMLFile) //nolint: errcheck

	return nil
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

	moduleFile := filepath.Join(moduleDir, indexHTMLFile)

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

	_, _ = fmt.Fprintln(h.output, color.HiGreenString("Render"), ":", filepath.Join(m.Path, indexHTMLFile)) //nolint: errcheck

	return nil
}

// NewHandlebarsRenderder creates a new HandlebarsRenderder.
func NewHandlebarsRenderder(
	homepageSrc, notFoundSrc, repositorySrc string,
	outputDir string,
	opts ...RendererOption,
) (*HandlebarsRenderder, error) {
	homepageTpl, err := raymond.Parse(homepageSrc)
	if err != nil {
		return nil, fmt.Errorf("could not parse homepage template: %w", err)
	}

	notFoundTpl, err := raymond.Parse(notFoundSrc)
	if err != nil {
		return nil, fmt.Errorf("could not parse 404 template: %w", err)
	}

	repositoryTpl, err := raymond.Parse(repositorySrc)
	if err != nil {
		return nil, fmt.Errorf("could not parse repository template: %w", err)
	}

	r := &HandlebarsRenderder{
		homepageTpl:   homepageTpl,
		notFoundTpl:   notFoundTpl,
		repositoryTpl: repositoryTpl,
		outputDir:     outputDir,
		output:        io.Discard,
	}

	for _, o := range opts {
		o.applyRendererOption(r)
	}

	return r, nil
}

// RendererOption is an option to configure renderer.
type RendererOption interface {
	applyRendererOption(r *HandlebarsRenderder)
}

type rendererOptionFunc func(r *HandlebarsRenderder)

func (f rendererOptionFunc) applyRendererOption(r *HandlebarsRenderder) {
	f(r)
}

// WithOutput sets the output writer.
func WithOutput(w io.Writer) RendererOption {
	return rendererOptionFunc(func(r *HandlebarsRenderder) {
		r.output = w
	})
}
