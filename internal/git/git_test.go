package git_test

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.nhat.io/vanityrender/internal/git"
)

var (
	currentTS = time.Now()
	nowMu     = sync.Mutex{}
)

func TestClone_Error_CouldNotClone(t *testing.T) {
	t.Parallel()

	_, _, err := git.Clone("not-found", "")

	expected := `could not clone repository: repository not found`

	assert.EqualError(t, err, expected)
}

func TestClone_Error_CouldNotResolveRef(t *testing.T) {
	t.Parallel()

	repo := mockRepository()(t)

	_, _, err := git.Clone(repo, "unknown")

	expected := `could not resolve ref "unknown": reference not found`

	assert.EqualError(t, err, expected)
}

func TestClone_Success(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario       string
		mockRepository func(t *testing.T) string
		ref            string
	}{
		{
			scenario:       "empty ref",
			mockRepository: mockRepository(),
		},
		{
			scenario:       "branch ref",
			mockRepository: mockRepository(),
			ref:            "master",
		},
		{
			scenario:       "tag ref",
			mockRepository: mockRepository(tagRepositoryHead("v0.5.0")),
			ref:            "v0.5.0",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			repo := tc.mockRepository(t)

			dir, r, err := git.Clone(repo, tc.ref)

			assert.NotEmpty(t, dir)
			assert.NotNil(t, r)
			assert.NoError(t, err)
		})
	}
}

func TestVersions_Success(t *testing.T) {
	t.Parallel()

	dir := mockRepository(initExampleModule())(t)

	testCases := []struct {
		scenario string
		ref      string
		expected []string
	}{
		{
			scenario: "without ref",
			expected: []string{
				"contrib/v0.1.0", "contrib/v0.2.0",
				"test/v0.1.0", "test/v0.2.0",
				"v0.1.0", "v0.1.1", "v0.2.0", "v0.3.0", "v0.4.0", "v0.5.0",
			},
		},
		{
			scenario: "with branch",
			ref:      "master",
			expected: []string{
				"contrib/v0.1.0", "contrib/v0.2.0",
				"test/v0.1.0", "test/v0.2.0",
				"v0.1.0", "v0.1.1", "v0.2.0", "v0.3.0", "v0.4.0", "v0.5.0",
			},
		},
		{
			scenario: "with newest tag",
			ref:      "v0.5.0",
			expected: []string{
				"contrib/v0.1.0", "contrib/v0.2.0",
				"test/v0.1.0", "test/v0.2.0",
				"v0.1.0", "v0.1.1", "v0.2.0", "v0.3.0", "v0.4.0", "v0.5.0",
			},
		},
		{
			scenario: "v0.4.0",
			ref:      "v0.4.0",
			expected: []string{
				"contrib/v0.1.0", "contrib/v0.2.0",
				"test/v0.1.0", "test/v0.2.0",
				"v0.1.0", "v0.1.1", "v0.2.0", "v0.3.0", "v0.4.0",
			},
		},
		{
			scenario: "v0.3.0",
			ref:      "v0.3.0",
			expected: []string{
				"contrib/v0.1.0",
				"test/v0.1.0",
				"v0.1.0", "v0.1.1", "v0.2.0", "v0.3.0",
			},
		},
		{
			scenario: "v0.2.0",
			ref:      "v0.2.0",
			expected: []string{
				"v0.1.0", "v0.1.1", "v0.2.0",
			},
		},
		{
			scenario: "v0.1.1",
			ref:      "v0.1.1",
			expected: []string{
				"v0.1.0", "v0.1.1",
			},
		},
		{
			scenario: "v0.1.0",
			ref:      "v0.1.0",
			expected: []string{
				"v0.1.0",
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			_, r, err := git.Clone(dir, tc.ref)
			require.NoError(t, err, "could not clone")

			actual, err := git.Versions(r)
			require.NoError(t, err, "could not get versions")

			assert.Equal(t, tc.expected, actual)
		})
	}
}

func mockRepository(mockers ...func(t *testing.T, r *gogit.Repository, dir string)) func(t *testing.T) string {
	return func(t *testing.T) string {
		t.Helper()

		// Create bare repo.
		barePath := t.TempDir()

		_, err := gogit.PlainInit(barePath, true)
		require.NoError(t, err, "could not init bare repository")

		// Prepare the repo.
		repoDir := t.TempDir()

		r, err := gogit.PlainInit(repoDir, false)
		require.NoError(t, err, "could not init repository")

		_, err = r.CreateRemote(&config.RemoteConfig{
			Name: "origin",
			URLs: []string{"file://" + barePath},
		})
		require.NoError(t, err, "could not create remote")

		writeFile(t, filepath.Join(repoDir, ".gitignore"), "")
		commitAndPush(t, r, "Initial commit")

		for _, mocks := range mockers {
			mocks(t, r, repoDir)
		}

		return barePath
	}
}

func tagRepositoryHead(tag string) func(t *testing.T, r *gogit.Repository, dir string) {
	return func(t *testing.T, r *gogit.Repository, dir string) {
		t.Helper()

		h, err := r.Head()
		require.NoError(t, err, "could not get head")

		_, err = r.CreateTag(tag, h.Hash(), &gogit.CreateTagOptions{
			Tagger:  &object.Signature{Name: "test", When: time.Now()},
			Message: fmt.Sprintf("tag %s", tag),
		})
		require.NoError(t, err, "could not tag")

		err = r.Push(&gogit.PushOptions{
			RefSpecs: []config.RefSpec{"refs/tags/*:refs/tags/*"},
		})
		require.NoError(t, err, "could not push")
	}
}

