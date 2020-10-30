package publish

import (
	"testing"

	"github.com/cidertool/cider/internal/client/clienttest"
	"github.com/cidertool/cider/internal/pipe"
	"github.com/cidertool/cider/pkg/config"
	"github.com/cidertool/cider/pkg/context"
	"github.com/stretchr/testify/assert"
)

func TestPublish_String(t *testing.T) {
	p := Pipe{}
	assert.Equal(t, "publishing from app store connect", p.String())
}

func TestPublish_Happy_Testflight(t *testing.T) {
	ctx := context.New(config.Project{
		"TEST": {},
	})
	ctx.AppsToRelease = []string{"TEST"}
	ctx.Credentials = &clienttest.Credentials{}
	ctx.PublishMode = context.PublishModeTestflight
	p := Pipe{}
	p.client = &clienttest.Client{}

	err := p.Run(ctx)
	assert.NoError(t, err)
}

func TestPublish_Happy_Store(t *testing.T) {
	ctx := context.New(config.Project{
		"TEST": {},
	})
	ctx.AppsToRelease = []string{"TEST"}
	ctx.Credentials = &clienttest.Credentials{}
	ctx.PublishMode = context.PublishModeAppStore
	p := Pipe{}
	p.client = &clienttest.Client{}

	err := p.Run(ctx)
	assert.NoError(t, err)
}

func TestPublish_Happy_NoApps(t *testing.T) {
	ctx := context.New(config.Project{})
	p := Pipe{}

	err := p.Run(ctx)
	assert.EqualError(t, err, pipe.ErrSkipNoAppsToPublish.Error())
}

func TestPublish_Err_NoPublishMode(t *testing.T) {
	ctx := context.New(config.Project{})
	ctx.AppsToRelease = []string{"TEST"}
	p := Pipe{}

	err := p.Run(ctx)
	assert.EqualError(t, err, errUnsupportedPublishMode{context.PublishMode("")}.Error())
}

func TestPublish_Err_AppMismatchTestflight(t *testing.T) {
	ctx := context.New(config.Project{
		"TEST_": {},
	})
	ctx.AppsToRelease = []string{"_TEST"}
	ctx.Credentials = &clienttest.Credentials{}
	ctx.PublishMode = context.PublishModeTestflight
	p := Pipe{}
	p.client = &clienttest.Client{}

	err := p.Run(ctx)
	assert.Error(t, err)
}

func TestPublish_Err_AppMismatchStore(t *testing.T) {
	ctx := context.New(config.Project{
		"TEST_": {},
	})
	ctx.AppsToRelease = []string{"_TEST"}
	ctx.Credentials = &clienttest.Credentials{}
	ctx.PublishMode = context.PublishModeAppStore
	p := Pipe{}
	p.client = &clienttest.Client{}

	err := p.Run(ctx)
	assert.Error(t, err)
}
