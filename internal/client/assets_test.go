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
	"testing"

	"github.com/cidertool/asc-go/asc"
	"github.com/cidertool/cider/pkg/config"
	"github.com/stretchr/testify/assert"
)

// Test UploadRoutingCoverage

func TestUploadRoutingCoverage_Happy(t *testing.T) {
	asset := newTestAsset(t, "TEST")
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

	err := client.UploadRoutingCoverage(ctx.Context, "TEST", config.File{
		Path: asset.Name,
	})
	assert.NoError(t, err)
}
