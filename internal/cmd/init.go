package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"tasker.jsas.dev/internal/config"
	"tasker.jsas.dev/internal/constants"
	"tasker.jsas.dev/internal/output"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Scaffold .tasker/ directory with config and task files",
	Long:  "Initializes a new Tasker configuration by detecting package managers and creating starter files.",
	RunE: func(_ *cobra.Command, _ []string) error {
		dir := "."

		cfgPath := filepath.Join(dir, constants.ConfigDir, constants.ConfigFile)
		if _, err := os.Stat(cfgPath); err == nil {
			return fmt.Errorf("%s/%s already exists, use 'tasker generate' to regenerate outputs", constants.ConfigDir, constants.ConfigFile)
		}

		managers := config.DetectManagers(dir)
		projectName := filepath.Base(absPath(dir))

		cfg := config.Config{
			Name:        projectName,
			Description: fmt.Sprintf("%s project tasks", projectName),
			Version:     "3",
			Environments: map[string]config.Environment{
				"dev": {
					Name:        "Development",
					Description: "Local development environment",
					Dotenv:      []string{".env", ".env.local"},
				},
				"test": {
					Name:        "Testing",
					Description: "Automated testing environment",
					Dotenv:      []string{".env", ".env.test"},
				},
				"prod": {
					Name:        "Production",
					Description: "Production environment",
					Dotenv:      []string{".env", ".env.production"},
				},
			},
			Vars: map[string]string{
				"PROJECT_NAME": projectName,
				"PROJECT_ROOT": "{{.ROOT_DIR}}",
			},
			Groups: make(map[string]config.Group),
		}

		for _, m := range managers {
			cfg.Groups[m.Key] = config.Group{
				Name:        m.Name,
				Description: m.Description,
			}
		}

		if len(cfg.Groups) == 0 {
			cfg.Groups["general"] = config.Group{
				Name:        "General",
				Description: "General project tasks",
			}
		}

		// Create .tasker/ directory
		configDir := filepath.Join(dir, constants.ConfigDir)
		if err := os.MkdirAll(configDir, constants.DirPermissions); err != nil {
			return fmt.Errorf("creating %s: %w", constants.ConfigDir, err)
		}

		// Write .tasker/config.yml with schema reference
		data, err := yaml.Marshal(cfg)
		if err != nil {
			return fmt.Errorf("marshaling config: %w", err)
		}

		header := "# yaml-language-server: $schema=schemas/tasker.schema.json\n"
		if err := os.WriteFile(cfgPath, append([]byte(header), data...), constants.FilePermissions); err != nil {
			return fmt.Errorf("writing %s: %w", constants.ConfigFile, err)
		}
		output.Successf("Created %s/%s", constants.ConfigDir, constants.ConfigFile)

		// Create .tasker/tasks/ directory and starter files
		tasksDir := filepath.Join(dir, constants.TasksDir)
		if err := os.MkdirAll(tasksDir, constants.DirPermissions); err != nil {
			return fmt.Errorf("creating %s: %w", constants.TasksDir, err)
		}

		for key, group := range cfg.Groups {
			filename := config.GroupFileName(key)
			groupPath := filepath.Join(tasksDir, filename)

			starter := generateStarterTasks(key, group)
			gfData, err := yaml.Marshal(starter)
			if err != nil {
				return fmt.Errorf("marshaling %s: %w", filename, err)
			}

			gfHeader := "# yaml-language-server: $schema=../schemas/tasks.schema.json\n"
			if err := os.WriteFile(groupPath, append([]byte(gfHeader), gfData...), constants.FilePermissions); err != nil {
				return fmt.Errorf("writing %s: %w", filename, err)
			}
			output.Successf("Created %s/%s", constants.TasksDir, filename)
		}

		// Export schemas for IDE support
		if err := exportSchemas(dir); err != nil {
			return fmt.Errorf("exporting schemas: %w", err)
		}
		output.Successf("Created %s/schemas/", constants.ConfigDir)

		fmt.Println()
		output.Info("Detected package managers:")
		if len(managers) > 0 {
			for _, m := range managers {
				fmt.Printf("  - %s (%s)\n", m.Name, m.File)
			}
		} else {
			output.Dim("  (none detected)")
		}

		fmt.Println()
		output.Info("Next steps:")
		fmt.Printf("  1. Edit %s/%s and %s/*.yml to define your tasks\n", constants.ConfigDir, constants.ConfigFile, constants.TasksDir)
		fmt.Println("  2. Run 'tasker generate' to create Taskfile.yml and Makefile")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func absPath(dir string) string {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return dir
	}
	return abs
}

func exportSchemas(dir string) error {
	schemasDir := filepath.Join(dir, constants.SchemasDir)
	if err := os.MkdirAll(schemasDir, constants.DirPermissions); err != nil {
		return err
	}

	fs := config.SchemaFS()
	entries, err := fs.ReadDir("schemas")
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		data, err := fs.ReadFile("schemas/" + entry.Name())
		if err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(schemasDir, entry.Name()), data, constants.FilePermissions); err != nil {
			return err
		}
	}
	return nil
}

func generateStarterTasks(key string, group config.Group) config.GroupFile {
	gf := config.GroupFile{
		Tasks: make(map[string]config.Task),
	}

	lowerName := strings.ToLower(group.Name)

	gf.Tasks["dev:build"] = config.Task{
		Name:        fmt.Sprintf("Build (%s dev)", lowerName),
		Description: fmt.Sprintf("Build %s for development", lowerName),
		Environment: "dev",
		Cmds:        []string{fmt.Sprintf("echo 'TODO: Add %s dev build command'", key)},
	}

	gf.Tasks["dev:run"] = config.Task{
		Name:        fmt.Sprintf("Run (%s dev)", lowerName),
		Description: fmt.Sprintf("Run %s in development mode", lowerName),
		Environment: "dev",
		Cmds:        []string{fmt.Sprintf("echo 'TODO: Add %s dev run command'", key)},
	}

	return gf
}
