package sitecache

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	xerrors "go.nhat.io/vanityrender/internal/errors"
	"go.nhat.io/vanityrender/internal/must"
	"go.nhat.io/vanityrender/internal/site"
)

const metadataFile = `metadata.v1.json`

var (
	// ErrChecksumMismatched indicates that the checksum of the metadata file does not match the checksum of the site.
	ErrChecksumMismatched = xerrors.Error("checksum mismatched")
	// ErrMetadataNotFound indicates that the metadata file is not found.
	ErrMetadataNotFound = xerrors.Error("metadata not found")
	// ErrMetadataInvalid indicates that the metadata file is invalid.
	ErrMetadataInvalid = xerrors.Error("invalid metadata")
)

var _ site.Hydrator = (*Hydrator)(nil)

type metadata struct {
	Checksum string `json:"checksum"`
	site.Site
}

// Hydrator hydrates configuration using the metadata file.
type Hydrator struct {
	client *http.Client

	output   io.Writer
	checksum string
	timeout  time.Duration
}

func (h *Hydrator) metadata(host string) (*metadata, error) {
	ctx, cancel := context.WithTimeout(context.Background(), h.timeout)
	defer cancel()

	url := fmt.Sprintf("http://%s/%s", host, metadataFile)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	must.NoError(err)

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close() // nolint: errcheck

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrMetadataNotFound
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %s", resp.Status) // nolint: goerr113
	}

	var m metadata

	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&m); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrMetadataInvalid, err) //nolint: errorlint
	}

	return &m, nil
}

// Hydrate hydrates configuration using the metadata file.
func (h *Hydrator) Hydrate(s *site.Site) error {
	if len(s.Hostname) == 0 {
		return nil
	}

	m, err := h.metadata(s.Hostname)
	if err != nil {
		return err
	}

	if h.checksum != m.Checksum {
		return ErrChecksumMismatched
	}

	*s = m.Site

	return nil
}

// NewMetadataHydrator hydrates configuration using the metadata file.
func NewMetadataHydrator(checksum string, opts ...HydratorOption) *Hydrator {
	h := &Hydrator{
		client:   &http.Client{},
		output:   io.Discard,
		checksum: checksum,
		timeout:  10 * time.Second,
	}

	for _, o := range opts {
		o.applyHydratorOption(h)
	}

	return h
}

// HydratorOption is an option to configure Hydrator.
type HydratorOption interface {
	applyHydratorOption(r *Hydrator)
}

type hydratorOptionFunc func(r *Hydrator)

func (f hydratorOptionFunc) applyHydratorOption(r *Hydrator) {
	f(r)
}
