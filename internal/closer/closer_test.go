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
