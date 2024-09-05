package module_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go.nhat.io/vanityrender/internal/module"
)

func TestPath_IsRoot(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario string
		path     module.Path
		expected bool
	}{
		{
			scenario: "empty string",
			path:     "",
			expected: false,
		},
		{
			scenario: "root",
			path:     ".",
			expected: true,
		},
		{
			scenario: "version",
			path:     "v2",
			expected: true,
		},
		{
			scenario: "submodule",
			path:     "contrib",
			expected: false,
		},
		{
			scenario: "submodule/version",
			path:     "contrib/v2",
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tc.expected, tc.path.IsRoot())
		})
	}
}

func TestVersion_LessThan(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario string
		left     module.Version
		right    module.Version
		expected bool
	}{
		{
			scenario: "major version less than",
			left:     module.NewVersion(1, 0, 0),
			right:    module.NewVersion(2, 0, 0),
			expected: true,
		},
		{
			scenario: "major version greater than",
			left:     module.NewVersion(2, 0, 0),
			right:    module.NewVersion(1, 0, 0),
			expected: false,
		},
		{
			scenario: "minor version less than",
			left:     module.NewVersion(1, 1, 0),
			right:    module.NewVersion(1, 2, 0),
			expected: true,
		},
		{
			scenario: "minor version greater than",
			left:     module.NewVersion(1, 2, 0),
			right:    module.NewVersion(1, 1, 0),
			expected: false,
		},
		{
			scenario: "patch version less than",
			left:     module.NewVersion(1, 2, 1),
			right:    module.NewVersion(1, 2, 2),
			expected: true,
		},
		{
			scenario: "patch version greater than",
			left:     module.NewVersion(1, 2, 2),
			right:    module.NewVersion(1, 2, 1),
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tc.expected, tc.left.LessThan(tc.right))
		})
	}
}

func TestVersion_String(t *testing.T) {
	t.Parallel()

	actual := module.NewVersion(1, 2, 3).String()
	expected := "v1.2.3"

	assert.Equal(t, expected, actual)
}

func TestPathVersion(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario        string
		value           string
		expectedPath    module.Path
		expectedVersion module.Version
	}{
		{
			scenario: "empty string",
			value:    "",
		},
		{
			scenario: "only path",
			value:    "contrib/module",
		},
		{
			scenario:        "only version - v0",
			value:           "v0.3.0",
			expectedPath:    ".",
			expectedVersion: module.NewVersion(0, 3, 0),
		},
		{
			scenario:        "only version - v1",
			value:           "v1.2.0",
			expectedPath:    ".",
			expectedVersion: module.NewVersion(1, 2, 0),
		},
		{
			scenario:        "only version - v2",
			value:           "v2.3.0",
			expectedPath:    "v2",
			expectedVersion: module.NewVersion(2, 3, 0),
		},
		{
			scenario:        "full path and version - v0",
			value:           "contrib/v0.3.0",
			expectedPath:    "contrib",
			expectedVersion: module.NewVersion(0, 3, 0),
		},
		{
			scenario:        "full path and version - v1",
			value:           "contrib/v1.2.0",
			expectedPath:    "contrib",
			expectedVersion: module.NewVersion(1, 2, 0),
		},
		{
			scenario:        "full path and version - v2",
			value:           "contrib/v2.3.0",
			expectedPath:    "contrib/v2",
			expectedVersion: module.NewVersion(2, 3, 0),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			actualPath, actualVersion := module.PathVersion(tc.value)

			assert.Equal(t, tc.expectedPath, actualPath)
			assert.Equal(t, tc.expectedVersion, actualVersion)
		})
	}
}

func TestPathWithVersion(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario string
		path     string
		version  module.Version
		expected string
	}{
		{
			scenario: "empty path",
			path:     "",
			version:  module.NewVersion(1, 2, 3),
			expected: "v1",
		},
		{
			scenario: "path is dot",
			path:     ".",
			version:  module.NewVersion(1, 2, 3),
			expected: "v1",
		},
		{
			scenario: "simple path",
			path:     "contrib",
			version:  module.NewVersion(1, 2, 3),
			expected: "contrib/v1",
		},
		{
			scenario: "deep path",
			path:     "contrib/test",
			version:  module.NewVersion(1, 2, 3),
			expected: "contrib/test/v1",
		},
		{
			scenario: "leading slash",
			path:     "/contrib/test",
			version:  module.NewVersion(1, 2, 3),
			expected: "contrib/test/v1",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			actual := module.PathWithVersion(tc.path, tc.version)

			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestPathWithoutVersion(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario string
		path     string
		expected string
	}{
		{
			scenario: "empty path",
			path:     "",
			expected: "",
		},
		{
			scenario: "dot path",
			path:     ".",
			expected: ".",
		},
		{
			scenario: "only version",
			path:     "v2",
			expected: ".",
		},
		{
			scenario: "simple path with version",
			path:     "contrib/v2",
			expected: "contrib",
		},
		{
			scenario: "deep path with version",
			path:     "contrib/test/v2",
			expected: "contrib/test",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			actual := module.PathWithoutVersion(tc.path)

			assert.Equal(t, tc.expected, actual)
		})
	}
}
