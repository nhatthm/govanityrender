package git

import (
	"fmt"
	"io"
	"os"
	"sort"
	"sync"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"

	"go.nhat.io/vanityrender/internal/module"
	"go.nhat.io/vanityrender/internal/must"
)

var (
	cloneOnce = make(map[cloneRequest]*sync.Once)
	cloned    = make(map[cloneRequest]func() (string, *git.Repository, error))
	cloneMu   = sync.Mutex{}
)

type cloneRequest struct {
	Repository string
	Ref        string
}

// Clone clones a repository.
func Clone(url string, ref string) (string, *git.Repository, error) {
	cloneMu.Lock()
	defer cloneMu.Unlock()

	req := cloneRequest{
		Repository: url,
		Ref:        ref,
	}

	if _, ok := cloneOnce[req]; !ok {
		cloneOnce[req] = new(sync.Once)
	}

	cloneOnce[req].Do(func() {
		dir, r, err := clone(url, ref)

		cloned[req] = func() (string, *git.Repository, error) {
			return dir, r, err
		}
	})

	return cloned[req]()
}

func clone(url string, ref string) (string, *git.Repository, error) {
	dir, err := os.MkdirTemp("", "")
	must.NoError(err)

	r, err := git.PlainClone(dir, false, &git.CloneOptions{
		URL:      url,
		Progress: io.Discard,
	})
	if err != nil {
		return "", nil, fmt.Errorf("could not clone repository: %w", err)
	}

	if len(ref) == 0 {
		return dir, r, nil
	}

	commit, err := r.ResolveRevision(plumbing.Revision(ref))
	if err != nil {
		return "", nil, fmt.Errorf("could not resolve ref %q: %w", ref, err)
	}

	w, err := r.Worktree()
	must.NoError(err)

	err = w.Checkout(&git.CheckoutOptions{
		Hash: *commit,
	})
	if err != nil {
		return "", nil, fmt.Errorf("could not checkout revision %s: %w", commit.String(), err)
	}

	return dir, r, nil
}

// Versions returns all the versions up to HEAD in the repository.
func Versions(r *git.Repository) ([]string, error) {
	h, err := r.Head()
	if err != nil {
		return nil, fmt.Errorf("could not get head: %w", err)
	}

	headC, err := r.CommitObject(h.Hash())
	if err != nil {
		return nil, fmt.Errorf("could not get head commit: %w", err)
	}

	tags, err := r.Tags()
	if err != nil {
		return nil, fmt.Errorf("could not list tags: %w", err)
	}

	var tagNames []string

	err = tags.ForEach(func(t *plumbing.Reference) error {
		hash, err := r.ResolveRevision(plumbing.Revision(t.Name()))
		if err != nil {
			return fmt.Errorf("could not resolve tag %q: %w", t.Name(), err)
		}

		tagC, err := r.CommitObject(*hash)
		if err != nil {
			return fmt.Errorf("could not get tag commit %q: %w", t.Name(), err)
		}

		// Ignore new tags.
		if tagC.Committer.When.After(headC.Committer.When) {
			return nil
		}

		if version := t.Name().Short(); module.PathVersionRegExp.MatchString(version) {
			tagNames = append(tagNames, version)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	sort.Strings(tagNames)

	return tagNames, nil
}
