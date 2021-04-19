/**
Copyright (C) 2020 Aaron Sky.

This file is part of Cider, a tool for automating submission
of apps to Apple's App Stores.

Cider is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

Cider is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with Cider.  If not, see <http://www.gnu.org/licenses/>.
*/

package config

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var errTestError = errors.New("test error")

func TestValidConfiguration(t *testing.T) {
	t.Parallel()

	f, err := Load("testdata/valid.yml")
	assert.NoError(t, err)
	assert.Len(t, f, 1)
}

func TestMissingConfiguration(t *testing.T) {
	t.Parallel()

	_, err := Load("testdata/doesnotexist.yml")
	assert.Error(t, err)
}

func TestInvalidConfiguration(t *testing.T) {
	t.Parallel()

	_, err := Load("testdata/invalid.yml")
	assert.Error(t, err)
}

func TestMarshalledIsValidConfiguration(t *testing.T) {
	t.Parallel()

	f, err := Load("testdata/valid.yml")
	assert.NoError(t, err)
	str, err := f.String()
	assert.NoError(t, err)
	f2, err := LoadReader(strings.NewReader(str))
	assert.NoError(t, err)
	assert.Equal(t, f, f2)
}

func TestBrokenFile(t *testing.T) {
	t.Parallel()

	_, err := LoadReader(errReader(0))
	assert.Error(t, err)
}

type errReader int

func (errReader) Read(p []byte) (int, error) {
	return 0, errTestError
}

func TestCopy(t *testing.T) {
	t.Parallel()

	p := Project{
		"App1": {},
		"App2": {},
		"App3": {},
	}
	pPrime, err := p.Copy()
	assert.NoError(t, err)
	assert.Equal(t, p, pPrime)
	assert.NotSame(t, p, pPrime)
}

func TestCopy_Err(t *testing.T) {
	t.Parallel()

	p := Project{
		"App1": {},
		"App2": {},
		"App3": {},
	}
	pPrime, err := p.Copy()
	assert.NoError(t, err)
	assert.Equal(t, p, pPrime)
	assert.NotSame(t, p, pPrime)
}

func TestAppsMatching(t *testing.T) {
	t.Parallel()

	p := Project{
		"App1": {},
		"App2": {},
		"App3": {},
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
