package semver

import (
	"testing"

	"github.com/cidertool/cider/pkg/config"
	"github.com/cidertool/cider/pkg/context"
	"github.com/stretchr/testify/assert"
)

func TestSemver(t *testing.T) {
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
