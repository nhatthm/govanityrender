package sitecache

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"go.nhat.io/vanityrender/internal/site"
)

var _ site.Renderder = (*Renderder)(nil)

// Renderder is a site.Renderder.
type Renderder struct {
	upstream site.Renderder

	outputDir string
	checksum  string
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

	metadataFile := filepath.Join(r.outputDir, metadataFile)

	return os.WriteFile(metadataFile, data, 0o644) // nolint: gosec
}

// NewRenderder initiates a new site.Renderder.
func NewRenderder(upstream site.Renderder, outputDir string, checksum string) *Renderder {
	return &Renderder{
		upstream:  upstream,
		outputDir: outputDir,
		checksum:  checksum,
	}
}
