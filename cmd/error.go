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

package cmd

type exitError struct {
	err     error
	code    int
	details string
}

func wrapErrorWithCode(err error, code int, details string) *exitError {
	return &exitError{
		err:     err,
		code:    code,
		details: details,
	}
}

func wrapError(err error, log string) *exitError {
	return wrapErrorWithCode(err, 1, log)
}

func (e *exitError) Error() string {
	return e.err.Error()
}
