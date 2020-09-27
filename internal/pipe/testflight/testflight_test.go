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
	ctx := context.New(config.Project{})
	ctx.Credentials = &clienttest.Credentials{}
	p := Pipe{}

	err := p.Publish(ctx)
	assert.NoError(t, err)
}
