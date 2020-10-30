package context

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

var errTestError = errors.New("TEST")

func TestInterruptOK(t *testing.T) {
	assert.NoError(t, NewInterrupt().Run(context.Background(), func() error {
		return nil
	}))
}

func TestInterruptErrors(t *testing.T) {
	assert.EqualError(t, NewInterrupt().Run(context.Background(), func() error {
		return errTestError
	}), errTestError.Error())
}