func initExampleModule() func(t *testing.T, r *gogit.Repository, dir string) {
	return func(t *testing.T, r *gogit.Repository, dir string) {
		t.Helper()

		// Init go module.
		writeGoMod(t, dir, "github.com/org/go-repository")
		writeFile(t, filepath.Join(dir, "example.go"), "package repository\n")

		commit(t, r, "Initial module")

		// Update code base.
		writeFile(t, filepath.Join(dir, "example.go"), `package repository

func Name() string {
	return "repository"
}
`)
		commitAndPush(t, r, "Add Name() function")

		// v0.1.0: Initial Release.
		tagHead(t, r, "v0.1.0")

		// Update code base.
		writeFile(t, filepath.Join(dir, "example.go"), `package repository

func Name() string {
	return "repository"
}

func ShortName() string {
	return "repo"
}
`)
		commitAndPush(t, r, "Add ShortName() function")

		// v0.1.0: Add ShortName() function.
		tagHead(t, r, "v0.1.1")

		// Update code base.
		writeFile(t, filepath.Join(dir, "example.go"), `package repository

func Name() string {
	return "github.com/org/repository"
}

func ShortName() string {
	return "repository"
}
`)
		commitAndPush(t, r, "Correct the names")

		// v0.2.0: Correct the names.
		tagHead(t, r, "v0.2.0")

		// Change module path.
		writeGoMod(t, dir, "host.tld/go-repository")
		writeFile(t, filepath.Join(dir, "example.go"), `package repository

func Name() string {
	return "host.tld/go-repository"
}

func ShortName() string {
	return "repository"
}
`)
		commitAndPush(t, r, "Change module path")

		// v0.3.0: Change module path.
		tagHead(t, r, "v0.3.0")

		// Add submodules.
		writeGoMod(t, filepath.Join(dir, "contrib"), "host.tld/go-repository/contrib")
		writeFile(t, filepath.Join(dir, "contrib", "contrib.go"), `package contrib

func Name() string {
	return "host.tld/go-repository/contrib"
}
`)
		writeGoMod(t, filepath.Join(dir, "test"), "host.tld/go-repository/test")
		writeFile(t, filepath.Join(dir, "test", "test.go"), `package test

func Name() string {
	return "host.tld/go-repository/test"
}
`)

		// contrib/v0.1.0: Initial Release.
		tagHead(t, r, "contrib/v0.1.0")
		// test/v0.1.0: Initial Release.
		tagHead(t, r, "test/v0.1.0")

		// Change module path.
		writeGoMod(t, dir, "host.tld/repository")
		writeFile(t, filepath.Join(dir, "example.go"), `package repository

func Name() string {
	return "host.tld/repository"
}

func ShortName() string {
	return "repository"
}
`)
		commitAndPush(t, r, "Change module path")

		// v0.4.0: Change module path.
		tagHead(t, r, "v0.4.0")

		// Change submodules path.
		writeGoMod(t, filepath.Join(dir, "contrib"), "host.tld/repository/contrib")
		writeFile(t, filepath.Join(dir, "contrib", "contrib.go"), `package contrib

func Name() string {
	return "host.tld/repository/contrib"
}
`)
		writeGoMod(t, filepath.Join(dir, "test"), "host.tld/repository/test")
		writeFile(t, filepath.Join(dir, "test", "test.go"), `package test

func Name() string {
	return "host.tld/repository/test"
}
`)
		// contrib/v0.2.0: Change module path.
		tagHead(t, r, "contrib/v0.2.0")
		// test/v0.2.0: Change module path.
		tagHead(t, r, "test/v0.2.0")

		// Update code base.
		writeFile(t, filepath.Join(dir, "example.go"), `package repository

func FQDN() string {
	return "host.tld/repository"
}

func Name() string {
	return "repository"
}

func ShortName() string {
	return "repo"
}
`)
		commitAndPush(t, r, "Add FQDN() function")

		// Add empty test file.
		writeFile(t, filepath.Join(dir, "test.go"), "package repository_test\n")
		commitAndPush(t, r, "Add empty test file")

		// Add README.md
		writeFile(t, filepath.Join(dir, "README.md"), "")
		commitAndPush(t, r, "Add README.md")

		tagHead(t, r, "v0.5.0")
	}
}

func commit(t *testing.T, r *gogit.Repository, message string) {
	t.Helper()

	w, err := r.Worktree()
	require.NoError(t, err, "could not get worktree")

	_, err = w.Add(".")
	require.NoError(t, err, "could not stage files")

	ts := now()

	_, err = w.Commit(message, &gogit.CommitOptions{
		Author:    sign(ts),
		Committer: sign(ts),
	})
	require.NoError(t, err, "could not commit")
}

func commitAndPush(t *testing.T, r *gogit.Repository, message string) {
	t.Helper()

	commit(t, r, message)

	err := r.Push(&gogit.PushOptions{})
	require.NoError(t, err, "could not push")
}

func tagHead(t *testing.T, r *gogit.Repository, tag string) {
	t.Helper()

	h, err := r.Head()
	require.NoError(t, err, "could not get head")

	_, err = r.CreateTag(tag, h.Hash(), &gogit.CreateTagOptions{
		Tagger:  sign(now()),
		Message: fmt.Sprintf("tag %s", tag),
	})
	require.NoError(t, err, "could not tag")

	err = r.Push(&gogit.PushOptions{
		RefSpecs: []config.RefSpec{"refs/tags/*:refs/tags/*"},
	})
	require.NoError(t, err, "could not push")
}

func sign(when time.Time) *object.Signature {
	return &object.Signature{
		Name:  "John Doe",
		Email: "john.doe@example.com",
		When:  when,
	}
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

func now() time.Time {
	nowMu.Lock()
	defer nowMu.Unlock()

	currentTS = currentTS.Add(time.Second)

	return currentTS
}
