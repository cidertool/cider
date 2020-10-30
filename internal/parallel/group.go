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

// Package parallel wraps a error group with a semaphore with configurable
// size, so you can control the number of tasks being executed simultaneously.
package parallel

import (
	"sync"

	"golang.org/x/sync/errgroup"
)

// Group is the Semphore ErrorGroup itself.
type Group interface {
	Go(func() error)
	Wait() error
}

// New returns a new Group of a given size.
func New(size int) Group {
	if size <= 1 {
		return &serialGroup{}
	}

	return &parallelGroup{
		ch: make(chan bool, size),
		g:  errgroup.Group{},
	}
}

type parallelGroup struct {
	ch chan bool
	g  errgroup.Group
}

// Go execs one function respecting the group and semaphore.
func (s *parallelGroup) Go(fn func() error) {
	s.g.Go(func() error {
		s.ch <- true
		defer func() {
			<-s.ch
		}()
		return fn()
	})
}

// Wait waits for the group to complete and return an error if any.
func (s *parallelGroup) Wait() error {
	return s.g.Wait()
}

type serialGroup struct {
	err     error
	errOnce sync.Once
}

// Go execs runs `fn` and saves the result if no error has been encountered.
func (s *serialGroup) Go(fn func() error) {
	if s.err != nil {
		return
	}

	if err := fn(); err != nil {
		s.errOnce.Do(func() {
			s.err = err
		})
	}
}

// Wait waits for Go to complete and returns the first error encountered.
func (s *serialGroup) Wait() error {
	return s.err
}
