package sitefragment

import (
	"errors"

	"go.nhat.io/vanityrender/internal/module"
	"go.nhat.io/vanityrender/internal/service/sitecache"
	"go.nhat.io/vanityrender/internal/site"
)

var _ site.Hydrator = (*Hydrator)(nil)

// Hydrator is a site.Hydrator.
type Hydrator struct {
	cache    site.Hydrator
	upstream site.Hydrator
	modules  []string
}

// Hydrate hydrates the site configuration.
func (h *Hydrator) Hydrate(s *site.Site) error {
	if len(h.modules) == 0 {
		return h.upstream.Hydrate(s)
	}

	originalRepos := repositoriesMap(s.Repositories)

	if err := h.cache.Hydrate(s); err != nil {
		if isCacheErrors(err) {
			return h.upstream.Hydrate(s)
		}

		return err // nolint: errcheck
	}

	return h.hydrateFragments(s, originalRepos)
}

func (h *Hydrator) hydrateFragments(s *site.Site, originalRepos map[string]site.Repository) error {
	indexes := make(map[string]int, len(h.modules))
	for _, m := range h.modules {
		indexes[module.PathWithoutVersion(m)] = -1
	}

	s2 := *s
	s2.Repositories = make([]site.Repository, 0, len(h.modules))

	for i, r := range s.Repositories {
		path := module.PathWithoutVersion(r.Path)

		if _, ok := indexes[path]; ok {
			o := originalRepos[path]
			o.Path = path
			o.Modules = nil

			indexes[path] = i

			s2.Repositories = append(s2.Repositories, o)
		}
	}

	if err := h.upstream.Hydrate(&s2); err != nil {
		return err // nolint: errcheck
	}

	for _, r := range s2.Repositories {
		i := indexes[module.PathWithoutVersion(r.Path)]
		s.Repositories[i] = r
	}

	return nil
}

// NewHydrator initiates a new site.Hydrator.
func NewHydrator(cache site.Hydrator, upstream site.Hydrator, modules ...string) *Hydrator {
	return &Hydrator{
		cache:    cache,
		upstream: upstream,
		modules:  modules,
	}
}

func repositoriesMap(repos []site.Repository) map[string]site.Repository {
	m := make(map[string]site.Repository, len(repos))

	for _, r := range repos {
		m[module.PathWithoutVersion(r.Path)] = r
	}

	return m
}

func isCacheErrors(err error) bool {
	if errors.Is(err, sitecache.ErrChecksumMismatched) ||
		errors.Is(err, sitecache.ErrMetadataNotFound) ||
		errors.Is(err, sitecache.ErrMetadataInvalid) {
		return true
	}

	return false
}
