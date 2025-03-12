package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version information
var (
	Version = "0.1.0"
	Commit  = "none"
	Date    = "unknown"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Co",
	Long:  `All software has versions. This is Co's.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Co version %s (commit: %s, built at: %s)\n", Version, Commit, Date)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
