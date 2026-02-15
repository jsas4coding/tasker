package cmd

import (
	"github.com/spf13/cobra"
	"tasker.jsas.dev/internal/bundler"
	"tasker.jsas.dev/internal/config"
	"tasker.jsas.dev/internal/output"
	"tasker.jsas.dev/internal/resolver"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Bundle .tasker/ config into Taskfile.yml and Makefile",
	Long:  "Reads the Tasker configuration and generates Taskfile.yml and Makefile.",
	RunE: func(_ *cobra.Command, _ []string) error {
		project, err := config.Load(".")
		if err != nil {
			return err
		}

		if errs := project.Validate(); len(errs) > 0 {
			for _, e := range errs {
				output.Errorf("%s", e)
			}
			return errValidation(len(errs))
		}

		resolved, err := resolver.Resolve(project)
		if err != nil {
			return err
		}

		if err := bundler.WriteTaskfile(resolved, "."); err != nil {
			return err
		}
		output.Success("Generated Taskfile.yml")

		if err := bundler.WriteMakefile(resolved, "."); err != nil {
			return err
		}
		output.Success("Generated Makefile")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
}
