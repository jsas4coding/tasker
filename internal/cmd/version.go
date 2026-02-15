package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Build information. These variables are set via ldflags at build time.
var (
	// Version is the application version (set at build time).
	Version = "dev"
	// BuildTime is the build timestamp (set at build time).
	BuildTime = "unknown"
	// GitCommit is the git commit hash (set at build time).
	GitCommit = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Printf("Tasker version information:\n")
		fmt.Printf("  Version:    %s\n", Version)
		fmt.Printf("  Build Time: %s\n", BuildTime)
		fmt.Printf("  Git Commit: %s\n", GitCommit)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
