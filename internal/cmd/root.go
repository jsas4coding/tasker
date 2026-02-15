// Package cmd provides CLI command implementations for the Tasker application.
package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tasker",
	Short: "Task bundler for Taskfile.yml and Makefile generation",
	Long: `Tasker reads structured configuration (.tasker/config.yml + .tasker/tasks/*.yml)
and bundles into a single Taskfile.yml and Makefile.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return generateCmd.RunE(cmd, args)
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}
