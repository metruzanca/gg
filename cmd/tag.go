package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	ggit "github.com/metruzanca/gg/internal/git"
	"github.com/metruzanca/gg/internal/tui"
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

var tagBumpCmd = &cobra.Command{
	Use:   "bump",
	Short: "Bump the latest tag's minor version",
	Long:  "Increment the minor component of the latest versioned tag, create it, then ask whether to push.",
	RunE: func(cmd *cobra.Command, args []string) error {
		tags, err := ggit.Tags()
		if err != nil {
			return fmt.Errorf("list tags: %w", err)
		}
		latest := ggit.LatestVersion(tags)
		if latest == "" {
			return fmt.Errorf("no versioned tags found to bump")
		}
		name, ok := ggit.BumpMinor(latest)
		if !ok {
			return fmt.Errorf("could not bump minor version of %q", latest)
		}
		if err := ggit.CreateTag(name); err != nil {
			return fmt.Errorf("create tag: %w", err)
		}
		fmt.Printf("Created tag %s (bumped minor from %s)\n", name, latest)

		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("Push %s and its commits to origin? [y/N]: ", name)
		answer, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("read answer: %w", err)
		}
		if strings.TrimSpace(strings.ToLower(answer)) == "y" {
			if err := ggit.PushTagWithCommits(name); err != nil {
				return fmt.Errorf("push: %w", err)
			}
			fmt.Printf("Pushed tag %s and its commits\n", name)
		} else {
			fmt.Printf("Tag %s created locally (not pushed)\n", name)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(tagCmd)
	tagCmd.AddCommand(tagBumpCmd)
}
