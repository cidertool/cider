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

package context

import (
	"os"
	"testing"
	"time"

	"github.com/cidertool/cider/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assert.NoError(t, os.Setenv("TEST", "DOG"))

	ctx := New(config.Project{})
	assert.Equal(t, "DOG", ctx.Env["TEST"])
}

func TestNewWithTimeout(t *testing.T) {
	ctx, cancel := NewWithTimeout(config.Project{}, time.Second)
	assert.NotEmpty(t, ctx.Env)
	cancel()
	<-ctx.Done()
	assert.EqualError(t, ctx.Err(), `context canceled`)
}

func TestEnv(t *testing.T) {
	var env = Env{"DOG": "FRIEND"}
	anotherEnv := env.Copy()
	assert.Equal(t, env, anotherEnv)
	assert.NotSame(t, &env, &anotherEnv)
	assert.Equal(t, []string{"DOG=FRIEND"}, env.Strings())
}

func TestPublishMode(t *testing.T) {
	var mode PublishMode
	mode = PublishModeAppStore
	assert.Equal(t, "appstore", mode.String())
	assert.Equal(t, "{appstore,testflight}", mode.Type())
	mode = PublishModeTestflight
	assert.Equal(t, "testflight", mode.String())
	assert.Equal(t, "{appstore,testflight}", mode.Type())
	mode = PublishMode("bad")
	assert.Equal(t, "bad", mode.String())
	assert.Equal(t, "{appstore,testflight}", mode.Type())

	var err error
	err = mode.Set("appstore")
	assert.NoError(t, err)
	assert.Equal(t, PublishModeAppStore, mode)
	err = mode.Set("testflight")
	assert.NoError(t, err)
	assert.Equal(t, PublishModeTestflight, mode)
	err = mode.Set("bad")
	assert.Error(t, err)
}
