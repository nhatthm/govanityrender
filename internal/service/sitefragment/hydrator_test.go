package sitefragment_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.nhat.io/vanityrender/internal/service/sitecache"
	"go.nhat.io/vanityrender/internal/service/sitefragment"
	"go.nhat.io/vanityrender/internal/site"
)

func TestHydrator_Hydrate(t *testing.T) { // nolint: maintidx
	t.Parallel()

	testCases := []struct {
		scenario       string
		mockCache      func(t *testing.T) site.Hydrator
		mockUpstream   func(t *testing.T) site.Hydrator
		input          site.Site
		modules        []string
		expectedResult site.Site
		expectedError  string
	}{
		{
			scenario:      "no modules - upstream error",
			mockUpstream:  mockHydrateError(errors.New("upstream error")),
			mockCache:     nopHydrator,
			expectedError: "upstream error",
		},
		{
			scenario: "no modules - upstream success",
			mockUpstream: mockHydrateSuccess(func(_ *testing.T, s *site.Site) {
				s.PageTitle = "test"
			}),
			mockCache:      nopHydrator,
			expectedResult: site.Site{PageTitle: "test"},
		},
		{
			scenario:      "has modules - cache error",
			mockCache:     mockHydrateError(errors.New("cache error")),
			mockUpstream:  nopHydrator,
			modules:       []string{"contrib"},
			expectedError: "cache error",
		},
		{
			scenario:  "has modules - cache error - checksum mismatched",
			mockCache: mockHydrateError(sitecache.ErrChecksumMismatched),
			mockUpstream: mockHydrateSuccess(func(_ *testing.T, s *site.Site) {
				s.PageTitle = "test"
			}),
			modules:        []string{"contrib"},
			expectedResult: site.Site{PageTitle: "test"},
		},
		{
			scenario:  "has modules - cache error - not found",
			mockCache: mockHydrateError(sitecache.ErrMetadataNotFound),
			mockUpstream: mockHydrateSuccess(func(_ *testing.T, s *site.Site) {
				s.PageTitle = "test"
			}),
			modules:        []string{"contrib"},
			expectedResult: site.Site{PageTitle: "test"},
		},
		{
			scenario:  "has modules - cache error - invalid metadata",
			mockCache: mockHydrateError(sitecache.ErrMetadataInvalid),
			mockUpstream: mockHydrateSuccess(func(_ *testing.T, s *site.Site) {
				s.PageTitle = "test"
			}),
			modules:        []string{"contrib"},
			expectedResult: site.Site{PageTitle: "test"},
		},
		{
			scenario:      "has modules - upstream error",
			mockCache:     mockHydrateSuccess(),
			mockUpstream:  mockHydrateError(errors.New("upstream error")),
			modules:       []string{"contrib"},
			expectedError: "upstream error",
		},
		{
			scenario: "has modules - success",
			mockCache: mockHydrateSuccess(func(_ *testing.T, s *site.Site) {
				s.Repositories = []site.Repository{
					{
						Name:           "Test",
						Path:           "test",
						RepositoryURL:  "https://github.com/org/go-test",
						RepositoryName: "github.com/org/go-test",
						LatestVersion:  "v0.6.0",
						Modules: []site.Module{
							{
								Path:          "test",
								ImportPrefix:  "test",
								VCS:           "git",
								RepositoryURL: "https://github.com/org/go-test",
								HomeURL:       "https://github.com/org/go-test",
								DirectoryURL:  "https://github.com/org/go-test/tree/master{/dir}",
								FileURL:       "https://github.com/org/go-test/blob/master{/dir}/{file}#L{line}",
							},
						},
					},
					{
						Name:           "Contrib",
						Path:           "contrib/v2",
						RepositoryURL:  "https://github.com/org/go-contrib",
						RepositoryName: "github.com/org/go-contrib",
						LatestVersion:  "v2.3.2",
						Modules: []site.Module{
							{
								Path:          "contrib",
								ImportPrefix:  "contrib",
								VCS:           "git",
								RepositoryURL: "https://github.com/org/go-contrib",
								HomeURL:       "https://github.com/org/go-contrib",
								DirectoryURL:  "https://github.com/org/go-contrib/tree/master{/dir}",
								FileURL:       "https://github.com/org/go-contrib/blob/master{/dir}/{file}#L{line}",
							},
							{
								Path:          "contrib/v2",
								ImportPrefix:  "contrib",
								VCS:           "git",
								RepositoryURL: "https://github.com/org/go-contrib",
								HomeURL:       "https://github.com/org/go-contrib",
								DirectoryURL:  "https://github.com/org/go-contrib/tree/master{/dir}",
								FileURL:       "https://github.com/org/go-contrib/blob/master{/dir}/{file}#L{line}",
							},
							{
								Path:          "contrib/module",
								ImportPrefix:  "contrib",
								VCS:           "git",
								RepositoryURL: "https://github.com/org/go-contrib",
								HomeURL:       "https://github.com/org/go-contrib",
								DirectoryURL:  "https://github.com/org/go-contrib/tree/master{/dir}",
								FileURL:       "https://github.com/org/go-contrib/blob/master{/dir}/{file}#L{line}",
							},
							{
								Path:          "contrib/module/v2",
								ImportPrefix:  "contrib",
								VCS:           "git",
								RepositoryURL: "https://github.com/org/go-contrib",
								HomeURL:       "https://github.com/org/go-contrib",
								DirectoryURL:  "https://github.com/org/go-contrib/tree/master{/dir}",
								FileURL:       "https://github.com/org/go-contrib/blob/master{/dir}/{file}#L{line}",
							},
						},
					},
					{
						Name:           "Mock",
						Path:           "mock",
						RepositoryURL:  "https://github.com/org/go-mock",
						RepositoryName: "github.com/org/go-mock",
						LatestVersion:  "v0.3.0",
						Modules: []site.Module{
							{
								Path:          "mock",
								ImportPrefix:  "mock",
								VCS:           "git",
								RepositoryURL: "https://github.com/org/go-mock",
								HomeURL:       "https://github.com/org/go-mock",
								DirectoryURL:  "https://github.com/org/go-mock/tree/master{/dir}",
								FileURL:       "https://github.com/org/go-mock/blob/master{/dir}/{file}#L{line}",
							},
						},
					},
					{
						Name:           "Experiment",
						Path:           "experiment",
						RepositoryURL:  "https://github.com/org/go-experiment",
						RepositoryName: "github.com/org/go-experiment",
						LatestVersion:  "v1.5.0",
						Modules: []site.Module{
							{
								Path:          "experiment",
								ImportPrefix:  "experiment",
								VCS:           "git",
								RepositoryURL: "https://github.com/org/go-experiment",
								HomeURL:       "https://github.com/org/go-experiment",
								DirectoryURL:  "https://github.com/org/go-experiment/tree/master{/dir}",
								FileURL:       "https://github.com/org/go-experiment/blob/master{/dir}/{file}#L{line}",
							},
						},
					},
				}
			}),
			mockUpstream: mockHydrateSuccess(func(t *testing.T, s *site.Site) {
				t.Helper()

				expected := []site.Repository{
					{
						Name:          "Contrib",
						Path:          "contrib",
						RepositoryURL: "https://github.com/org/go-contrib",
					},
					{
						Name:          "Mock",
						Path:          "mock",
						RepositoryURL: "https://github.com/org/go-mock",
					},
				}

				assert.Equal(t, expected, s.Repositories)

				s.Repositories = []site.Repository{
					{
						Name:           "Contrib",
						Path:           "contrib/v3",
						RepositoryURL:  "https://github.com/org/go-contrib",
						RepositoryName: "github.com/org/go-contrib",
						LatestVersion:  "v3.0.0",
						Modules: []site.Module{
							{
								Path:          "contrib",
								ImportPrefix:  "contrib",
								VCS:           "git",
								RepositoryURL: "https://github.com/org/go-contrib",
								HomeURL:       "https://github.com/org/go-contrib",
								DirectoryURL:  "https://github.com/org/go-contrib/tree/master{/dir}",
								FileURL:       "https://github.com/org/go-contrib/blob/master{/dir}/{file}#L{line}",
							},
							{
								Path:          "contrib/v2",
								ImportPrefix:  "contrib",
								VCS:           "git",
								RepositoryURL: "https://github.com/org/go-contrib",
								HomeURL:       "https://github.com/org/go-contrib",
								DirectoryURL:  "https://github.com/org/go-contrib/tree/master{/dir}",
								FileURL:       "https://github.com/org/go-contrib/blob/master{/dir}/{file}#L{line}",
							},
							{
								Path:          "contrib/v3",
								ImportPrefix:  "contrib",
								VCS:           "git",
								RepositoryURL: "https://github.com/org/go-contrib",
								HomeURL:       "https://github.com/org/go-contrib",
								DirectoryURL:  "https://github.com/org/go-contrib/tree/master{/dir}",
								FileURL:       "https://github.com/org/go-contrib/blob/master{/dir}/{file}#L{line}",
							},
							{
								Path:          "contrib/module",
								ImportPrefix:  "contrib",
								VCS:           "git",
								RepositoryURL: "https://github.com/org/go-contrib",
								HomeURL:       "https://github.com/org/go-contrib",
								DirectoryURL:  "https://github.com/org/go-contrib/tree/master{/dir}",
								FileURL:       "https://github.com/org/go-contrib/blob/master{/dir}/{file}#L{line}",
							},
							{
								Path:          "contrib/module/v2",
								ImportPrefix:  "contrib",
								VCS:           "git",
								RepositoryURL: "https://github.com/org/go-contrib",
								HomeURL:       "https://github.com/org/go-contrib",
								DirectoryURL:  "https://github.com/org/go-contrib/tree/master{/dir}",
								FileURL:       "https://github.com/org/go-contrib/blob/master{/dir}/{file}#L{line}",
							},
						},
					},
					{
						Name:           "Mock",
						Path:           "mock",
						RepositoryURL:  "https://github.com/org/go-mock",
						RepositoryName: "github.com/org/go-mock",
						LatestVersion:  "v0.5.0",
						Modules: []site.Module{
							{
								Path:          "mock",
								ImportPrefix:  "mock",
								VCS:           "git",
								RepositoryURL: "https://github.com/org/go-mock",
								HomeURL:       "https://github.com/org/go-mock",
								DirectoryURL:  "https://github.com/org/go-mock/tree/master{/dir}",
								FileURL:       "https://github.com/org/go-mock/blob/master{/dir}/{file}#L{line}",
							},
						},
					},
				}
			}),
			modules: []string{"contrib", "mock"},
			input: site.Site{
				Repositories: []site.Repository{
					{
						Name:          "Test",
						Path:          "test",
						RepositoryURL: "https://github.com/org/go-test",
					},
					{
						Name:          "Contrib",
						Path:          "contrib",
						RepositoryURL: "https://github.com/org/go-contrib",
					},
					{
						Name:          "Mock",
						Path:          "mock",
						RepositoryURL: "https://github.com/org/go-mock",
					},
					{
						Name:          "Experiment",
						Path:          "experiment",
						RepositoryURL: "https://github.com/org/go-experiment",
					},
				},
			},
			expectedResult: site.Site{
				Repositories: []site.Repository{
					{
						Name:           "Test",
						Path:           "test",
						RepositoryURL:  "https://github.com/org/go-test",
						RepositoryName: "github.com/org/go-test",
						LatestVersion:  "v0.6.0",
						Modules: []site.Module{
							{
								Path:          "test",
								ImportPrefix:  "test",
								VCS:           "git",
								RepositoryURL: "https://github.com/org/go-test",
								HomeURL:       "https://github.com/org/go-test",
								DirectoryURL:  "https://github.com/org/go-test/tree/master{/dir}",
								FileURL:       "https://github.com/org/go-test/blob/master{/dir}/{file}#L{line}",
							},
						},
					},
					{
						Name:           "Contrib",
						Path:           "contrib/v3",
						RepositoryURL:  "https://github.com/org/go-contrib",
						RepositoryName: "github.com/org/go-contrib",
						LatestVersion:  "v3.0.0",
						Modules: []site.Module{
							{
								Path:          "contrib",
								ImportPrefix:  "contrib",
								VCS:           "git",
								RepositoryURL: "https://github.com/org/go-contrib",
								HomeURL:       "https://github.com/org/go-contrib",
								DirectoryURL:  "https://github.com/org/go-contrib/tree/master{/dir}",
								FileURL:       "https://github.com/org/go-contrib/blob/master{/dir}/{file}#L{line}",
							},
							{
								Path:          "contrib/v2",
								ImportPrefix:  "contrib",
								VCS:           "git",
								RepositoryURL: "https://github.com/org/go-contrib",
								HomeURL:       "https://github.com/org/go-contrib",
								DirectoryURL:  "https://github.com/org/go-contrib/tree/master{/dir}",
								FileURL:       "https://github.com/org/go-contrib/blob/master{/dir}/{file}#L{line}",
							},
							{
								Path:          "contrib/v3",
								ImportPrefix:  "contrib",
								VCS:           "git",
								RepositoryURL: "https://github.com/org/go-contrib",
								HomeURL:       "https://github.com/org/go-contrib",
								DirectoryURL:  "https://github.com/org/go-contrib/tree/master{/dir}",
								FileURL:       "https://github.com/org/go-contrib/blob/master{/dir}/{file}#L{line}",
							},
							{
								Path:          "contrib/module",
								ImportPrefix:  "contrib",
								VCS:           "git",
								RepositoryURL: "https://github.com/org/go-contrib",
								HomeURL:       "https://github.com/org/go-contrib",
								DirectoryURL:  "https://github.com/org/go-contrib/tree/master{/dir}",
								FileURL:       "https://github.com/org/go-contrib/blob/master{/dir}/{file}#L{line}",
							},
							{
								Path:          "contrib/module/v2",
								ImportPrefix:  "contrib",
								VCS:           "git",
								RepositoryURL: "https://github.com/org/go-contrib",
								HomeURL:       "https://github.com/org/go-contrib",
								DirectoryURL:  "https://github.com/org/go-contrib/tree/master{/dir}",
								FileURL:       "https://github.com/org/go-contrib/blob/master{/dir}/{file}#L{line}",
							},
						},
					},
					{
						Name:           "Mock",
						Path:           "mock",
						RepositoryURL:  "https://github.com/org/go-mock",
						RepositoryName: "github.com/org/go-mock",
						LatestVersion:  "v0.5.0",
						Modules: []site.Module{
							{
								Path:          "mock",
								ImportPrefix:  "mock",
								VCS:           "git",
								RepositoryURL: "https://github.com/org/go-mock",
								HomeURL:       "https://github.com/org/go-mock",
								DirectoryURL:  "https://github.com/org/go-mock/tree/master{/dir}",
								FileURL:       "https://github.com/org/go-mock/blob/master{/dir}/{file}#L{line}",
							},
						},
					},
					{
						Name:           "Experiment",
						Path:           "experiment",
						RepositoryURL:  "https://github.com/org/go-experiment",
						RepositoryName: "github.com/org/go-experiment",
						LatestVersion:  "v1.5.0",
						Modules: []site.Module{
							{
								Path:          "experiment",
								ImportPrefix:  "experiment",
								VCS:           "git",
								RepositoryURL: "https://github.com/org/go-experiment",
								HomeURL:       "https://github.com/org/go-experiment",
								DirectoryURL:  "https://github.com/org/go-experiment/tree/master{/dir}",
								FileURL:       "https://github.com/org/go-experiment/blob/master{/dir}/{file}#L{line}",
							},
						},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			h := sitefragment.NewHydrator(tc.mockCache(t), tc.mockUpstream(t), tc.modules)

			err := h.Hydrate(&tc.input)

			assert.Equal(t, tc.expectedResult, tc.input)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

var nopHydrator = func(t *testing.T) site.Hydrator {
	t.Helper()

	return hydrateFunc(func(s *site.Site) error {
		return nil
	})
}

type hydrateFunc func(s *site.Site) error

func (f hydrateFunc) Hydrate(s *site.Site) error {
	return f(s)
}

func mockHydrateError(err error) func(t *testing.T) site.Hydrator {
	return func(t *testing.T) site.Hydrator {
		t.Helper()

		return hydrateFunc(func(s *site.Site) error {
			return err
		})
	}
}

func mockHydrateSuccess(mocks ...func(t *testing.T, s *site.Site)) func(t *testing.T) site.Hydrator {
	return func(t *testing.T) site.Hydrator {
		t.Helper()

		return hydrateFunc(func(s *site.Site) error {
			for _, m := range mocks {
				m(t, s)
			}

			return nil
		})
	}
}
