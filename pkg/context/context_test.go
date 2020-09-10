package context

import (
	"os"
	"testing"
	"time"

	"github.com/cidertool/cider/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assert.NoError(t, os.Setenv("TEST", "DOG"))
	ctx := New(config.Project{})
	assert.Equal(t, "DOG", ctx.Env["TEST"])
}

func TestNewWithTimeout(t *testing.T) {
	ctx, cancel := NewWithTimeout(config.Project{}, time.Second)
	assert.NotEmpty(t, ctx.Env)
	cancel()
	<-ctx.Done()
	assert.EqualError(t, ctx.Err(), `context canceled`)
}

func TestEnv(t *testing.T) {
	var env = Env{"DOG": "FRIEND"}
	anotherEnv := env.Copy()
	assert.Equal(t, env, anotherEnv)
	assert.NotSame(t, &env, &anotherEnv)
	assert.Equal(t, []string{"DOG=FRIEND"}, env.Strings())
}

func TestPublishMode(t *testing.T) {
	var mode PublishMode
	mode = PublishModeAppStore
	assert.Equal(t, "appstore", mode.String())
	assert.Equal(t, "{appstore,testflight}", mode.Type())
	mode = PublishModeTestflight
	assert.Equal(t, "testflight", mode.String())
	assert.Equal(t, "{appstore,testflight}", mode.Type())
	mode = PublishMode("bad")
	assert.Equal(t, "bad", mode.String())
	assert.Equal(t, "{appstore,testflight}", mode.Type())

	var err error
	err = mode.Set("appstore")
	assert.NoError(t, err)
	assert.Equal(t, PublishModeAppStore, mode)
	err = mode.Set("testflight")
	assert.NoError(t, err)
	assert.Equal(t, PublishModeTestflight, mode)
	err = mode.Set("bad")
	assert.Error(t, err)
}
