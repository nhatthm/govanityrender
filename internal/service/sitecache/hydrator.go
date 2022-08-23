package sitecache

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.nhat.io/vanityrender/internal/must"
	"go.nhat.io/vanityrender/internal/site"
)

const metadataFile = `metadata.v1.json`

var _ site.Hydrator = (*Hydrator)(nil)

type metadata struct {
	Checksum string `json:"checksum"`
	site.Site
}

// Hydrator hydrates configuration using the metadata file.
type Hydrator struct {
	client *http.Client

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

	// Silently ignore all the errors.
	if resp.StatusCode != http.StatusOK {
		return nil, nil // nolint: nilnil // Ignore the error.
	}

	var m metadata

	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&m); err != nil {
		return nil, nil // nolint: nilnil,nilerr // Ignore the error.
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

	if m == nil || h.checksum != m.Checksum {
		return nil
	}

	*s = m.Site

	return nil
}

// NewMetadataHydrator hydrates configuration using the metadata file.
func NewMetadataHydrator(checksum string) *Hydrator {
	return &Hydrator{
		client:   &http.Client{},
		checksum: checksum,
		timeout:  10 * time.Second,
	}
}
