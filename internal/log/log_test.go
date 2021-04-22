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

package log

import (
	"testing"

	alog "github.com/apex/log"
	"github.com/fatih/color"
	"github.com/stretchr/testify/assert"
)

func TestLog(t *testing.T) {
	t.Parallel()

	log := New()

	log.SetPadding(4)

	log.SetColorMode(false)
	assert.False(t, color.NoColor)
	log.SetColorMode(true)
	assert.True(t, color.NoColor)

	log.SetDebug(true)
	assert.Equal(t, log.Level, alog.DebugLevel)
	log.SetDebug(false)
	assert.Equal(t, log.Level, alog.InfoLevel)
}
