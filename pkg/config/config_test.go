package config

import (
	"errors"
	"strings"
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

func TestMarshalledIsValidConfiguration(t *testing.T) {
	f, err := Load("testdata/valid.yml")
	assert.NoError(t, err)
	str, err := f.String()
	assert.NoError(t, err)
	f2, err := LoadReader(strings.NewReader(str))
	assert.NoError(t, err)
	assert.Equal(t, f, f2)
}

func TestBrokenFile(t *testing.T) {
	_, err := LoadReader(errReader(0))
	assert.Error(t, err)
}

type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("test error")
}

func TestAppsMatching(t *testing.T) {
	p := Project{
		Apps: map[string]App{
			"App1": {},
			"App2": {},
			"App3": {},
		},
	}
	var matches []string
	matches = p.AppsMatching([]string{"App1", "App2", "App3"}, false)
	assert.ElementsMatch(t, matches, []string{"App1", "App2", "App3"})
	matches = p.AppsMatching([]string{"App1", "App2"}, false)
	assert.ElementsMatch(t, matches, []string{"App1", "App2"})
	matches = p.AppsMatching([]string{"App1", "App2", "App4"}, false)
	assert.ElementsMatch(t, matches, []string{"App1", "App2"})
	matches = p.AppsMatching([]string{"", ""}, false)
	assert.ElementsMatch(t, matches, []string{})
	matches = p.AppsMatching([]string{"App1", "App2", "App4"}, true)
	assert.ElementsMatch(t, matches, []string{"App1", "App2", "App3"})
	matches = p.AppsMatching([]string{}, true)
	assert.ElementsMatch(t, matches, []string{"App1", "App2", "App3"})
}
