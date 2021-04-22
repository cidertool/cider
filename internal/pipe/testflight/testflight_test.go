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

package testflight

import (
	"testing"

	"github.com/cidertool/cider/internal/client/clienttest"
	"github.com/cidertool/cider/internal/pipe"
	"github.com/cidertool/cider/pkg/config"
	"github.com/cidertool/cider/pkg/context"
	"github.com/stretchr/testify/assert"
)

func TestTestflight_Happy(t *testing.T) {
	t.Parallel()

	ctx := context.New(config.Project{
		"TEST": {
			BundleID: "com.test.TEST",
			Testflight: config.Testflight{
				ReviewDetails: &config.ReviewDetails{
					Contact: &config.ContactPerson{
						Email:     "test@example.com",
						FirstName: "Person",
						LastName:  "Personson",
						Phone:     "1555555555",
					},
					DemoAccount: &config.DemoAccount{},
					Notes:       "TEST",
					Attachments: []config.File{
						{Path: "TEST"},
					},
				},
			},
		},
	})
	ctx.AppsToRelease = []string{"TEST"}

	p := Pipe{}
	p.Client = &clienttest.Client{}

	assert.Equal(t, "committing to testflight", p.String())

	err := p.Publish(ctx)
	assert.NoError(t, err)
}

func TestTestflight_Happy_Skips(t *testing.T) {
	t.Parallel()

	ctx := context.New(config.Project{
		"TEST": {
			BundleID: "com.test.TEST",
		},
	})
	ctx.AppsToRelease = []string{"TEST"}
	ctx.SkipUpdateMetadata = true
	ctx.SkipSubmit = true

	p := Pipe{}
	p.Client = &clienttest.Client{}

	err := p.Publish(ctx)
	assert.EqualError(t, err, pipe.ErrSkipSubmitEnabled.Error())
}

func TestTestflight_Happy_NoApps(t *testing.T) {
	t.Parallel()

	ctx := context.New(config.Project{})
	ctx.Credentials = &clienttest.Credentials{}

	p := Pipe{}

	err := p.Publish(ctx)
	assert.NoError(t, err)
}
