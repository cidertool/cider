package docs

import (
	"path/filepath"

	"github.com/apex/log"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

func runDocsManCmd(cmd *cobra.Command, args []string) error {
	var path string
	if len(args) == 0 {
		path = defaultDocsPath
	} else {
		path = args[0]
	}

	path = filepath.Join(path, "man")

	log.WithField("path", path).Info("generating man documentation")

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
