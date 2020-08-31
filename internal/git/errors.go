package git

import (
	"errors"
	"fmt"
	"path/filepath"
)

// ErrDirty happens when the repo has uncommitted/unstashed changes.
type ErrDirty struct {
	Status string
}

func (e ErrDirty) Error() string {
	return fmt.Sprintf("git is currently in a dirty state, please check in your pipeline what can be changing the following files:\n%v", e.Status)
}

// ErrWrongRef happens when the HEAD reference is different from the tag being built.
type ErrWrongRef struct {
	Commit, Tag string
}

func (e ErrWrongRef) Error() string {
	return fmt.Sprintf("git tag %v was not made against commit %v", e.Tag, e.Commit)
}

// ErrNoTag happens if the underlying git repository doesn't contain any tags
var ErrNoTag = errors.New("git doesn't contain any tags")

// ErrNotRepository happens if you try to run applereleaser against a folder
// which is not a git repository.
type ErrNotRepository struct {
	Dir string
}

func (e ErrNotRepository) Error() string {
	return fmt.Sprintf("the directory at %s is not a git repository", filepath.Clean(e.Dir))
}

// ErrNoGit happens when git is not present in PATH.
var ErrNoGit = errors.New("git not present in PATH")
