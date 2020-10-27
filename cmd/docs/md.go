package docs

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/apex/log"
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
		"cider_docs.md":        {exclude: true},
		"cider_docs_config.md": {exclude: true},
		"cider_docs_man.md":    {exclude: true},
		"cider_docs_md.md":     {exclude: true},
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

	err := doc.GenMarkdownTreeCustom(cmd.Root(), dir, prepender, linkHandler)
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
