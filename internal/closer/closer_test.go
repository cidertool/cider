package closer

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

var errAlreadyClosed = errors.New("already closed")

func TestClose(t *testing.T) {
	var expectingErr bool

	onCloseErr = func(err error) {
		if !expectingErr {
			assert.FailNow(t, err.Error())
		}
	}

	r := &resource{name: "a"}
	Close(r)

	r = &resource{name: "b"}
	Close(r)

	expectingErr = true

	Close(r)
}

type resource struct {
	name   string
	closed bool
}

func (r *resource) Close() error {
	if r.closed {
		return errAlreadyClosed
	}

	r.closed = true

	return nil
}
