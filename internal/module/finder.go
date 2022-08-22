package module

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"

	"golang.org/x/mod/modfile"
)

var fileGoModVersionRegExp = regexp.MustCompile(`/(v\d+)$`)

// Finder finds modules.
type Finder interface {
	Find(loc, ref string) (map[Path]Version, error)
}

// FindVersions returns the module versions in the given path.
func FindVersions(dir string) ([]string, error) {
	var result []string

	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && info.Name() == goMod {
			f, err := parseGoMod(path)
			if err != nil {
				return err
			}

			rel, err := filepath.Rel(dir, path)
			if err != nil {
				return fmt.Errorf("could not get relative path: %w", err)
			}

			modulePath := filepath.Dir(rel)
			version := "v0.0.0"

			if f.Module != nil && fileGoModVersionRegExp.MatchString(f.Module.Mod.Path) {
				m := fileGoModVersionRegExp.FindStringSubmatch(f.Module.Mod.Path)
				version = fmt.Sprintf("%s.0.0", m[1])
			}

			pathVersion := version
			if modulePath != "." {
				pathVersion = fmt.Sprintf("%s/%s", modulePath, version)
			}

			result = append(result, pathVersion)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("could not walk directory: %w", err)
	}

	sort.Strings(result)

	return result, nil
}

func parseGoMod(file string) (*modfile.File, error) {
	data, err := os.ReadFile(filepath.Clean(file))
	if err != nil {
		return nil, fmt.Errorf("could not read file %q: %w", file, err)
	}

	f, err := modfile.Parse(file, data, nil)
	if err != nil {
		panic(fmt.Errorf("could not parse mod file: %w", err))
	}

	return f, nil
}
