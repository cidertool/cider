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

package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name                           string
		version, commit, date, builtBy string
		out                            string
	}{
		{
			name: "all empty",
			out:  "",
		},
		{
			name:    "complete",
			version: "1.2.3",
			date:    "12/12/12",
			commit:  "aaaa",
			builtBy: "me",
			out:     "1.2.3\ncommit: aaaa\nbuilt at: 12/12/12\nbuilt by: me",
		},
		{
			name:    "only version",
			version: "1.2.3",
			out:     "1.2.3",
		},
		{
			name:    "version and date",
			version: "1.2.3",
			date:    "12/12/12",
			out:     "1.2.3\nbuilt at: 12/12/12",
		},
		{
			name:    "version, date, built by",
			version: "1.2.3",
			date:    "12/12/12",
			builtBy: "me",
			out:     "1.2.3\nbuilt at: 12/12/12\nbuilt by: me",
		},
	}

	for _, tt := range testCases {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.out+licenseDisclaimer, buildVersion(tt.version, tt.commit, tt.date, tt.builtBy))
		})
	}
}
