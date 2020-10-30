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
	"fmt"
	"net/http"
	"time"

	"github.com/cidertool/asc-go/asc"
)

const twentyMinuteTokenLifetime = time.Minute * 20

// Credentials stores credentials used by clients.
type Credentials interface {
	Client() *http.Client
}

type credentials struct {
	*asc.AuthTransport
}

// NewCredentials returns a new store object for App Store Connect credentials.
func NewCredentials(keyID, issuerID string, privateKey []byte) (Credentials, error) {
	token, err := asc.NewTokenConfig(keyID, issuerID, twentyMinuteTokenLifetime, privateKey)
	if err != nil {
		err = fmt.Errorf("failed to authorize with App Store Connect: %w", err)
	}

	return credentials{token}, err
}
