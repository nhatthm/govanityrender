package git

import (
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

	versions := []string{"v0.0.0"}

	taggedVersions, err := Versions(r)
	if err != nil {
		return nil, err
	}

	versions = append(versions, taggedVersions...)

	goModVersions, err := module.FindVersions(dir)
	if err != nil {
		return nil, err // nolint: wrapcheck
	}

	versions = append(versions, goModVersions...)

	result := make(map[module.Path]module.Version, len(versions))

	for _, s := range versions {
		k, v := module.PathVersion(s)

		if curVersion, ok := result[k]; !ok || curVersion.LessThan(v) {
			result[k] = v
		}
	}

	return result, nil
}

// NewModuleFinder returns a new module finder.
func NewModuleFinder() *ModuleFinder {
	return &ModuleFinder{}
}
