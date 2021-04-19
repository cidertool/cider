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

package defaults

import (
	"errors"
	"testing"

	"github.com/cidertool/cider/internal/defaults"
	"github.com/cidertool/cider/pkg/config"
	"github.com/cidertool/cider/pkg/context"
	"github.com/stretchr/testify/assert"
)

var errTestError = errors.New("TEST")

func TestDefaults(t *testing.T) {
	t.Parallel()

	ctx := context.New(config.Project{})
	pipe := Pipe{}

	var err error

	assert.Equal(t, "setting defaults", pipe.String())

	err = pipe.Run(ctx)
	assert.NoError(t, err)

	pipe.defaulters = []defaults.Defaulter{
		mockDefaulter{},
	}
	err = pipe.Run(ctx)
	assert.Error(t, err)
}

type mockDefaulter struct{}

func (d mockDefaulter) String() string {
	return ""
}

func (d mockDefaulter) Default(ctx *context.Context) error {
	return errTestError
}
