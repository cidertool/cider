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

package semver

import (
	"testing"

	"github.com/cidertool/cider/pkg/config"
	"github.com/cidertool/cider/pkg/context"
	"github.com/stretchr/testify/assert"
)

func TestSemver(t *testing.T) {
	t.Parallel()

	ctx := context.New(config.Project{})
	pipe := Pipe{}

	var err error

	assert.Equal(t, "parsing version", pipe.String())

	err = pipe.Run(ctx)
	assert.Error(t, err)
	assert.Empty(t, ctx.Semver)

	ctx.Version = "1.0.1"
	err = pipe.Run(ctx)
	assert.NoError(t, err)
	assert.Equal(t, context.Semver{
		Major:      1,
		Minor:      0,
		Patch:      1,
		RawVersion: "1.0.1",
	}, ctx.Semver)

	ctx.Version = "1.1.1-patch90"
	err = pipe.Run(ctx)
	assert.NoError(t, err)
	assert.Equal(t, context.Semver{
		Major:      1,
		Minor:      1,
		Patch:      1,
		RawVersion: "1.1.1-patch90",
		Prerelease: "patch90",
	}, ctx.Semver)

	ctx.Version = "aa.ee.bb"
	ctx.Semver = context.Semver{}
	err = pipe.Run(ctx)
	assert.Error(t, err)
	assert.Empty(t, ctx.Semver)
}
