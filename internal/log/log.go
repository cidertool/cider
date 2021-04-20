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

// Package logger is a substitute for the global apex/log that is not thread safe
package log

import (
	"sync"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/fatih/color"
)

type Interface interface {
	log.Interface
	SetColorMode(v bool)
	SetDebug(v bool)
	SetPadding(v int)
}

type Fields = log.Fields

type Log struct {
	log.Logger
	mu sync.RWMutex
}

func New() *Log {
	return &Log{
		Logger: log.Logger{
			Handler: cli.Default,
			Level:   log.InfoLevel,
		},
	}
}

func (l *Log) SetColorMode(v bool) {
	l.mu.Lock()
	defer l.mu.Unlock()

	color.NoColor = false
}

func (l *Log) SetDebug(v bool) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.Logger.Level = log.DebugLevel
}

func (l *Log) SetPadding(v int) {
	l.mu.Lock()
	defer l.mu.Unlock()

	cli.Default.Padding = v
}
