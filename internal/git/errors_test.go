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

package git

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrDirtyMessage(t *testing.T) {
	t.Parallel()

	err := ErrDirty{"TEST"}
	expected := "git is currently in a dirty state, please check in your pipeline what can be changing the following files:\nTEST"
	assert.Equal(t, expected, err.Error())
}

func TestErrWrongRefMessage(t *testing.T) {
	t.Parallel()

	err := ErrWrongRef{"TEST", "TEST"}
	expected := "git tag TEST was not made against commit TEST"
	assert.Equal(t, expected, err.Error())
}

func TestErrNotRepositoryMessage(t *testing.T) {
	t.Parallel()

	err := ErrNotRepository{"TEST"}
	expected := "the directory at TEST is not a git repository"
	assert.Equal(t, expected, err.Error())
}
