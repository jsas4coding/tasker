package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"tasker.jsas.dev/internal/bundler"
	"tasker.jsas.dev/internal/config"
	"tasker.jsas.dev/internal/output"
	"tasker.jsas.dev/internal/resolver"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Show structured task list with groups and environments",
	RunE: func(_ *cobra.Command, _ []string) error {
		project, err := config.Load(".")
		if err != nil {
			return err
		}

		resolved, err := resolver.Resolve(project)
		if err != nil {
			return err
		}

		resolver.InjectBuiltins(resolved)

		output.Section(fmt.Sprintf("%s - %s", project.Config.Name, project.Config.Description))

		if len(project.Config.Environments) > 0 {
			fmt.Println("Environments:")
			for key, env := range project.Config.Environments {
				dotenv := strings.Join(env.Dotenv, ", ")
				fmt.Printf("  %-10s %-20s %s\n", key, env.Name, dotenv)
			}
			fmt.Println()
		}

		groupKeys := resolver.SortedGroupKeys(resolved.Groups)
		for _, groupKey := range groupKeys {
			groupName, groupDesc := bundler.GroupMetadata(project.Config, groupKey)
			tasks := resolved.Groups[groupKey]

			fmt.Printf("%s  %s\n", groupName, groupDesc)
			for _, rt := range tasks {
				fmt.Printf("  %-30s %-25s %s\n", rt.FullKey, rt.Name, rt.Description)
			}
			fmt.Println()
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
