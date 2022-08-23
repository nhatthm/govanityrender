package github

import (
	"context"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"go.nhat.io/vanityrender/internal/module"
	"go.nhat.io/vanityrender/internal/site"
)

const (
	gitHubDomain = `https://github.com/`

	defaultNumWorkers = 5
)

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

	numWorkers int
}

// Hydrate hydrates the configuration.
func (h *Hydrator) Hydrate(s *site.Site) error {
	ch := make(chan *site.Repository, h.numWorkers*2)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg := sync.WaitGroup{}
	errMu := sync.Mutex{}
	err := (error)(nil)

	wg.Add(h.numWorkers)

	for i := 0; i < h.numWorkers; i++ {
		go func() {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					return

				case r, ok := <-ch:
					if !ok {
						return
					}

					if hErr := h.hydrateRepository(r); hErr != nil {
						errMu.Lock()
						err = hErr
						errMu.Unlock()

						cancel()
					}
				}
			}
		}()
	}

	go func() {
		defer close(ch)

		for i := range s.Repositories {
			ch <- &s.Repositories[i]
		}
	}()

	wg.Wait()

	return err
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
		r.Path = module.PathWithVersion(r.Path, latestVersion)
	}

	r.LatestVersion = latestVersion.String()
	r.Modules = modules

	return nil
}

// NewHydrator initiates a new config.Hydrator.
func NewHydrator(finder module.Finder) *Hydrator {
	return &Hydrator{
		finder:     finder,
		numWorkers: defaultNumWorkers,
	}
}

func repositoryURL(url string) string {
	result := repositoryURLSanitizer.Replace(url)

	return fmt.Sprintf("https://%s", strings.TrimRight(result, "/"))
}
