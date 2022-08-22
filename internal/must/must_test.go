package must_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.nhat.io/vanityrender/internal/must"
)

func TestNoError_Panic(t *testing.T) {
	t.Parallel()

	assert.Panics(t, func() {
		must.NoError(errors.New("error"))
	})
}

func TestNoError_NoPanic(t *testing.T) {
	t.Parallel()

	assert.NotPanics(t, func() {
		must.NoError(nil)
	})
}
