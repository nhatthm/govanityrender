package git

import (
	"fmt"
	"io"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func clone(url string, ref string) (string, error) {
	dir, err := os.MkdirTemp("", "")
	if err != nil {
		return "", fmt.Errorf("could not create temporary directory: %w", err)
	}

	g, err := git.PlainClone(dir, false, &git.CloneOptions{
		URL:      url,
		Progress: io.Discard,
	})
	if err != nil {
		return "", fmt.Errorf("could not clone repository: %w", err)
	}

	commit, err := g.ResolveRevision(plumbing.Revision(ref))
	if err != nil {
		return "", fmt.Errorf("could not resolve revision %q: %w", ref, err)
	}

	w, err := g.Worktree()
	if err != nil {
		return "", fmt.Errorf("could not get worktree: %w", err)
	}

	err = w.Checkout(&git.CheckoutOptions{
		Hash: *commit,
	})
	if err != nil {
		return "", fmt.Errorf("could not checkout revision %s: %w", commit.String(), err)
	}

	return dir, nil
}
