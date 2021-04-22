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

// Package log is a substitute for the global apex/log that is not thread safe.
package log

import (
	"os"
	"sync"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/fatih/color"
)

// Interface is an extension of log.Interface.
type Interface interface {
	log.Interface
	SetColorMode(v bool)
	SetDebug(v bool)
	SetPadding(v int)
}

// Fields re-exports log.Fields from github.com/apex/log.
type Fields = log.Fields

// Log is a thread-safe wrapper for log.Logger.
type Log struct {
	log.Logger
	mu sync.RWMutex
}

// New creates a new Log instance.
func New() *Log {
	return &Log{
		Logger: log.Logger{
			Handler: cli.New(os.Stderr),
			Level:   log.InfoLevel,
		},
	}
}

// SetColorMode sets the global color mode for the logger.
func (l *Log) SetColorMode(v bool) {
	l.mu.Lock()
	defer l.mu.Unlock()

	color.NoColor = v
}

// SetDebug sets the log level to Debug or Info.
func (l *Log) SetDebug(v bool) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if v {
		l.Level = log.DebugLevel
	} else {
		l.Level = log.InfoLevel
	}
}

// SetPadding sets the padding of the log handler in a thread-safe way.
func (l *Log) SetPadding(v int) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if handler, ok := l.Handler.(*cli.Handler); ok {
		handler.Padding = v
		l.Handler = handler
	}
}
