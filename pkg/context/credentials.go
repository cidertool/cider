package context

import (
	"fmt"
	"net/http"
	"time"

	"github.com/aaronsky/asc-go/asc"
)

// Credentials stores credentials used by clients.
type Credentials interface {
	Client() *http.Client
}

type credentials struct {
	*asc.AuthTransport
}

// NewCredentials returns a new store object for App Store Connect credentials.
func NewCredentials(keyID, issuerID string, privateKey []byte) (Credentials, error) {
	token, err := asc.NewTokenConfig(keyID, issuerID, time.Minute*20, privateKey)
	if err != nil {
		err = fmt.Errorf("failed to authorize with App Store Connect: %w", err)
	}
	return credentials{token}, err
}
