package cmd

import (
	"github.com/metru/gg/internal/tui"
	"github.com/spf13/cobra"
)

var tagCmd = &cobra.Command{
	Use:   "tag",
	Short: "Manage git tags",
	Long: `Open a TUI listing current tags.
Press n to create a new tag: bump version numbers with the
arrow keys and type to edit the tag name.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return tui.RunTags()
	},
}

func init() {
	rootCmd.AddCommand(tagCmd)
}