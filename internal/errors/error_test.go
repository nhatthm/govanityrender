package errors_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go.nhat.io/vanityrender/internal/errors"
)

func TestError_Error(t *testing.T) {
	t.Parallel()

	e := errors.Error("error message")

	actual := e.Error()
	expected := "error message"

	assert.Equal(t, expected, actual)
}
