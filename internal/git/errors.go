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

package git

import (
	"errors"
	"fmt"
	"path/filepath"
)

const errDirtyHeading = "git is currently in a dirty state, please check in your pipeline what can be changing the following files"

// ErrDirty happens when the repo has uncommitted/unstashed changes.
type ErrDirty struct {
	Status string
}

func (e ErrDirty) Error() string {
	return fmt.Sprintf("%s:\n%v", errDirtyHeading, e.Status)
}

// ErrWrongRef happens when the HEAD reference is different from the tag being built.
type ErrWrongRef struct {
	Commit, Tag string
}

func (e ErrWrongRef) Error() string {
	return fmt.Sprintf("git tag %v was not made against commit %v", e.Tag, e.Commit)
}

// ErrNoTag happens if the underlying git repository doesn't contain any tags.
var ErrNoTag = errors.New("git doesn't contain any tags")

// ErrNotRepository happens if you try to run Cider against a folder
// which is not a git repository.
type ErrNotRepository struct {
	Dir string
}

func (e ErrNotRepository) Error() string {
	return fmt.Sprintf("the directory at %s is not a git repository", filepath.Clean(e.Dir))
}

// ErrNoGit happens when git is not present in PATH.
var ErrNoGit = errors.New("git not present in PATH")

// ErrNoRemoteOrigin happens when the repository has no remote named "origin".
var ErrNoRemoteOrigin = errors.New("repository doesn't have an `origin` remote")
