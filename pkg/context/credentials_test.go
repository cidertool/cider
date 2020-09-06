package context

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
	assert.NotNil(t, cred.Client())
	_, err = NewCredentials("kid", "iss", []byte("nothing"))
	assert.Error(t, err)
}
