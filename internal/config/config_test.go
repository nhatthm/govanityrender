package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.nhat.io/vanityrender/internal/config"
)

func TestFromFile(t *testing.T) {
	t.Parallel()

	const (
		payloadBroken      = `{`
		payloadMissingHost = `{}`
		payloadOK          = `{
    "page_title": "",
    "host": "go.nhat.io",
    "repositories": [
        {
            "name": "Vanity Renderder",
            "path": "vanityrender",
            "repository": "https://github.com/nhatthm/govanityrender"
        }
    ]
}`
	)

	testCases := []struct {
		scenario             string
		file                 string
		expectedResult       config.Config
		expectedError        error
		expectedErrorMessage string
	}{
		{
			scenario:             "file not found",
			file:                 "not-found",
			expectedError:        config.ErrCouldNotReadConfigFile,
			expectedErrorMessage: "could not read config file: open not-found: no such file or directory",
		},
		{
			scenario:             "payload is broken",
			file:                 testFile(t, "broken.json", payloadBroken),
			expectedError:        config.ErrInvalidConfig,
			expectedErrorMessage: "invalid config: unexpected EOF",
		},
		{
			scenario:             "missing host",
			file:                 testFile(t, "missing_host.json", payloadMissingHost),
			expectedError:        config.ErrMissingHost,
			expectedErrorMessage: "missing host",
		},
		{
			scenario: "success",
			file:     testFile(t, "success.json", payloadOK),
			expectedResult: config.Config{
				PageTitle: "go.nhat.io",
				Host:      "go.nhat.io",
				Repositories: []config.Repository{
					{
						Name:       "Vanity Renderder",
						Path:       "vanityrender",
						Repository: "https://github.com/nhatthm/govanityrender",
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			actual, err := config.FromFile(tc.file)

			assert.Equal(t, tc.expectedResult, actual)

			if tc.expectedError == nil {
				assert.NoError(t, err)
			} else {
				assert.ErrorIs(t, err, tc.expectedError)
				assert.EqualError(t, err, tc.expectedErrorMessage)
			}
		})
	}
}

func testFile(t *testing.T, name, content string) string {
	t.Helper()

	file := filepath.Join(t.TempDir(), name)

	err := os.WriteFile(file, []byte(content), 0o644) // nolint: gosec
	require.NoError(t, err)

	return file
}
