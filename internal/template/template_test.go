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

package template

import (
	"testing"

	"github.com/cidertool/cider/pkg/config"
	"github.com/cidertool/cider/pkg/context"
	"github.com/stretchr/testify/assert"
)

func TestTemplate(t *testing.T) {
	tmpl, err := New(context.New(config.Project{})).
		WithEnv(map[string]string{
			"DOG": "HAPPY",
			"CAT": "GRUMPY",
		}).
		WithShellEnv("HORSE=EVIL", "CAT=SPOOKY").
		WithFields(Fields{
			"customKey": 0,
		}).
		Apply(`My {{ .env.CAT }} cat fought {{ .customKey }} {{ .env.HORSE }} horses.`)
	assert.NoError(t, err)
	assert.Equal(t, "My SPOOKY cat fought 0 EVIL horses.", tmpl)
}

func TestInvalidTemplate(t *testing.T) {
	_, err := New(context.New(config.Project{})).Apply(`{{ .timestamp`)
	assert.Error(t, err)
}

func TestEmptyTemplate(t *testing.T) {
	tmpl, err := New(context.New(config.Project{})).Apply("")
	assert.NoError(t, err)
	assert.Empty(t, tmpl)
}
