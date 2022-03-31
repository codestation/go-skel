package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"megpoid.xyz/go/go-skel/internal/version"
)

const versionFormatter = "GoApp version: %s, commit: %s, date: %s, modified: %t\n"

func printVersion() {
	fmt.Printf(versionFormatter, version.Tag, version.Revision, version.LastCommit, version.Modified)
}

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display version",
	Long:  `Prints the version and build info`,
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		printVersion()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
