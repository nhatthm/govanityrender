package git

import (
	"sort"

	"go.nhat.io/vanityrender/internal/module"
)

// ModuleFinder finds modules in a repository.
type ModuleFinder struct{}

// Find finds modules in a repository.
func (f *ModuleFinder) Find(loc, ref string) (map[module.Path]module.Version, error) {
	dir, r, err := Clone(loc, ref)
	if err != nil {
		return nil, err
	}

	versions, err := Versions(r)
	if err != nil {
		return nil, err
	}

	goModVersions, err := module.FindVersions(dir)
	if err != nil {
		return nil, err // nolint: wrapcheck
	}

	versions = append(versions, goModVersions...)
	sort.Strings(versions)

	result := make(map[module.Path]module.Version, len(versions))

	for _, s := range versions {
		k, v := module.PathVersion(s)

		result[k] = v
	}

	return result, nil
}

// NewModuleFinder returns a new module finder.
func NewModuleFinder() *ModuleFinder {
	return &ModuleFinder{}
}
