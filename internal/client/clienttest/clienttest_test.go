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

package clienttest_test

import (
	"testing"

	"github.com/cidertool/cider/internal/client/clienttest"
	"github.com/cidertool/cider/pkg/config"
	"github.com/cidertool/cider/pkg/context"
	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	t.Parallel()

	ctx := context.New(config.Project{})
	c := clienttest.Client{}

	app, err := c.GetAppForBundleID(ctx, "TEST")
	assert.NoError(t, err)
	assert.NotNil(t, app)

	info, err := c.GetAppInfo(ctx, "TEST")
	assert.NoError(t, err)
	assert.NotNil(t, info)

	build, err := c.GetBuild(ctx, nil)
	assert.NoError(t, err)
	assert.NotNil(t, build)

	initial, err := c.ReleaseForAppIsInitial(ctx, "TEST")
	assert.NoError(t, err)
	assert.False(t, initial)

	err = c.UpdateBetaAppLocalizations(ctx, "TEST", config.TestflightLocalizations{})
	assert.NoError(t, err)

	err = c.UpdateBetaBuildDetails(ctx, "TEST", config.Testflight{})
	assert.NoError(t, err)

	err = c.UpdateBetaBuildLocalizations(ctx, "TEST", config.TestflightLocalizations{})
	assert.NoError(t, err)

	err = c.UpdateBetaLicenseAgreement(ctx, "TEST", config.Testflight{})
	assert.NoError(t, err)

	err = c.AssignBetaGroups(ctx, "TEST", "TEST", []config.BetaGroup{})
	assert.NoError(t, err)

	err = c.AssignBetaTesters(ctx, "TEST", "TEST", []config.BetaTester{})
	assert.NoError(t, err)

	err = c.UpdateBetaReviewDetails(ctx, "TEST", config.ReviewDetails{})
	assert.NoError(t, err)

	err = c.SubmitBetaApp(ctx, "TEST")
	assert.NoError(t, err)

	err = c.UpdateApp(ctx, "TEST", "TEST", "TEST", config.App{})
	assert.NoError(t, err)

	err = c.UpdateAppLocalizations(ctx, "TEST", config.AppLocalizations{})
	assert.NoError(t, err)

	version, err := c.CreateVersionIfNeeded(ctx, "TEST", "TEST", config.Version{})
	assert.NoError(t, err)
	assert.NotNil(t, version)

	err = c.UpdateVersionLocalizations(ctx, "TEST", config.VersionLocalizations{})
	assert.NoError(t, err)

	err = c.UpdateIDFADeclaration(ctx, "TEST", config.IDFADeclaration{})
	assert.NoError(t, err)

	err = c.UploadRoutingCoverage(ctx, "TEST", config.File{})
	assert.NoError(t, err)

	err = c.UpdateReviewDetails(ctx, "TEST", config.ReviewDetails{})
	assert.NoError(t, err)

	err = c.EnablePhasedRelease(ctx, "TEST")
	assert.NoError(t, err)

	err = c.SubmitApp(ctx, "TEST")
	assert.NoError(t, err)

	proj, err := c.Project()
	assert.NoError(t, err)
	assert.NotNil(t, proj)
}

func TestCredentials(t *testing.T) {
	t.Parallel()

	c := clienttest.Credentials{}
	assert.NotNil(t, c.Client())
}
