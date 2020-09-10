// Package env is a pipe that loads environment variables
package env

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/cidertool/cider/pkg/context"
)

// ErrMissingEnvVar indicates an error when a required variable is missing in the environment.
var ErrMissingEnvVar = errors.New("missing required environment variable")

// Pipe is a global hook pipe.
type Pipe struct{}

// String is the name of this pipe.
func (Pipe) String() string {
	return "loading environment variables"
}

// Run executes the hooks.
func (p Pipe) Run(ctx *context.Context) error {
	keyID, err := loadEnv("ASC_KEY_ID", true)
	if err != nil {
		return err
	}
	issuerID, err := loadEnv("ASC_ISSUER_ID", true)
	if err != nil {
		return err
	}
	privateKey, err := loadEnvFromPath("ASC_PRIVATE_KEY_PATH", true)
	if err != nil {
		return err
	}
	creds, err := context.NewCredentials(keyID, issuerID, []byte(privateKey))
	if err != nil {
		return err
	}
	ctx.Credentials = creds
	return nil
}

func loadEnv(env string, required bool) (string, error) {
	val := os.Getenv(env)
	if val == "" && required {
		return "", fmt.Errorf("key %s not found: %w", env, ErrMissingEnvVar)
	}
	return val, nil
}

func loadEnvFromPath(env string, required bool) (string, error) {
	val, err := loadEnv(env, required)
	if err != nil {
		return "", err
	}
	f, err := os.Open(filepath.Clean(val))
	if err != nil && required {
		return "", err
	} else if err != nil {
		return "", nil
	}
	bytes, err := ioutil.ReadAll(f)
	return string(bytes), err
}
