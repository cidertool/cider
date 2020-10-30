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
