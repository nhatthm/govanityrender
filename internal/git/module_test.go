package git_test

import (
	"path/filepath"
	"testing"

	gogit "github.com/go-git/go-git/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.nhat.io/vanityrender/internal/git"
	"go.nhat.io/vanityrender/internal/module"
)

func TestModuleFinder_Find_Error_CouldNotClone(t *testing.T) {
	t.Parallel()

	f := git.NewModuleFinder()
	_, err := f.Find("not-found", "")

	expected := `could not clone repository: repository not found`

	assert.EqualError(t, err, expected)
}

func TestModuleFinder_Find_Success(t *testing.T) {
	t.Parallel()

	dir := mockRepository(initExampleModule(), bumpExampleModule())(t)
	f := git.NewModuleFinder()

	actual, err := f.Find(dir, "")
	require.NoError(t, err, "could not find modules")

	expected := map[module.Path]module.Version{
		".":          module.NewVersionFromString("v1.0.0"),
		"v2":         module.NewVersionFromString("v2.10.0"),
		"contrib":    module.NewVersionFromString("v0.2.0"),
		"contrib/v2": module.NewVersionFromString("v2.0.0"),
		"test":       module.NewVersionFromString("v0.2.0"),
	}

	assert.Equal(t, expected, actual)
}

func bumpExampleModule() func(t *testing.T, r *gogit.Repository, dir string) {
	return func(t *testing.T, r *gogit.Repository, dir string) {
		t.Helper()

		// Bump main module version to v1.0.0.
		writeFile(t, filepath.Join(dir, "VERSION"), "v1.0.0")
		commitAndPush(t, r, "Add VERSION")

		tagHead(t, r, "v1.0.0")

		// Bump contrib module version to v2.0.0.
		writeGoMod(t, filepath.Join(dir, "contrib"), "host.tld/repository/contrib/v2")
		commitAndPush(t, r, "Bump contrib to v2.0.0")

		tagHead(t, r, "contrib/v2.0.0")

		// Bump main module version to v2.2.0.
		writeFile(t, filepath.Join(dir, "VERSION"), "v2.2.0")
		commitAndPush(t, r, "Bump VERSION")

		tagHead(t, r, "v2.2.0")

		// Bump main module version to v2.10.0.
		writeFile(t, filepath.Join(dir, "VERSION"), "v2.10.0")
		commitAndPush(t, r, "Bump VERSION")

		tagHead(t, r, "v2.10.0")
	}
}
