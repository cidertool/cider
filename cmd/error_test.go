package cmd

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrors(t *testing.T) {
	err := wrapError(fmt.Errorf("TEST"), "TEST")
	assert.Error(t, err)
	assert.Equal(t, "TEST", err.Error())
	assert.Equal(t, 1, err.code)
	assert.Equal(t, "TEST", err.details)
}
