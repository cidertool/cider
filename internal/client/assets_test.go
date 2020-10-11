package client

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/cidertool/asc-go/asc"
	"github.com/cidertool/cider/pkg/config"
	"github.com/stretchr/testify/assert"
)

// Test UploadRoutingCoverage

func TestUploadRoutingCoverage_Happy(t *testing.T) {
	ctx, client := newTestContext(
		response{
			Response: asc.RoutingAppCoverageResponse{
				Data: asc.RoutingAppCoverage{
					Attributes: &asc.RoutingAppCoverageAttributes{
						UploadOperations: []asc.UploadOperation{},
					},
				},
			},
		},
		response{
			Response: asc.RoutingAppCoverageResponse{
				Data: asc.RoutingAppCoverage{
					Attributes: &asc.RoutingAppCoverageAttributes{
						UploadOperations: []asc.UploadOperation{},
					},
				},
			},
		},
		response{
			Response: asc.RoutingAppCoverageResponse{
				Data: asc.RoutingAppCoverage{
					Attributes: &asc.RoutingAppCoverageAttributes{
						UploadOperations: []asc.UploadOperation{},
					},
				},
			},
		},
		response{
			Response: asc.RoutingAppCoverageResponse{
				Data: asc.RoutingAppCoverage{
					Attributes: &asc.RoutingAppCoverageAttributes{
						UploadOperations: []asc.UploadOperation{},
					},
				},
			},
		},
	)
	defer ctx.Close()

	var path = filepath.Join(t.TempDir(), "TEST")
	err := ioutil.WriteFile(path, []byte("TEST"), 0600)
	assert.NoError(t, err)

	err = client.UploadRoutingCoverage(ctx.Context, "TEST", config.File{
		Path: path,
	})
	assert.NoError(t, err)
}
