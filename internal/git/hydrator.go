package git

import (
	"fmt"
	"io/fs"
	"path/filepath"

	"go.nhat.io/vanityrender/internal/config"
)

var goMod = `go.mod`

var _ config.Hydrator = (*Hydrator)(nil)

// Hydrator is a config.Hydrator.
type Hydrator struct{}

// Hydrate hydrates the configuration.
func (h *Hydrator) Hydrate(cfg *config.Config) error {
	for i := range cfg.Repositories {
		if err := h.hydrateRepository(&cfg.Repositories[i]); err != nil {
			return err
		}
	}

	return nil
}

func (h *Hydrator) hydrateRepository(r *config.Repository) error {
	dir, err := clone(r.Repository, r.Ref)
	if err != nil {
		return err
	}

	modules := make([]config.Module, 0)

	err = filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && info.Name() == goMod {
			rel, err := filepath.Rel(dir, path)
			if err != nil {
				return fmt.Errorf("could not get relative path: %w", err)
			}

			if path = filepath.Dir(rel); path != "." {
				modules = append(modules, config.Module{
					Path: path,
				})
			}
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("could not walk directory: %w", err)
	}

	r.Modules = modules

	return nil
}

// NewHydrator initiates a new config.Hydrator.
func NewHydrator() *Hydrator {
	return &Hydrator{}
}
