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

package middleware

import (
	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/cidertool/cider/pkg/context"
	"github.com/fatih/color"
)

// Padding is a logging initial padding.
type Padding int

// DefaultInitialPadding is the default padding in the log library.
const DefaultInitialPadding Padding = 3

// ExtraPadding is the double of the DefaultInitialPadding.
const ExtraPadding Padding = DefaultInitialPadding * 2

// Logging pretty prints the given action and its title.
// You can have different padding levels by providing different initial
// paddings. The middleware will print the title in the given padding and the
// action logs in padding+default padding.
// The default padding in the log library is 3.
// The middleware always resets to the default padding.
func Logging(title string, next Action, padding Padding) Action {
	return func(ctx *context.Context) error {
		defer func() {
			cli.Default.Padding = int(DefaultInitialPadding)
		}()

		cli.Default.Padding = int(padding)

		log.Info(color.New(color.Bold).Sprint(title))

		cli.Default.Padding = int(padding + DefaultInitialPadding)

		return next(ctx)
	}
}
