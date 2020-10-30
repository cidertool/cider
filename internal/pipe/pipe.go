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
