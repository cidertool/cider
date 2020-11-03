package cmd

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRootCmd(t *testing.T) {
	exit := func(code int) {
		assert.Equal(t, 0, code)
	}
	err := os.Setenv("CI", "1")
	assert.NoError(t, err)
	Execute("TEST", exit, []string{"help", "--debug"})
}

func TestRootCmd_Error(t *testing.T) {
	exit := func(code int) {
		assert.NotEqual(t, 0, code)
	}
	Execute("TEST", exit, []string{"check"})
}
