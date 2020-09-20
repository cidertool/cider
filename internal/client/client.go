// Package client provides a full-featured App Store Connect API client
package client

import (
	"fmt"
	"strings"

	"github.com/cidertool/asc-go/asc"
	"github.com/cidertool/cider/pkg/config"
	"github.com/cidertool/cider/pkg/context"
)

const (
	validProcessingState = "VALID"
)

// Client is an abstraction of an App Store Connect API client's functionality.
type Client interface {
	// GetAppForBundleID returns the App resource matching the given bundle ID
	GetAppForBundleID(ctx *context.Context, bundleID string) (*asc.App, error)
	GetAppInfo(ctx *context.Context, appID string) (*asc.AppInfo, error)
	// GetBuild returns the Build resource for the given app, depending on the value set for
	// ctx.Build. Returns an error if the selected build is still processing.
	GetBuild(ctx *context.Context, app *asc.App) (*asc.Build, error)
	// ReleaseForAppIsInitial returns true if the App resource has never released before,
	// i.e. has one or less associated App Store Version relationships.
	ReleaseForAppIsInitial(ctx *context.Context, appID string) (bool, error)

	// Testflight

	// UpdateBetaAppLocalizations updates an App's beta app localizations, and creates any new ones that do not exist.
	// It will not delete or update any locales that are associated with the app but are not configured in cider.
	UpdateBetaAppLocalizations(ctx *context.Context, appID string, config config.TestflightLocalizations) error
	// UpdateBetaBuildDetails updates an App's beta build details, or creates new ones if they do not yet exist.
	UpdateBetaBuildDetails(ctx *context.Context, buildID string, config config.Testflight) error
	// UpdateBetaBuildLocalizations updates an App's beta build localizations, and creates any new ones that do not exist.
	// It will not delete or update any locales that are associated with the app but are not configured in cider.
	UpdateBetaBuildLocalizations(ctx *context.Context, buildID string, config config.TestflightLocalizations) error
	// UpdateBetaLicenseAgreement updates an App's beta license agreement, or creates a new one if one does not yet exist.
	UpdateBetaLicenseAgreement(ctx *context.Context, appID string, config config.Testflight) error
	AssignBetaGroups(ctx *context.Context, appID string, buildID string, groups []config.BetaGroup) error
	AssignBetaTesters(ctx *context.Context, appID string, buildID string, testers []config.BetaTester) error
	// UpdateBetaReviewDetails updates an App's beta review details, or creates new ones if they do not yet exist.
	UpdateBetaReviewDetails(ctx *context.Context, appID string, config config.ReviewDetails) error
	// SubmitBetaApp submits the given beta build for review
	SubmitBetaApp(ctx *context.Context, buildID string) error

	// App Store

	UpdateApp(ctx *context.Context, appID string, appInfoID string, versionID string, config config.App) error
	UpdateAppLocalizations(ctx *context.Context, appID string, config config.AppLocalizations) error
	CreateVersionIfNeeded(ctx *context.Context, appID string, buildID string, config config.Version) (*asc.AppStoreVersion, error)
	UpdateVersionLocalizations(ctx *context.Context, versionID string, config config.VersionLocalizations) error
	UpdateIDFADeclaration(ctx *context.Context, versionID string, config config.IDFADeclaration) error
	UploadRoutingCoverage(ctx *context.Context, versionID string, config config.File) error
	// UpdateReviewDetails updates an App's review details, or creates new ones if they do not yet exist.
	UpdateReviewDetails(ctx *context.Context, versionID string, config config.ReviewDetails) error
	EnablePhasedRelease(ctx *context.Context, versionID string) error
	// SubmitApp submits the given app store version for review
	SubmitApp(ctx *context.Context, versionID string) error

	Project() (*config.Project, error)
}

// New returns a new Client.
func New(ctx *context.Context) Client {
	client := asc.NewClient(ctx.Credentials.Client())
	return &ascClient{client: client}
}

