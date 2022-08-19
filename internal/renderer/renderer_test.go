package renderer_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.nhat.io/vanityrender/internal/config"
	"go.nhat.io/vanityrender/internal/renderer"
	"go.nhat.io/vanityrender/templates"
)

func TestNewHandlebarsRenderder(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario      string
		homepageSrc   string
		librarySrc    string
		expectedError string
	}{
		{
			scenario:      "homepage template is broken",
			homepageSrc:   `{{`,
			expectedError: "could not parse homepage template: Parse error on line 1:",
		},
		{
			scenario:      "library template is broken",
			homepageSrc:   `{{ message }}`,
			librarySrc:    `{{`,
			expectedError: "could not parse repository template: Parse error on line 1:",
		},
		{
			scenario:    "success",
			homepageSrc: `{{ message }}`,
			librarySrc:  `{{ message }}`,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			actual, err := renderer.NewHandlebarsRenderder(tc.homepageSrc, tc.librarySrc, "")

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

	cfg := config.Config{
		PageTitle: "go.nhat.io",
		Host:      "go.nhat.io",
		Repositories: []config.Repository{
			{
				Library:    "Vanity Renderder",
				Path:       "vanityrender",
				Repository: "github.com/nhatthm/govanityrender",
			},
			{
				Library:    "Testcontainers Registry",
				Path:       "testcontainers-registry",
				Repository: "github.com/nhatthm/testcontainers-go-registry",
				Modules: []config.Module{
					{Path: "elasticsearch"},
					{Path: "mongo"},
					{Path: "mssql"},
					{Path: "mysql"},
					{Path: "postgres"},
				},
			},
			{
				Library:    "Testcontainers Registry",
				Path:       "testcontainers-go-registry",
				Repository: "github.com/nhatthm/testcontainers-go-registry",
				Deprecated: "Use go.nhat.io/testcontainers-registry instead",
			},
		},
	}

	r, err := renderer.NewHandlebarsRenderder(templates.EmbeddedHomepage(), templates.EmbeddedRepository(), outputDir)
	require.NoError(t, err)

	err = r.Render(cfg)
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

		expectedContent, err := os.ReadFile(filepath.Clean(path))
		if err != nil {
			return fmt.Errorf("could not read file %q: %w", path, err)
		}

		actualContent, err := os.ReadFile(filepath.Clean(actualFile))
		if err != nil {
			return fmt.Errorf("could not read file %q: %w", actualFile, err)
		}

		assert.Equal(t, string(expectedContent), string(actualContent))

		return nil
	})
	require.NoError(t, err)

	err = filepath.Walk(actualDir, func(path string, info os.FileInfo, err error) error {
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
