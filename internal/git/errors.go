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
