package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"tasker.jsas.dev/internal/config"
	"tasker.jsas.dev/internal/output"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate Tasker configuration without generating",
	RunE: func(_ *cobra.Command, _ []string) error {
		project, err := config.Load(".")
		if err != nil {
			return err
		}

		errs := project.Validate()
		if len(errs) > 0 {
			for _, e := range errs {
				output.Errorf("%s", e)
			}
			return errValidation(len(errs))
		}

		output.Success("Configuration is valid")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)
}

func errValidation(count int) error {
	return fmt.Errorf("validation failed with %d error(s)", count)
}
