package store

import (
	"testing"

	"github.com/cidertool/cider/internal/client/clienttest"
	"github.com/cidertool/cider/internal/pipe"
	"github.com/cidertool/cider/pkg/config"
	"github.com/cidertool/cider/pkg/context"
	"github.com/stretchr/testify/assert"
)

func TestStore_Happy(t *testing.T) {
	ctx := context.New(config.Project{
		"TEST": {
			BundleID: "com.test.TEST",
			Versions: config.Version{
				PhasedReleaseEnabled: true,
				IDFADeclaration: &config.IDFADeclaration{
					HonorsLimitedAdTracking: true,
				},
				RoutingCoverage: &config.File{
					Path: "TEST",
				},
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

	assert.Equal(t, "committing to app store", p.String())

	err := p.Publish(ctx)
	assert.NoError(t, err)
}

func TestStore_Happy_Skips(t *testing.T) {
	ctx := context.New(config.Project{
		"TEST": {
			BundleID: "com.test.TEST",
		},
	})
	ctx.AppsToRelease = []string{"TEST"}
	ctx.SkipUpdatePricing = true
	ctx.SkipUpdateMetadata = true
	ctx.SkipSubmit = true
	p := Pipe{}
	p.Client = &clienttest.Client{}

	err := p.Publish(ctx)
	assert.EqualError(t, err, pipe.ErrSkipSubmitEnabled.Error())
}

func TestStore_Happy_NoApps(t *testing.T) {
	ctx := context.New(config.Project{})
	ctx.Credentials = &clienttest.Credentials{}
	p := Pipe{}

	err := p.Publish(ctx)
	assert.NoError(t, err)
}
