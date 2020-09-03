package client

import (
	"crypto/md5"
	"fmt"
	"io"

	"github.com/aaronsky/applereleaser/pkg/config"
	"github.com/aaronsky/applereleaser/pkg/context"
	"github.com/aaronsky/asc-go/asc"
)

// Client is an abstraction of an App Store Connect API client's functionality
type Client interface {
	GetAppForBundleID(ctx *context.Context, id string) (*asc.App, error)
	GetRelevantBuild(ctx *context.Context, app *asc.App) (*asc.Build, error)
	UpdateBetaAppLocalizations(ctx *context.Context, app *asc.App, config config.TestflightLocalizations) error
	UpdateBetaBuildDetails(ctx *context.Context, build *asc.Build, config config.TestflightForApp) error
	UpdateBetaBuildLocalizations(ctx *context.Context, build *asc.Build, config config.TestflightLocalizations) error
	UpdateBetaLicenseAgreement(ctx *context.Context, app *asc.App, config config.TestflightForApp) error
	AssignBetaGroups(ctx *context.Context, build *asc.Build, groups []string) error
	AssignBetaTesters(ctx *context.Context, build *asc.Build, testers []config.BetaTester) error
	UpdateBetaReviewDetails(ctx *context.Context, app *asc.App, config config.ReviewDetails) error
	SubmitBetaApp(ctx *context.Context, build *asc.Build) error
	UpdateAppLocalizations(ctx *context.Context, app *asc.App, config config.AppLocalizations) error
	CreateVersionIfNeeded(ctx *context.Context, app *asc.App, build *asc.Build, config config.Version) (*asc.AppStoreVersion, error)
	UpdateVersionLocalizations(ctx *context.Context, version *asc.AppStoreVersion, config config.VersionLocalizations) error
	UpdateIDFADeclaration(ctx *context.Context, version *asc.AppStoreVersion, config config.IDFADeclaration) error
	UploadRoutingCoverage(ctx *context.Context, version *asc.AppStoreVersion, config config.File) error
	UpdatePreviewSets(ctx *context.Context, previewSets []asc.AppPreviewSet, appStoreVersionLocalizationID string, config config.PreviewSets) error
	UpdateScreenshotSets(ctx *context.Context, screenshotSets []asc.AppScreenshotSet, appStoreVersionLocalizationID string, config config.ScreenshotSets) error
	UpdateReviewDetails(ctx *context.Context, version *asc.AppStoreVersion, config config.ReviewDetails) error
	SubmitApp(ctx *context.Context, version *asc.AppStoreVersion) error
}

// New returns a new Client
func New(ctx *context.Context) Client {
	client := asc.NewClient(ctx.Credentials.Client())
	return &ascClient{client: client}
}

type ascClient struct {
	client *asc.Client
}

func (c *ascClient) GetAppForBundleID(ctx *context.Context, id string) (*asc.App, error) {
	resp, _, err := c.client.Apps.ListApps(ctx, &asc.ListAppsQuery{
		FilterBundleID: []string{id},
	})
	if err != nil {
		return nil, fmt.Errorf("app not found matching %s: %w", id, err)
	} else if len(resp.Data) == 0 {
		return nil, fmt.Errorf("app not found matching %s", id)
	}
	return &resp.Data[0], nil
}

func (c *ascClient) GetRelevantBuild(ctx *context.Context, app *asc.App) (*asc.Build, error) {
	if ctx.Version == "" {
		return nil, fmt.Errorf("no version provided to lookup build with")
	}
	resp, _, err := c.client.Builds.ListBuilds(ctx, &asc.ListBuildsQuery{
		FilterApp:                      []string{app.ID},
		FilterPreReleaseVersionVersion: []string{ctx.Version},
	})
	if err != nil {
		return nil, fmt.Errorf("build not found matching app %s and version %s: %w", *app.Attributes.BundleID, ctx.Version, err)
	} else if len(resp.Data) == 0 {
		return nil, fmt.Errorf("build not found matching app %s and version %s", *app.Attributes.BundleID, ctx.Version)
	}
	build := resp.Data[0]
	if build.Attributes == nil {
		return nil, fmt.Errorf("build %s has no attributes", build.ID)
	} else if build.Attributes.ProcessingState == nil {
		return nil, fmt.Errorf("build %s has no processing state", build.ID)
	} else if *build.Attributes.ProcessingState != "VALID" {
		return nil, fmt.Errorf("latest build %s has a processing state of %s. it would be dangerous to proceed", build.ID, *build.Attributes.ProcessingState)
	}
	return &build, nil
}

func md5Checksum(f io.Reader) (string, error) {
	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
