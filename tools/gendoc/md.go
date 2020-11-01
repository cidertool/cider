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

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/apex/log"
	commands "github.com/cidertool/cider/pkg/cmd"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

const docsMdFrontmatterTemplate = `---
layout: page
parent: Commands
title: %s
nav_order: %d
nav_exclude: %t
---

`

type pageNavField struct {
	order   int
	exclude bool
}

func runDocsMdCmd(cmd *cobra.Command, args []string) error {
	var orderRoot, orderInit, orderRelease, orderCheck, orderCompletions = 0, 1, 2, 3, 4

	var pageNavFields = map[string]pageNavField{
		"cider.md":             {order: orderRoot},
		"cider_init.md":        {order: orderInit},
		"cider_release.md":     {order: orderRelease},
		"cider_check.md":       {order: orderCheck},
		"cider_completions.md": {order: orderCompletions},
	}

	var dir string
	if len(args) == 0 {
		dir = defaultDocsPath
	} else {
		dir = args[0]
	}

	dir = filepath.Join(dir, "commands")

	prepender := func(filename string) string {
		base := filepath.Base(filename)
		return fmt.Sprintf(docsMdFrontmatterTemplate, pageTitle(base), pageNavFields[base].order, pageNavFields[base].exclude)
	}

	linkHandler := func(name string) string {
		base := strings.TrimSuffix(name, filepath.Ext(name))
		return "/commands/" + strings.ToLower(base) + "/"
	}

	log.WithField("path", dir).Info("generating Markdown documentation")

	err := doc.GenMarkdownTreeCustom(commands.NewRoot("dev", os.Exit).Cmd, dir, prepender, linkHandler)
	if err != nil {
		log.Error("generation failed")
	} else {
		log.Info("generation completed successfully")
	}

	return err
}

func pageTitle(s string) string {
	s = strings.TrimSuffix(s, filepath.Ext(s))
	if s != "cider" {
		s = strings.ReplaceAll(s, "cider", "")
	}

	s = strings.ReplaceAll(s, "_", " ")
	s = strings.TrimSpace(s)

	return s
}
