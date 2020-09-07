package client

import (
	"net/http"
	"testing"

	"github.com/aaronsky/asc-go/asc"
	"github.com/stretchr/testify/assert"
)

func TestGetApp(t *testing.T) {
	expectedBundleID := "com.app.bundleid"

	ctx := newTestContext(response{
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
	client := New(ctx.Context)

	// happy
	app, err := client.GetAppForBundleID(ctx.Context, expectedBundleID)
	assert.NoError(t, err)
	assert.NotNil(t, app)
	assert.Equal(t, expectedBundleID, *app.Attributes.BundleID)

	// err raise err
	ctx.SetResponses(response{
		StatusCode:  http.StatusNotFound,
		RawResponse: `{}`,
	})
	app, err = client.GetAppForBundleID(ctx.Context, expectedBundleID)
	assert.Error(t, err)
	assert.Nil(t, app)

	// err no apps
	ctx.SetResponses(response{
		RawResponse: `{"data":[]}`,
	})
	app, err = client.GetAppForBundleID(ctx.Context, expectedBundleID)
	assert.Error(t, err)
	assert.Nil(t, app)
}

func TestGetAppInfo(t *testing.T) {
	expectedState := asc.AppStoreVersionStatePrepareForSubmission
	app := asc.App{}

	ctx := newTestContext(response{
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
	client := New(ctx.Context)

	// happy
	info, err := client.GetAppInfo(ctx.Context, &app)
	assert.NoError(t, err)
	assert.NotNil(t, info)
	assert.Equal(t, expectedState, *info.Attributes.AppStoreState)

	// err raise err
	ctx.SetResponses(response{
		StatusCode:  http.StatusNotFound,
		RawResponse: `{}`,
	})
	info, err = client.GetAppInfo(ctx.Context, &app)
	assert.Error(t, err)
	assert.Nil(t, info)

	// err no applicable data
	ctx.SetResponses(response{
		RawResponse: `{}`,
	})
	info, err = client.GetAppInfo(ctx.Context, &app)
	assert.Error(t, err)
	assert.Nil(t, info)
}

func TestGetRelevantBuild(t *testing.T) {
	expectedProcessingState := "VALID"
	app := asc.App{
		Attributes: &asc.AppAttributes{
			BundleID: asc.String("com.app.bundleid"),
		},
	}

	ctx := newTestContext(response{
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
	client := New(ctx.Context)

	// happy
	ctx.Context.Version = "1.0"
	build, err := client.GetRelevantBuild(ctx.Context, &app)
	assert.NoError(t, err)
	assert.NotNil(t, build)
	assert.Equal(t, expectedProcessingState, *build.Attributes.ProcessingState)

	// err no version
	ctx.Context.Version = ""
	ctx.SetResponses(response{
		RawResponse: `{"data":[]}`,
	})
	build, err = client.GetRelevantBuild(ctx.Context, &app)
	assert.Error(t, err)
	assert.Nil(t, build)
	ctx.Context.Version = "1.0"

	// err raise err
	ctx.SetResponses(response{
		StatusCode:  http.StatusNotFound,
		RawResponse: `{}`,
	})
	build, err = client.GetRelevantBuild(ctx.Context, &app)
	assert.Error(t, err)
	assert.Nil(t, build)

	// err no builds
	ctx.SetResponses(response{
		RawResponse: `{"data":[]}`,
	})
	build, err = client.GetRelevantBuild(ctx.Context, &app)
	assert.Error(t, err)
	assert.Nil(t, build)

	// err no attributes
	ctx.SetResponses(response{
		RawResponse: `{"data":[{}]}`,
	})
	build, err = client.GetRelevantBuild(ctx.Context, &app)
	assert.Error(t, err)
	assert.Nil(t, build)

	// err no processing state
	ctx.SetResponses(response{
		RawResponse: `{"data":[{"attributes":{}}]}`,
	})
	build, err = client.GetRelevantBuild(ctx.Context, &app)
	assert.Error(t, err)
	assert.Nil(t, build)

	// err invalid processing state
	ctx.SetResponses(response{
		RawResponse: `{"data":[{"attributes":{"processingState":"PROCESSING"}}]}`,
	})
	build, err = client.GetRelevantBuild(ctx.Context, &app)
	assert.Error(t, err)
	assert.Nil(t, build)
}

func TestReleaseForAppIsInitial(t *testing.T) {
	app := asc.App{}

	ctx := newTestContext(response{
		Response: asc.AppStoreVersionsResponse{
			Data: []asc.AppStoreVersion{
				{},
			},
		},
	})
	defer ctx.Close()
	client := New(ctx.Context)

	// happy
	initial, err := client.ReleaseForAppIsInitial(ctx.Context, &app)
	assert.NoError(t, err)
	assert.True(t, initial)

	// err raise err
	ctx.SetResponses(response{
		StatusCode:  http.StatusNotFound,
		RawResponse: `{}`,
	})
	initial, err = client.ReleaseForAppIsInitial(ctx.Context, &app)
	assert.Error(t, err)
	assert.False(t, initial)

	// not initial
	ctx.SetResponses(response{
		Response: asc.AppStoreVersionsResponse{
			Data: []asc.AppStoreVersion{
				{},
				{},
				{},
			},
		},
	})
	initial, err = client.ReleaseForAppIsInitial(ctx.Context, &app)
	assert.NoError(t, err)
	assert.False(t, initial)
}
