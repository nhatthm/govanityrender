package github_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.nhat.io/vanityrender/internal/github"
	"go.nhat.io/vanityrender/internal/module"
	"go.nhat.io/vanityrender/internal/site"
)

func TestHydrator_Hydrate(t *testing.T) {
	t.Parallel()

	pathVersions := map[module.Path]module.Version{
		".":          module.NewVersionFromString("v1.0.0"),
		"v2":         module.NewVersionFromString("v2.10.0"),
		"contrib":    module.NewVersionFromString("v0.2.0"),
		"contrib/v2": module.NewVersionFromString("v2.0.0"),
		"test":       module.NewVersionFromString("v0.2.0"),
	}

	expectedModules := []site.Module{
		{
			Path:          "repository",
			ImportPrefix:  "repository",
			VCS:           "git",
			RepositoryURL: "https://github.com/org/repository",
			HomeURL:       "https://github.com/org/repository",
			DirectoryURL:  "https://github.com/org/repository/tree/master{/dir}",
			FileURL:       "https://github.com/org/repository/blob/master{/dir}/{file}#L{line}",
		},
		{
			Path:          "repository/contrib",
			ImportPrefix:  "repository",
			VCS:           "git",
			RepositoryURL: "https://github.com/org/repository",
			HomeURL:       "https://github.com/org/repository",
			DirectoryURL:  "https://github.com/org/repository/tree/master{/dir}",
			FileURL:       "https://github.com/org/repository/blob/master{/dir}/{file}#L{line}",
		},
		{
			Path:          "repository/contrib/v2",
			ImportPrefix:  "repository",
			VCS:           "git",
			RepositoryURL: "https://github.com/org/repository",
			HomeURL:       "https://github.com/org/repository",
			DirectoryURL:  "https://github.com/org/repository/tree/master{/dir}",
			FileURL:       "https://github.com/org/repository/blob/master{/dir}/{file}#L{line}",
		},
		{
			Path:          "repository/test",
			ImportPrefix:  "repository",
			VCS:           "git",
			RepositoryURL: "https://github.com/org/repository",
			HomeURL:       "https://github.com/org/repository",
			DirectoryURL:  "https://github.com/org/repository/tree/master{/dir}",
			FileURL:       "https://github.com/org/repository/blob/master{/dir}/{file}#L{line}",
		},
		{
			Path:          "repository/v2",
			ImportPrefix:  "repository",
			VCS:           "git",
			RepositoryURL: "https://github.com/org/repository",
			HomeURL:       "https://github.com/org/repository",
			DirectoryURL:  "https://github.com/org/repository/tree/master{/dir}",
			FileURL:       "https://github.com/org/repository/blob/master{/dir}/{file}#L{line}",
		},
	}

	testCases := []struct {
		scenario       string
		moduleFinder   module.Finder
		site           site.Site
		expectedResult site.Site
		expectedError  string
	}{
		{
			scenario: "ignore non-github",
			site: site.Site{
				Repositories: []site.Repository{{
					RepositoryURL: "unknown",
				}},
			},
			expectedResult: site.Site{
				Repositories: []site.Repository{{
					RepositoryURL: "unknown",
				}},
			},
		},
		{
			scenario:     "error",
			moduleFinder: mockModuleFinderError(errors.New("find error")),
			site: site.Site{
				Repositories: []site.Repository{{
					RepositoryURL: "https://github.com/org/repository",
				}},
			},
			expectedResult: site.Site{
				Repositories: []site.Repository{{
					RepositoryURL: "https://github.com/org/repository",
				}},
			},
			expectedError: "find error",
		},
		{
			scenario:     "success - http",
			moduleFinder: mockModuleFinder(pathVersions),
			site: site.Site{
				Repositories: []site.Repository{{
					RepositoryURL: "http://github.com/org/repository",
					Path:          "repository",
				}},
			},
			expectedResult: site.Site{
				Repositories: []site.Repository{{
					RepositoryURL:  "https://github.com/org/repository",
					RepositoryName: "github.com/org/repository",
					Path:           "repository/v2",
					Modules:        expectedModules,
					LatestVersion:  "v2.10.0",
				}},
			},
		},
		{
			scenario:     "success - https",
			moduleFinder: mockModuleFinder(pathVersions),
			site: site.Site{
				Repositories: []site.Repository{{
					RepositoryURL: "https://github.com/org/repository",
					Path:          "repository",
				}},
			},
			expectedResult: site.Site{
				Repositories: []site.Repository{{
					RepositoryURL:  "https://github.com/org/repository",
					RepositoryName: "github.com/org/repository",
					Path:           "repository/v2",
					Modules:        expectedModules,
					LatestVersion:  "v2.10.0",
				}},
			},
		},
		{
			scenario:     "success - https with git",
			moduleFinder: mockModuleFinder(pathVersions),
			site: site.Site{
				Repositories: []site.Repository{{
					RepositoryURL: "https://github.com/org/repository.git",
					Path:          "repository",
				}},
			},
			expectedResult: site.Site{
				Repositories: []site.Repository{{
					RepositoryURL:  "https://github.com/org/repository",
					RepositoryName: "github.com/org/repository",
					Path:           "repository/v2",
					Modules:        expectedModules,
					LatestVersion:  "v2.10.0",
				}},
			},
		},
		{
			scenario:     "success - git",
			moduleFinder: mockModuleFinder(pathVersions),
			site: site.Site{
				Repositories: []site.Repository{{
					RepositoryURL:  "git@github.com:org/repository.git",
					RepositoryName: "github.com/org/repository",
					Path:           "repository",
				}},
			},
			expectedResult: site.Site{
				Repositories: []site.Repository{{
					RepositoryURL:  "https://github.com/org/repository",
					RepositoryName: "github.com/org/repository",
					Path:           "repository/v2",
					Modules:        expectedModules,
					LatestVersion:  "v2.10.0",
				}},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			err := github.NewHydrator(tc.moduleFinder).Hydrate(&tc.site)

			assert.Equal(t, tc.expectedResult, tc.site)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

type moduleFinderFunc func(loc, ref string) (map[module.Path]module.Version, error)

func (f moduleFinderFunc) Find(loc, ref string) (map[module.Path]module.Version, error) {
	return f(loc, ref)
}

func mockModuleFinderError(err error) moduleFinderFunc {
	return func(string, string) (map[module.Path]module.Version, error) {
		return nil, err
	}
}

func mockModuleFinder(versions map[module.Path]module.Version) moduleFinderFunc {
	return func(string, string) (map[module.Path]module.Version, error) {
		return versions, nil
	}
}
