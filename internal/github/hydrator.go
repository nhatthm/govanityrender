package github

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"go.nhat.io/vanityrender/internal/module"
	"go.nhat.io/vanityrender/internal/site"
)

const gitHubDomain = `https://github.com/`

var repositoryURLSanitizer = strings.NewReplacer(
	"https://", "",
	"http://", "",
	".git", "",
	"git@", "",
	"github.com:", "github.com/",
)

var _ site.Hydrator = (*Hydrator)(nil)

// Hydrator is a config.Hydrator.
type Hydrator struct {
	finder module.Finder
}

// Hydrate hydrates the configuration.
func (h *Hydrator) Hydrate(s *site.Site) error {
	for i := range s.Repositories {
		if err := h.hydrateRepository(&s.Repositories[i]); err != nil {
			return err
		}
	}

	return nil
}

func (h *Hydrator) hydrateRepository(r *site.Repository) error {
	repoURL := repositoryURL(r.RepositoryURL)

	if !strings.Contains(repoURL, gitHubDomain) {
		return nil
	}

	pathVersions, err := h.finder.Find(repoURL, r.Ref)
	if err != nil {
		return err // nolint: wrapcheck
	}

	r.RepositoryURL = repoURL
	r.RepositoryName = repositoryURLSanitizer.Replace(repoURL)

	modules := make([]site.Module, 0, len(pathVersions))
	latestVersion := module.Version{}

	for path, version := range pathVersions {
		modulePath := r.Path
		if string(path) != "." {
			modulePath = filepath.Join(r.Path, string(path))
		}

		modules = append(modules, site.Module{
			Path:          modulePath,
			ImportPrefix:  r.Path,
			VCS:           "git",
			RepositoryURL: r.RepositoryURL,
			HomeURL:       r.RepositoryURL,
			DirectoryURL:  fmt.Sprintf("%s/tree/master{/dir}", r.RepositoryURL),
			FileURL:       fmt.Sprintf("%s/blob/master{/dir}/{file}#L{line}", r.RepositoryURL),
		})

		if path.IsRoot() && latestVersion.LessThan(version) {
			latestVersion = version
		}
	}

	sort.Slice(modules, func(i, j int) bool {
		return modules[i].Path < modules[j].Path
	})

	if latestVersion.Major > 1 {
		r.Path = fmt.Sprintf("%s/v%d", r.Path, latestVersion.Major)
	}

	r.LatestVersion = latestVersion.String()
	r.Modules = modules

	return nil
}

// NewHydrator initiates a new config.Hydrator.
func NewHydrator(finder module.Finder) *Hydrator {
	return &Hydrator{
		finder: finder,
	}
}

func repositoryURL(url string) string {
	result := repositoryURLSanitizer.Replace(url)

	return fmt.Sprintf("https://%s", strings.TrimRight(result, "/"))
}
