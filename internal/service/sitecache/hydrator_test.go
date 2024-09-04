package sitecache_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.nhat.io/vanityrender/internal/service/sitecache"
	"go.nhat.io/vanityrender/internal/site"
)

func TestMetadataHydrator_Hydrate_Error_CouldNotSendRequest(t *testing.T) {
	t.Parallel()

	h := sitecache.NewMetadataHydrator("123")

	actual := h.Hydrate(&site.Site{Hostname: "https://localhost"})
	expected := `failed to send request: Get "http://https//localhost/metadata.v1.json":`

	assert.ErrorContains(t, actual, expected)
}

func TestMetadataHydrator_Hydrate_MissingHostname(t *testing.T) {
	t.Parallel()

	h := sitecache.NewMetadataHydrator("123")

	actual := site.Site{
		PageTitle: "test",
	}

	err := h.Hydrate(&actual)
	require.NoError(t, err)

	expected := site.Site{
		PageTitle: "test",
	}

	assert.Equal(t, expected, actual)
}

func TestMetadataHydrator_Hydrate(t *testing.T) {
	t.Parallel()

	const checksum = "123"

	testCases := []struct {
		scenario       string
		mockServer     func(t *testing.T) string
		expectedResult site.Site
		expectedError  string
	}{
		{
			scenario: "404",
			mockServer: mockMetadataServer(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			}),
			expectedResult: site.Site{},
			expectedError:  `metadata not found`,
		},
		{
			scenario: "500",
			mockServer: mockMetadataServer(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			}),
			expectedResult: site.Site{},
			expectedError:  `unexpected status code: 500 Internal Server Error`,
		},
		{
			scenario: "could not decode metadata",
			mockServer: mockMetadataServer(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`43`)) // nolint: errcheck
			}),
			expectedResult: site.Site{},
			expectedError:  `invalid metadata: json: cannot unmarshal number into Go value of type sitecache.metadata`,
		},
		{
			scenario: "checksum mismatched",
			mockServer: mockMetadataServer(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"checksum": "456", "page_title":"Test"}`)) // nolint: errcheck
			}),
			expectedResult: site.Site{},
			expectedError:  "checksum mismatched",
		},
		{
			scenario: "success",
			mockServer: mockMetadataServer(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"checksum": "123", "page_title":"Test"}`)) // nolint: errcheck
			}),
			expectedResult: site.Site{
				PageTitle: "Test",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			h := sitecache.NewMetadataHydrator(checksum)
			actual := site.Site{
				Hostname: tc.mockServer(t),
			}

			err := h.Hydrate(&actual)

			tc.expectedResult.Hostname = actual.Hostname

			assert.Equal(t, tc.expectedResult, actual)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func mockMetadataServer(h http.HandlerFunc) func(t *testing.T) string {
	return mockServer(func(srv *http.ServeMux) {
		srv.Handle("/metadata.v1.json", h)
	})
}

func mockServer(mocks ...func(srv *http.ServeMux)) func(t *testing.T) string {
	return func(t *testing.T) string {
		t.Helper()

		mux := http.NewServeMux()

		for _, mock := range mocks {
			mock(mux)
		}

		srv := httptest.NewServer(mux)

		t.Cleanup(func() {
			srv.Close()
		})

		return strings.ReplaceAll(srv.URL, "http://", "")
	}
}
