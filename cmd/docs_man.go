package cmd

import (
	"path/filepath"

	"github.com/apex/log"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

func newDocsManCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "man",
		Short: "Generate man documentation for Cider.",
		Args:  cobra.MaximumNArgs(1),
		RunE:  runDocsManCmd,
	}
}

func runDocsManCmd(cmd *cobra.Command, args []string) error {
	var path string
	if len(args) == 0 {
		path = "docs"
	} else {
		path = args[0]
	}
	path = filepath.Join(path, "man")

	log.WithField("path", path).Info("generating Cider documentation")
	err := doc.GenManTreeFromOpts(cmd.Root(), doc.GenManTreeOptions{
		Path: path,
	})
	if err != nil {
		log.Error("generation failed")
	} else {
		log.Info("generation completed successfully")
	}
	return err
}
