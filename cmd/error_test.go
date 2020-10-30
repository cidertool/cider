package cmd

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

var errTestError = errors.New("TEST")

func TestErrors(t *testing.T) {
	err := wrapError(errTestError, "TEST")
	assert.Error(t, err)
	assert.Equal(t, "TEST", err.Error())
	assert.Equal(t, 1, err.code)
	assert.Equal(t, "TEST", err.details)
}
