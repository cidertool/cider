package context

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInterruptOK(t *testing.T) {
	assert.NoError(t, NewInterrupt().Run(context.Background(), func() error {
		return nil
	}))
}

func TestInterruptErrors(t *testing.T) {
	var err = errors.New("some error")

	assert.EqualError(t, NewInterrupt().Run(context.Background(), func() error {
		return err
	}), err.Error())
}
