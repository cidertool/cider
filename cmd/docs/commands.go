package docs

import (
	"github.com/spf13/cobra"
)

const defaultDocsPath = "docs"

// CmdConfig returns the cobra.Command for the man subcommand.
func CmdConfig() *cobra.Command {
	return &cobra.Command{
		Use:   "config",
		Short: "Generate configuration file documentation for Cider.",
		Args:  cobra.MaximumNArgs(1),
		RunE:  runDocsConfigCmd,
	}
}

// CmdMan returns the cobra.Command for the man subcommand.
func CmdMan() *cobra.Command {
	return &cobra.Command{
		Use:   "man",
		Short: "Generate man documentation for Cider.",
		Args:  cobra.MaximumNArgs(1),
		RunE:  runDocsManCmd,
	}
}

// CmdMarkdown returns the cobra.Command for the man subcommand.
func CmdMarkdown() *cobra.Command {
	return &cobra.Command{
		Use:   "md",
		Short: "Generate Markdown documentation for Cider.",
		Args:  cobra.MaximumNArgs(1),
		RunE:  runDocsMdCmd,
	}
}
