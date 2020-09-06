package git

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrDirtyMessage(t *testing.T) {
	err := ErrDirty{"TEST"}
	expected := "git is currently in a dirty state, please check in your pipeline what can be changing the following files:\nTEST"
	assert.Equal(t, expected, err.Error())
}

func TestErrWrongRefMessage(t *testing.T) {
	err := ErrWrongRef{"TEST", "TEST"}
	expected := "git tag TEST was not made against commit TEST"
	assert.Equal(t, expected, err.Error())
}

func TestErrNotRepositoryMessage(t *testing.T) {
	err := ErrNotRepository{"TEST"}
	expected := "the directory at TEST is not a git repository"
	assert.Equal(t, expected, err.Error())
}
