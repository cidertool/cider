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

// Package pipe declares utilities and errors for pipes
package pipe

import (
	"errors"
	"fmt"
)

// ErrSkipGitEnabled happens if --skip-git is set. It means that the part of a Piper that
// extracts metadata from the Git repository was not run.
var ErrSkipGitEnabled = Skip("inspecting git state is disabled")

// ErrSkipNoAppsToPublish happens when there are no apps in the configuration to publish
// or update metadata for. It will be raised when the configuration is effectively empty.
var ErrSkipNoAppsToPublish = Skip("no apps selected to publish")

// ErrSkipSubmitEnabled happens if --skip-submit is set.
// It means that the part of a Piper that submits to Apple for review was not run.
var ErrSkipSubmitEnabled = Skip("submission is disabled")

// ErrMissingApp happens when an app is selected in the interface that is not defined in the configuration.
type ErrMissingApp struct {
	Name string
}

func (e ErrMissingApp) Error() string {
	return fmt.Sprintf("no app defined in configuration matching the name %s", e.Name)
}

// IsSkip returns true if the error is an ErrSkip.
func IsSkip(err error) bool {
	var serr ErrSkip
	ok := errors.As(err, &serr)

	return ok
}

// ErrSkip occurs when a pipe is skipped for some reason.
type ErrSkip struct {
	reason string
}

// Error implements the error interface. returns the reason the pipe was skipped.
func (e ErrSkip) Error() string {
	return e.reason
}

// Skip skips this pipe with the given reason.
func Skip(reason string) ErrSkip {
	return ErrSkip{reason: reason}
}