type ascClient struct {
	client *asc.Client
}

func (c *ascClient) GetAppForBundleID(ctx *context.Context, bundleID string) (*asc.App, error) {
	resp, _, err := c.client.Apps.ListApps(ctx, &asc.ListAppsQuery{
		FilterBundleID: []string{bundleID},
	})
	if err != nil {
		return nil, fmt.Errorf("app not found matching %s: %w", bundleID, err)
	} else if len(resp.Data) == 0 {
		return nil, fmt.Errorf("app not found matching %s", bundleID)
	}
	return &resp.Data[0], nil
}

func (c *ascClient) GetAppInfo(ctx *context.Context, appID string) (*asc.AppInfo, error) {
	resp, _, err := c.client.Apps.ListAppInfosForApp(ctx, appID, nil)
	if err != nil {
		return nil, err
	}
	for _, info := range resp.Data {
		if info.Attributes == nil {
			continue
		} else if info.Attributes.AppStoreState == nil {
			continue
		}
		state := *info.Attributes.AppStoreState
		if state == asc.AppStoreVersionStatePrepareForSubmission {
			return &info, nil
		}
	}

	return nil, fmt.Errorf("app info not found matching %s", appID)
}

type errBuildNotFound struct {
	AppID         string
	BuildVersion  string
	VersionString string
	InnerErr      error
}

func (e errBuildNotFound) Error() string {
	str := strings.Builder{}
	str.WriteString("build not found")
	listingFields := false
	if e.AppID != "" {
		if !listingFields {
			str.WriteString(" matching ")
		}
		str.WriteString(fmt.Sprintf("app=%s", e.AppID))
		listingFields = true
	}
	if e.VersionString != "" {
		if listingFields {
			str.WriteString(", ")
		} else {
			str.WriteString(" matching ")
		}
		str.WriteString(fmt.Sprintf("version=%s", e.VersionString))
		listingFields = true
	}
	if e.BuildVersion != "" {
		if listingFields {
			str.WriteString(", ")
		} else {
			str.WriteString(" matching ")
		}
		str.WriteString(fmt.Sprintf("build=%s", e.BuildVersion))
	}
	if e.InnerErr != nil {
		str.WriteString(fmt.Errorf(": %w", e.InnerErr).Error())
	}
	return str.String()
}

func (c *ascClient) GetBuild(ctx *context.Context, app *asc.App) (*asc.Build, error) {
	if ctx.Version == "" {
		return nil, fmt.Errorf("no version provided to lookup build with")
	}
	query := asc.ListBuildsQuery{
		FilterApp:                      []string{app.ID},
		FilterPreReleaseVersionVersion: []string{ctx.Version},
	}
	if ctx.Build != "" {
		query.FilterVersion = []string{ctx.Build}
	}
	resp, _, err := c.client.Builds.ListBuilds(ctx, &query)
	if err != nil || len(resp.Data) == 0 {
		return nil, errBuildNotFound{
			AppID:         *app.Attributes.BundleID,
			VersionString: ctx.Version,
			BuildVersion:  ctx.Build,
			InnerErr:      err,
		}
	}
	build := resp.Data[0]
	if build.Attributes == nil {
		return nil, fmt.Errorf("build %s has no attributes", build.ID)
	}
	if build.Attributes.ProcessingState == nil {
		return nil, fmt.Errorf("build %s has no processing state", build.ID)
	}
	if *build.Attributes.ProcessingState != validProcessingState {
		return nil, fmt.Errorf("latest build %s has a processing state of %s. it would be dangerous to proceed", build.ID, *build.Attributes.ProcessingState)
	}
	return &build, nil
}

func (c *ascClient) ReleaseForAppIsInitial(ctx *context.Context, appID string) (bool, error) {
	resp, _, err := c.client.Apps.ListAppStoreVersionsForApp(ctx, appID, nil)
	if err != nil {
		return false, err
	}
	return len(resp.Data) <= 1, nil
}
