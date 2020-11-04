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

package pipe

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPipeSkip(t *testing.T) {
	skip := Skip("TEST")

	var err error = ErrSkip{reason: "TEST"}

	assert.Error(t, skip)
	assert.Error(t, err)
	assert.True(t, IsSkip(err))
	assert.Equal(t, err, skip)
	assert.EqualError(t, err, skip.Error())
}

func TestErrMissingApp(t *testing.T) {
	err := ErrMissingApp{Name: "TEST"}
	assert.Error(t, err)
	assert.EqualError(t, err, "no app defined in configuration matching the name TEST")
}
