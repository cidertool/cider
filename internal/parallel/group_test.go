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

package parallel

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var errTestError = errors.New("TEST")

func TestGroup(t *testing.T) {
	var g = New(4)

	var lock sync.Mutex

	var counter int

	for i := 0; i < 10; i++ {
		g.Go(func() error {
			time.Sleep(10 * time.Millisecond)
			lock.Lock()
			counter++
			lock.Unlock()
			return nil
		})
	}
	assert.NoError(t, g.Wait())
	assert.Equal(t, counter, 10)
}

func TestGroupOrder(t *testing.T) {
	var num = 10

	var g = New(1)

	var output = []int{}

	for i := 0; i < num; i++ {
		i := i

		g.Go(func() error {
			output = append(output, i)
			return nil
		})
	}
	assert.NoError(t, g.Wait())
	assert.Equal(t, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, output)
}

func TestGroupOrderError(t *testing.T) {
	var g = New(1)

	var output = []int{}

	for i := 0; i < 10; i++ {
		i := i

		g.Go(func() error {
			output = append(output, i)
			return errTestError
		})
	}
	assert.EqualError(t, g.Wait(), errTestError.Error())
	assert.Equal(t, []int{0}, output)
}
