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

package client

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/cidertool/asc-go/asc"
	"github.com/stretchr/testify/assert"
)

// Test GetApp

func TestGetApp_Happy(t *testing.T) {
	expectedBundleID := "com.app.bundleid"
	ctx, client := newTestContext(response{
		Response: asc.AppsResponse{
			Data: []asc.App{
				{
					Attributes: &asc.AppAttributes{
						BundleID: &expectedBundleID,
					},
				},
			},
		},
	})

	defer ctx.Close()

	app, err := client.GetAppForBundleID(ctx.Context, expectedBundleID)
	assert.NoError(t, err)
	assert.NotNil(t, app)
	assert.Equal(t, expectedBundleID, *app.Attributes.BundleID)
}

func TestGetApp_Err(t *testing.T) {
	ctx, client := newTestContext(response{
		StatusCode:  http.StatusNotFound,
		RawResponse: `{}`,
	})
	defer ctx.Close()

	app, err := client.GetAppForBundleID(ctx.Context, "com.app.bundleid")
	assert.Error(t, err)
	assert.Nil(t, app)
}

func TestGetApp_ErrNoApps(t *testing.T) {
	ctx, client := newTestContext(response{
		RawResponse: `{"data":[]}`,
	})
	defer ctx.Close()

	app, err := client.GetAppForBundleID(ctx.Context, "com.app.bundleid")
	assert.Error(t, err)
	assert.Nil(t, app)
}

// Test GetAppInfo

func TestGetAppInfo_Happy(t *testing.T) {
	expectedState := asc.AppStoreVersionStatePrepareForSubmission
	app := asc.App{}
	ctx, client := newTestContext(response{
		Response: asc.AppInfosResponse{
			Data: []asc.AppInfo{
				{},
				{
					Attributes: &asc.AppInfoAttributes{},
				},
				{
					Attributes: &asc.AppInfoAttributes{
						AppStoreState: &expectedState,
					},
				},
			},
		},
	})

	defer ctx.Close()

	info, err := client.GetAppInfo(ctx.Context, app.ID)
	assert.NoError(t, err)
	assert.NotNil(t, info)
	assert.Equal(t, expectedState, *info.Attributes.AppStoreState)
}

func TestGetAppInfo_Err(t *testing.T) {
	app := asc.App{}
	ctx, client := newTestContext(response{
		StatusCode:  http.StatusNotFound,
		RawResponse: `{}`,
	})

	defer ctx.Close()

	info, err := client.GetAppInfo(ctx.Context, app.ID)
	assert.Error(t, err)
	assert.Nil(t, info)
}

func TestGetAppInfo_ErrNoData(t *testing.T) {
	app := asc.App{}
	ctx, client := newTestContext(response{
		RawResponse: `{}`,
	})

	defer ctx.Close()

	info, err := client.GetAppInfo(ctx.Context, app.ID)
	assert.Error(t, err)
	assert.Nil(t, info)
}

// Test GetBuild

const (
	testGetBuildVersion = "1.0"
)

func TestGetBuild_Happy(t *testing.T) {
	expectedProcessingState := validProcessingState
	app := asc.App{
		Attributes: &asc.AppAttributes{
			BundleID: asc.String("com.app.bundleid"),
		},
	}

	ctx, client := newTestContext(response{
		Response: asc.BuildsResponse{
			Data: []asc.Build{
				{
					Attributes: &asc.BuildAttributes{
						ProcessingState: &expectedProcessingState,
					},
				},
			},
		},
	})
	defer ctx.Close()

	ctx.Context.Version = testGetBuildVersion
	build, err := client.GetBuild(ctx.Context, &app)
	assert.NoError(t, err)
	assert.NotNil(t, build)
	assert.Equal(t, expectedProcessingState, *build.Attributes.ProcessingState)
}

func TestGetBuild_HappyOverrideBuild(t *testing.T) {
	expectedProcessingState := validProcessingState
	app := asc.App{
		Attributes: &asc.AppAttributes{
			BundleID: asc.String("com.app.bundleid"),
		},
	}

	ctx, client := newTestContext(response{
		Response: asc.BuildsResponse{
			Data: []asc.Build{
				{
					Attributes: &asc.BuildAttributes{
						ProcessingState: &expectedProcessingState,
					},
				},
			},
		},
	})
	defer ctx.Close()

	ctx.Context.Version = testGetBuildVersion
	ctx.Context.Build = "3"
	build, err := client.GetBuild(ctx.Context, &app)
	assert.NoError(t, err)
	assert.NotNil(t, build)
	assert.Equal(t, expectedProcessingState, *build.Attributes.ProcessingState)
}

