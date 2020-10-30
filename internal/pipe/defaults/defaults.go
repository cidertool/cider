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

// Package defaults runs all defaulter pipelines
package defaults

import (
	"github.com/cidertool/cider/internal/defaults"
	"github.com/cidertool/cider/internal/middleware"
	"github.com/cidertool/cider/pkg/context"
)

// Pipe that sets the defaults.
type Pipe struct {
	defaulters []defaults.Defaulter
}

func (Pipe) String() string {
	return "setting defaults"
}

// Run the pipe.
func (p Pipe) Run(ctx *context.Context) error {
	if len(p.defaulters) == 0 {
		p.defaulters = defaults.Defaulters
	}

	for _, defaulter := range p.defaulters {
		if err := middleware.Logging(
			defaulter.String(),
			middleware.ErrHandler(defaulter.Default),
			middleware.ExtraPadding,
		)(ctx); err != nil {
			return err
		}
	}

	return nil
}
