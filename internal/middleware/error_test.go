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

package middleware

import (
	"errors"
	"testing"

	"github.com/cidertool/cider/internal/pipe"
	"github.com/cidertool/cider/pkg/config"
	"github.com/cidertool/cider/pkg/context"
	"github.com/stretchr/testify/assert"
)

var errTestError = errors.New("TEST")

func TestErrHandler_WrapsError(t *testing.T) {
	t.Parallel()

	ctx := context.New(config.Project{})
	wrapped := ErrHandler(func(ctx *context.Context) error {
		return errTestError
	})
	err := wrapped(ctx)
	assert.Error(t, err)
}

func TestErrHandler_IgnoresNoError(t *testing.T) {
	t.Parallel()

	ctx := context.New(config.Project{})
	wrapped := ErrHandler(func(ctx *context.Context) error {
		return nil
	})
	err := wrapped(ctx)
	assert.NoError(t, err)
}

func TestErrHandler_HandlesSkip(t *testing.T) {
	t.Parallel()

	ctx := context.New(config.Project{})
	wrapped := ErrHandler(func(ctx *context.Context) error {
		return pipe.Skip("TEST")
	})
	err := wrapped(ctx)
	assert.NoError(t, err)
}
