package site_test

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.nhat.io/vanityrender/internal/site"
	"go.nhat.io/vanityrender/templates"
)

func TestNewHandlebarsRenderder(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario      string
		homepageSrc   string
		NotFoundSrc   string
		repositorySrc string
		expectedError string
	}{
		{
			scenario:      "homepage template is broken",
			homepageSrc:   `{{`,
			expectedError: "could not parse homepage template: Parse error on line 1:",
		},
		{
			scenario:      "not found template is broken",
			homepageSrc:   `{{ message }}`,
			NotFoundSrc:   `{{`,
			expectedError: "could not parse 404 template: Parse error on line 1:",
		},
		{
			scenario:      "repository template is broken",
			homepageSrc:   `{{ message }}`,
			NotFoundSrc:   `{{ message }}`,
			repositorySrc: `{{`,
			expectedError: "could not parse repository template: Parse error on line 1:",
		},
		{
			scenario:      "success",
			homepageSrc:   `{{ message }}`,
			NotFoundSrc:   `{{ message }}`,
			repositorySrc: `{{ message }}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			actual, err := site.NewHandlebarsRenderder(tc.homepageSrc, tc.NotFoundSrc, tc.repositorySrc, "")

			if tc.expectedError == "" {
				assert.NotNil(t, actual)
				assert.NoError(t, err)
			} else {
				assert.Nil(t, actual)
				assert.ErrorContains(t, err, tc.expectedError)
			}
		})
	}
}

func TestHandlebarsRenderder_Render(t *testing.T) {
	t.Parallel()

	outputDir := t.TempDir()

	s := site.Site{
		PageTitle:       "go.nhat.io",
		PageDescription: "Open Source Go Modules",
		Hostname:        "go.nhat.io",
		SourceURL:       "github.com/nhatthm/govanityrender",
		Repositories: []site.Repository{
			{
				Name:           "Hidden",
				Path:           "hidden",
				RepositoryURL:  "https://github.com/nhatthm/hidden",
				RepositoryName: "github.com/nhatthm/hidden",
				Hidden:         true,
				LatestVersion:  "v0.1.0",
				Modules: []site.Module{{
					Path:          "hidden",
					ImportPrefix:  "hidden",
					VCS:           "git",
					RepositoryURL: "https://github.com/nhatthm/hidden",
					HomeURL:       "https://github.com/nhatthm/hidden",
					DirectoryURL:  "https://github.com/nhatthm/hidden/tree/master{/dir}",
					FileURL:       "https://github.com/nhatthm/hidden/blob/master{/dir}/{file}#L{line}",
				}},
			},
			{
				Name:           "Vanity Renderder",
				Path:           "vanityrender",
				RepositoryURL:  "https://github.com/nhatthm/govanityrender",
				RepositoryName: "github.com/nhatthm/govanityrender",
				LatestVersion:  "v0.1.0",
				Modules: []site.Module{{
					Path:          "vanityrender",
					ImportPrefix:  "vanityrender",
					VCS:           "git",
					RepositoryURL: "https://github.com/nhatthm/govanityrender",
					HomeURL:       "https://github.com/nhatthm/govanityrender",
					DirectoryURL:  "https://github.com/nhatthm/govanityrender/tree/master{/dir}",
					FileURL:       "https://github.com/nhatthm/govanityrender/blob/master{/dir}/{file}#L{line}",
				}},
			},
			{
				Name:           "Testcontainers Registry",
				Path:           "testcontainers-registry",
				RepositoryURL:  "https://github.com/nhatthm/testcontainers-go-registry",
				RepositoryName: "github.com/nhatthm/testcontainers-go-registry",
				LatestVersion:  "v0.6.0",
				Modules: []site.Module{{
					Path:          "testcontainers-registry",
					ImportPrefix:  "testcontainers-registry",
					VCS:           "git",
					RepositoryURL: "https://github.com/nhatthm/testcontainers-go-registry",
					HomeURL:       "https://github.com/nhatthm/testcontainers-go-registry",
					DirectoryURL:  "https://github.com/nhatthm/testcontainers-go-registry/tree/master{/dir}",
					FileURL:       "https://github.com/nhatthm/testcontainers-go-registry/blob/master{/dir}/{file}#L{line}",
				}, {
					Path:          "testcontainers-registry/elasticsearch",
					ImportPrefix:  "testcontainers-registry",
					VCS:           "git",
					RepositoryURL: "https://github.com/nhatthm/testcontainers-go-registry",
					HomeURL:       "https://github.com/nhatthm/testcontainers-go-registry",
					DirectoryURL:  "https://github.com/nhatthm/testcontainers-go-registry/tree/master{/dir}",
					FileURL:       "https://github.com/nhatthm/testcontainers-go-registry/blob/master{/dir}/{file}#L{line}",
				}, {
					Path:          "testcontainers-registry/mongo",
					ImportPrefix:  "testcontainers-registry",
					VCS:           "git",
					RepositoryURL: "https://github.com/nhatthm/testcontainers-go-registry",
					HomeURL:       "https://github.com/nhatthm/testcontainers-go-registry",
					DirectoryURL:  "https://github.com/nhatthm/testcontainers-go-registry/tree/master{/dir}",
					FileURL:       "https://github.com/nhatthm/testcontainers-go-registry/blob/master{/dir}/{file}#L{line}",
				}, {
					Path:          "testcontainers-registry/mssql",
					ImportPrefix:  "testcontainers-registry",
					VCS:           "git",
					RepositoryURL: "https://github.com/nhatthm/testcontainers-go-registry",
					HomeURL:       "https://github.com/nhatthm/testcontainers-go-registry",
					DirectoryURL:  "https://github.com/nhatthm/testcontainers-go-registry/tree/master{/dir}",
					FileURL:       "https://github.com/nhatthm/testcontainers-go-registry/blob/master{/dir}/{file}#L{line}",
				}, {
					Path:          "testcontainers-registry/mysql",
					ImportPrefix:  "testcontainers-registry",
					VCS:           "git",
					RepositoryURL: "https://github.com/nhatthm/testcontainers-go-registry",
					HomeURL:       "https://github.com/nhatthm/testcontainers-go-registry",
					DirectoryURL:  "https://github.com/nhatthm/testcontainers-go-registry/tree/master{/dir}",
					FileURL:       "https://github.com/nhatthm/testcontainers-go-registry/blob/master{/dir}/{file}#L{line}",
				}, {
					Path:          "testcontainers-registry/postgres",
					ImportPrefix:  "testcontainers-registry",
					VCS:           "git",
					RepositoryURL: "https://github.com/nhatthm/testcontainers-go-registry",
					HomeURL:       "https://github.com/nhatthm/testcontainers-go-registry",
					DirectoryURL:  "https://github.com/nhatthm/testcontainers-go-registry/tree/master{/dir}",
					FileURL:       "https://github.com/nhatthm/testcontainers-go-registry/blob/master{/dir}/{file}#L{line}",
				}},
			},
			{
				Name:           "Testcontainers Registry",
				Path:           "testcontainers-go-registry",
				RepositoryURL:  "https://github.com/nhatthm/testcontainers-go-registry",
				RepositoryName: "github.com/nhatthm/testcontainers-go-registry",
				Deprecated:     "Use go.nhat.io/testcontainers-registry instead",
				LatestVersion:  "v0.4.0",
				Modules: []site.Module{{
					Path:          "testcontainers-go-registry",
					ImportPrefix:  "testcontainers-go-registry",
					VCS:           "git",
					RepositoryURL: "https://github.com/nhatthm/testcontainers-go-registry",
					HomeURL:       "https://github.com/nhatthm/testcontainers-go-registry",
					DirectoryURL:  "https://github.com/nhatthm/testcontainers-go-registry/tree/master{/dir}",
					FileURL:       "https://github.com/nhatthm/testcontainers-go-registry/blob/master{/dir}/{file}#L{line}",
				}},
			},
		},
	}

	r, err := site.NewHandlebarsRenderder(templates.EmbeddedHomepage(), templates.EmbeddedNotFound(), templates.EmbeddedRepository(), outputDir)
	require.NoError(t, err)

	err = r.Render(s)
	require.NoError(t, err)

	assertOutput(t, "../../resources/fixtures/render_success", outputDir)
}

func assertOutput(t *testing.T, expectedDir, actualDir string) {
	t.Helper()

	expectedFiles := make(map[string]struct{})
	actualFiles := make(map[string]struct{})

	err := filepath.Walk(expectedDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relativePath, err := filepath.Rel(expectedDir, path)
		if err != nil {
			return fmt.Errorf("could not get relative path %q: %w", path, err)
		}

		expectedFiles[relativePath] = struct{}{}

		actualFile := filepath.Join(actualDir, relativePath)

		if info.IsDir() {
			assert.DirExists(t, actualFile)

			return nil
		}

		if !assert.FileExists(t, actualFile) {
			return nil
		}

		assert.Equal(t, fileContent(t, path), fileContent(t, actualFile))

		return nil
	})
	require.NoError(t, err)

	err = filepath.Walk(actualDir, func(path string, _ os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relativePath, err := filepath.Rel(actualDir, path)
		if err != nil {
			return fmt.Errorf("could not get relative path %q: %w", path, err)
		}

		actualFiles[relativePath] = struct{}{}

		return nil
	})
	require.NoError(t, err)

	assert.Equal(t, expectedFiles, actualFiles)
}

func fileContent(t *testing.T, path string) string {
	t.Helper()

	data, err := os.ReadFile(filepath.Clean(path))
	require.NoErrorf(t, err, "could not read file: %s", path, err)

	return string(bytes.TrimRight(data, "\n"))
}
