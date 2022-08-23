package sitecache

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/fatih/color"

	"go.nhat.io/vanityrender/internal/site"
)

var _ site.Renderder = (*Renderder)(nil)

// Renderder is a site.Renderder.
type Renderder struct {
	upstream site.Renderder

	outputDir string
	checksum  string
	output    io.Writer
}

// Render renders the site.
func (r *Renderder) Render(s site.Site) error {
	if err := r.upstream.Render(s); err != nil {
		return err // nolint: errcheck
	}

	if err := r.renderMetadata(s); err != nil {
		return fmt.Errorf("could not render metadata: %w", err)
	}

	return nil
}

func (r *Renderder) renderMetadata(s site.Site) error {
	m := metadata{
		Checksum: r.checksum,
		Site:     s,
	}

	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return fmt.Errorf("could not marshal metadata: %w", err)
	}

	cacheFile := filepath.Join(r.outputDir, metadataFile)

	if err := os.WriteFile(cacheFile, data, 0o644); err != nil { // nolint: gosec
		return err
	}

	_, _ = fmt.Fprintln(r.output, color.HiGreenString("Render"), ":", metadataFile)

	return nil
}

// NewRenderder initiates a new site.Renderder.
func NewRenderder(upstream site.Renderder, outputDir string, checksum string, opts ...RendererOption) *Renderder {
	r := &Renderder{
		upstream:  upstream,
		outputDir: outputDir,
		checksum:  checksum,
		output:    io.Discard,
	}

	for _, o := range opts {
		o.applyRendererOption(r)
	}

	return r
}

// RendererOption is an option to configure renderer.
type RendererOption interface {
	applyRendererOption(r *Renderder)
}

type rendererOptionFunc func(r *Renderder)

func (f rendererOptionFunc) applyRendererOption(r *Renderder) {
	f(r)
}
