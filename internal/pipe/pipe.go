package pipe

// ErrSkipUpdateMetadataEnabled happens if --skip-update-metadata is set.
// It means that the part of a Piper that updates metadata in App Store Connect
// was not run.
var ErrSkipUpdateMetadataEnabled = Skip("updating metadata is disabled")

// ErrSkipSubmitEnabled happens if --skip-submit is set.
// It means that the part of a Piper that submits to Apple for review was not run.
var ErrSkipSubmitEnabled = Skip("submission is disabled")

// IsSkip returns true if the error is an ErrSkip.
func IsSkip(err error) bool {
	_, ok := err.(ErrSkip)
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
