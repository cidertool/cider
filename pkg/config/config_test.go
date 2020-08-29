package config

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidConfiguration(t *testing.T) {
	f, err := Load("testdata/valid.yml")
	assert.NoError(t, err)
	assert.Equal(t, "My Project", f.Name)
}

func TestMissingConfiguration(t *testing.T) {
	_, err := Load("testdata/doesnotexist.yml")
	assert.Error(t, err)
}

func TestInvalidConfiguration(t *testing.T) {
	_, err := Load("testdata/invalid.yml")
	assert.Error(t, err)
}

func TestBrokenFile(t *testing.T) {
	_, err := LoadReader(errReader(0))
	assert.Error(t, err)
}

type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("test error")
}
