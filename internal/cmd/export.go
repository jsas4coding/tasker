package cmd

import (
	"github.com/spf13/cobra"
	"tasker.jsas.dev/internal/bundler"
	"tasker.jsas.dev/internal/config"
	"tasker.jsas.dev/internal/output"
	"tasker.jsas.dev/internal/resolver"
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export resolved config as Tasker.json",
	Long:  "Loads, validates, and resolves the Tasker configuration, then writes Tasker.json.",
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

		resolver.InjectBuiltins(resolved)

		if err := bundler.WriteTaskerJSON(resolved, "."); err != nil {
			return err
		}
		output.Success("Generated Tasker.json")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)
}
