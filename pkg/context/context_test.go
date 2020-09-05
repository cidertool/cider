package context

import (
	"os"
	"testing"
	"time"

	"github.com/aaronsky/applereleaser/pkg/config"
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

func TestNewCredentials(t *testing.T) {
	privateKey := []byte(`
-----BEGIN PRIVATE KEY-----
MIGHAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBG0wawIBAQQgTHOfkv1Dj2Yp8hyT
cqY/BRRnYsoLzQoT04EW9d57dd+hRANCAAReUxyqGRpXQI2Fe833j7gQTnZ002VO
FSEpv2QUFs0+dXz04SWVmmzFErM0/iQyCYom0V1IMOWgV/8xvFN6+AeX
-----END PRIVATE KEY-----
`)
	cred, err := NewCredentials("kid", "iss", privateKey)
	assert.NoError(t, err)
	assert.NotNil(t, cred)
	assert.NotNil(t, cred.AuthTransport)
	cred, err = NewCredentials("kid", "iss", []byte("nothing"))
	assert.Error(t, err)
}

func TestEnv(t *testing.T) {
	var env = Env{"DOG": "FRIEND"}
	anotherEnv := env.Copy()
	assert.Equal(t, env, anotherEnv)
	assert.NotSame(t, &env, &anotherEnv)
	assert.Equal(t, []string{"DOG=FRIEND"}, env.Strings())
}