func TestGetBuild_ErrNoVersion(t *testing.T) {
	app := asc.App{
		Attributes: &asc.AppAttributes{
			BundleID: asc.String("com.app.bundleid"),
		},
	}

	ctx, client := newTestContext(response{
		RawResponse: `{"data":[]}`,
	})
	defer ctx.Close()

	build, err := client.GetBuild(ctx.Context, &app)
	assert.Error(t, err)
	assert.Equal(t, "no version provided to lookup build with", err.Error())
	assert.Nil(t, build)
}

func TestGetBuild_Err(t *testing.T) {
	app := asc.App{
		Attributes: &asc.AppAttributes{
			BundleID: asc.String("com.app.bundleid"),
		},
	}

	ctx, client := newTestContext(response{
		StatusCode:  http.StatusNotFound,
		RawResponse: `{}`,
	})
	defer ctx.Close()

	ctx.Context.Version = testGetBuildVersion
	build, err := client.GetBuild(ctx.Context, &app)
	assert.Error(t, err)
	assert.Equal(t, fmt.Sprintf("build not found matching app=com.app.bundleid, version=1.0: GET %s/v1/builds: 404\n", ctx.server.URL), err.Error())
	assert.Nil(t, build)
}

func TestGetBuild_ErrNoBuilds(t *testing.T) {
	app := asc.App{
		Attributes: &asc.AppAttributes{
			BundleID: asc.String("com.app.bundleid"),
		},
	}

	ctx, client := newTestContext(response{
		RawResponse: `{"data":[]}`,
	})
	defer ctx.Close()

	ctx.Context.Version = testGetBuildVersion
	build, err := client.GetBuild(ctx.Context, &app)
	assert.Error(t, err)
	assert.Equal(t, "build not found matching app=com.app.bundleid, version=1.0", err.Error())
	assert.Nil(t, build)
}

func TestGetBuild_ErrNoAttributes(t *testing.T) {
	app := asc.App{
		Attributes: &asc.AppAttributes{
			BundleID: asc.String("com.app.bundleid"),
		},
	}

	ctx, client := newTestContext(response{
		RawResponse: `{"data":[{}]}`,
	})
	defer ctx.Close()

	ctx.Context.Version = testGetBuildVersion
	build, err := client.GetBuild(ctx.Context, &app)
	assert.Error(t, err)
	assert.Equal(t, "build  has no attributes", err.Error())
	assert.Nil(t, build)
}

func TestGetBuild_ErrNoProcessingState(t *testing.T) {
	app := asc.App{
		Attributes: &asc.AppAttributes{
			BundleID: asc.String("com.app.bundleid"),
		},
	}

	ctx, client := newTestContext(response{
		RawResponse: `{"data":[{"attributes":{}}]}`,
	})
	defer ctx.Close()

	ctx.Context.Version = testGetBuildVersion
	build, err := client.GetBuild(ctx.Context, &app)
	assert.Error(t, err)
	assert.Equal(t, "build  has no processing state", err.Error())
	assert.Nil(t, build)
}

func TestGetBuild_ErrInvalidProcessingState(t *testing.T) {
	app := asc.App{
		Attributes: &asc.AppAttributes{
			BundleID: asc.String("com.app.bundleid"),
		},
	}

	ctx, client := newTestContext(response{
		RawResponse: `{"data":[{"attributes":{"processingState":"PROCESSING"}}]}`,
	})
	defer ctx.Close()

	ctx.Context.Version = testGetBuildVersion
	build, err := client.GetBuild(ctx.Context, &app)
	assert.Error(t, err)
	assert.Equal(t, "latest build  has a processing state of PROCESSING. it would be dangerous to proceed", err.Error())
	assert.Nil(t, build)
}

// Test ReleaseForAppIsInitial

func TestReleaseForAppIsInitial_HappyInitial(t *testing.T) {
	app := asc.App{}
	ctx, client := newTestContext(response{
		Response: asc.AppStoreVersionsResponse{
			Data: []asc.AppStoreVersion{
				{},
			},
		},
	})

	defer ctx.Close()

	initial, err := client.ReleaseForAppIsInitial(ctx.Context, app.ID)
	assert.NoError(t, err)
	assert.True(t, initial)
}

func TestReleaseForAppIsInitial_Err(t *testing.T) {
	app := asc.App{}
	ctx, client := newTestContext(response{
		StatusCode:  http.StatusNotFound,
		RawResponse: `{}`,
	})

	defer ctx.Close()

	initial, err := client.ReleaseForAppIsInitial(ctx.Context, app.ID)
	assert.Error(t, err)
	assert.False(t, initial)
}

func TestReleaseForAppIsInitial_HappyNotInitial(t *testing.T) {
	app := asc.App{}
	ctx, client := newTestContext(response{
		Response: asc.AppStoreVersionsResponse{
			Data: []asc.AppStoreVersion{
				{},
				{},
				{},
			},
		},
	})

	defer ctx.Close()

	initial, err := client.ReleaseForAppIsInitial(ctx.Context, app.ID)
	assert.NoError(t, err)
	assert.False(t, initial)
}
