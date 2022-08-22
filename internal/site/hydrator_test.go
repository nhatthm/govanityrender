package site_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.nhat.io/vanityrender/internal/site"
)

func TestHydrate(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario string
		hydrator site.Hydrator
		expected error
	}{
		{
			scenario: "has error",
			hydrator: hydrateFunc(func(s *site.Site) error {
				return errors.New("error")
			}),
			expected: errors.New("error"),
		},
		{
			scenario: "no error",
			hydrator: hydrateFunc(func(s *site.Site) error {
				return nil
			}),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			actual := site.Hydrate(&site.Site{}, tc.hydrator)

			assert.Equal(t, tc.expected, actual)
		})
	}
}

type hydrateFunc func(s *site.Site) error

func (f hydrateFunc) Hydrate(s *site.Site) error {
	return f(s)
}
