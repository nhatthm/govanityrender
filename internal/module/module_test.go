package module_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go.nhat.io/vanityrender/internal/module"
)

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
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			actualPath, actualVersion := module.PathVersion(tc.value)

			assert.Equal(t, tc.expectedPath, actualPath)
			assert.Equal(t, tc.expectedVersion, actualVersion)
		})
	}
}
