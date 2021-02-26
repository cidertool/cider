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

// Package env is a pipe that loads environment variables
package env

import (
	"errors"
	"fmt"
	"io"
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

	privateKey, err := loadEnv("ASC_PRIVATE_KEY", true)

	if err != nil {
		privateKey, err = loadEnvFromPath("ASC_PRIVATE_KEY_PATH", true)
		if err != nil {
			return err
		}
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

	bytes, err := io.ReadAll(f)

	return string(bytes), err
}
