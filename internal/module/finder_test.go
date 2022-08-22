package module_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.nhat.io/vanityrender/internal/module"
)

func TestFindVersions_Success(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario   string
		mockModule func(t *testing.T) string
		expected   []string
	}{
		{
			scenario:   "only v0 - no submodules",
			mockModule: mockModuleV0,
			expected:   []string{"v0.0.0"},
		},
		{
			scenario:   "only v2 - no submodules",
			mockModule: mockModuleV2,
			expected:   []string{"v2.0.0"},
		},
		{
			scenario:   "v0 with submodules",
			mockModule: mockModuleV0WithSubmodules,
			expected:   []string{"contrib/v0.0.0", "test/v3.0.0", "v0.0.0"},
		},
		{
			scenario:   "v2 with submodules",
			mockModule: mockModuleV2WithSubmodules,
			expected:   []string{"contrib/v0.0.0", "test/v3.0.0", "v2.0.0"},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			actual, err := module.FindVersions(tc.mockModule(t))
			require.NoError(t, err)

			assert.Equal(t, tc.expected, actual)
		})
	}
}

func mockModuleV0(t *testing.T) string {
	t.Helper()

	dir := t.TempDir()

	writeGoMod(t, dir, "example.com/module")

	return dir
}

func mockModuleV0WithSubmodules(t *testing.T) string {
	t.Helper()

	dir := t.TempDir()

	writeGoMod(t, dir, "example.com/module")
	writeGoMod(t, filepath.Join(dir, "contrib"), "example.com/module/contrib")
	writeGoMod(t, filepath.Join(dir, "test"), "example.com/module/test/v3")

	return dir
}

func mockModuleV2(t *testing.T) string {
	t.Helper()

	dir := t.TempDir()

	writeGoMod(t, dir, "example.com/module/v2")

	return dir
}

func mockModuleV2WithSubmodules(t *testing.T) string {
	t.Helper()

	dir := t.TempDir()

	writeGoMod(t, dir, "example.com/module/v2")
	writeGoMod(t, filepath.Join(dir, "contrib"), "example.com/module/contrib")
	writeGoMod(t, filepath.Join(dir, "test"), "example.com/module/test/v3")

	return dir
}

func writeFile(t *testing.T, path, data string) {
	t.Helper()

	dir := filepath.Dir(path)

	err := os.MkdirAll(dir, 0o755) // nolint: gosec
	require.NoError(t, err)

	err = os.WriteFile(filepath.Clean(path), []byte(data), 0o644) // nolint: gosec
	require.NoError(t, err)
}

func writeGoMod(t *testing.T, dir, module string) {
	t.Helper()

	data := `module %s

go 1.18
`

	data = fmt.Sprintf(data, module)

	writeFile(t, filepath.Join(dir, "go.mod"), data)
}
