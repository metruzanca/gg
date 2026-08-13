package cmd

import (
	"github.com/metru/gg/internal/tui"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gg",
	Short: "gg - a minimal git TUI",
	Long: `gg is a terminal UI for git built with Bubble Tea.
It shows staged and unstaged changes, lets you stage/unstage files,
and create conventional commits.

Run without arguments to open the TUI.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return tui.Run()
	},
}

func Execute() error {
	return rootCmd.Execute()
}
