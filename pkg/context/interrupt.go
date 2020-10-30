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

package context

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

type errReceivedSignal struct {
	Signal os.Signal
}

func (e errReceivedSignal) Error() string {
	return fmt.Sprintf("received: %s", e.Signal)
}

// Task is function that can be executed by an interrupt.
type Task func() error

// Interrupt tracks signals from the OS to determine whether to interrupt.
type Interrupt struct {
	signals chan os.Signal
	errs    chan error
}

// NewInterrupt creates an interrupt instance.
func NewInterrupt() *Interrupt {
	return &Interrupt{
		signals: make(chan os.Signal, 1),
		errs:    make(chan error, 1),
	}
}

// Run executes a given task with a given context, dealing with its timeouts,
// cancels and SIGTERM and SIGINT signals.
// It will return an error if the context is canceled, if deadline exceeds,
// if a SIGTERM or SIGINT is received and of course if the task itself fails.
func (i *Interrupt) Run(ctx context.Context, task Task) error {
	go func() {
		i.errs <- task()
	}()
	signal.Notify(i.signals, syscall.SIGINT, syscall.SIGTERM)
	select {
	case err := <-i.errs:
		return err
	case <-ctx.Done():
		return ctx.Err()
	case sig := <-i.signals:
		return errReceivedSignal{Signal: sig}
	}
}
