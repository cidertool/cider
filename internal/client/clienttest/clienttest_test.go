package clienttest_test

import (
	"testing"

	"github.com/cidertool/cider/internal/client/clienttest"
	"github.com/cidertool/cider/pkg/config"
	"github.com/cidertool/cider/pkg/context"
	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
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
	c := clienttest.Credentials{}
	assert.NotNil(t, c.Client())
}
