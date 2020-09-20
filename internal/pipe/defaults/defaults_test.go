package defaults

import (
	"errors"
	"testing"

	"github.com/cidertool/cider/internal/defaults"
	"github.com/cidertool/cider/pkg/config"
	"github.com/cidertool/cider/pkg/context"
	"github.com/stretchr/testify/assert"
)

func TestDefaults(t *testing.T) {
	ctx := context.New(config.Project{})
	pipe := Pipe{}
	var err error

	assert.Equal(t, "setting defaults", pipe.String())

	err = pipe.Run(ctx)
	assert.NoError(t, err)

	pipe.defaulters = []defaults.Defaulter{
		mockDefaulter{},
	}
	err = pipe.Run(ctx)
	assert.Error(t, err)
}

type mockDefaulter struct{}

func (d mockDefaulter) String() string {
	return ""
}

func (d mockDefaulter) Default(ctx *context.Context) error {
	return errors.New("TEST")
}
