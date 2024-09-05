package sitecache_test

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.nhat.io/vanityrender/internal/service/sitecache"
	"go.nhat.io/vanityrender/internal/site"
)

func TestRenderder_Render_UpstreamError(t *testing.T) {
	t.Parallel()

	upstream := mockRenderError(errors.New("upstream error"))
	h := sitecache.NewRenderder(upstream, "", "")

	err := h.Render(site.Site{})

	expected := `upstream error`

	assert.EqualError(t, err, expected)
}

func TestRenderder_Render_CouldNotWriteMetadata(t *testing.T) {
	t.Parallel()

	h := sitecache.NewRenderder(mockRender(), "unknown", "123")

	err := h.Render(site.Site{})

	expected := `could not render metadata: open unknown/metadata.v1.json: no such file or directory`

	assert.EqualError(t, err, expected)
}

func TestRenderder_Render_Success(t *testing.T) {
	t.Parallel()

	outputDir := t.TempDir()
	h := sitecache.NewRenderder(mockRender(), outputDir, "123")

	err := h.Render(site.Site{
		PageTitle: "test",
	})
	require.NoError(t, err)

	metadataFile := filepath.Join(outputDir, "metadata.v1.json")

	actual := fileContent(t, metadataFile)
	expected := `{
  "checksum": "123",
  "page_title": "test",
  "page_description": "",
  "hostname": "",
  "source_url": "",
  "repositories": null
}`

	assert.Equal(t, expected, actual)
}

func fileContent(t *testing.T, path string) string {
	t.Helper()

	data, err := os.ReadFile(filepath.Clean(path))
	require.NoErrorf(t, err, "could not read file: %s", path, err)

	return string(bytes.TrimRight(data, "\n"))
}

type renderFunc func(s site.Site) error

func (r renderFunc) Render(s site.Site) error {
	return r(s)
}

func mockRenderError(err error) renderFunc {
	return func(site.Site) error {
		return err
	}
}

func mockRender(handlers ...func(s site.Site)) renderFunc {
	return func(s site.Site) error {
		for _, h := range handlers {
			h(s)
		}

		return nil
	}
}
