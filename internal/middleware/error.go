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
	"github.com/cidertool/cider/internal/pipe"
	"github.com/cidertool/cider/pkg/context"
)

// ErrHandler handles an action error, ignoring and logging pipe skipped
// errors.
func ErrHandler(action Action) Action {
	return func(ctx *context.Context) error {
		var err = action(ctx)
		if err == nil {
			return nil
		}

		if pipe.IsSkip(err) {
			log.WithError(err).Warn("pipe skipped")
			return nil
		}

		return err
	}
}
