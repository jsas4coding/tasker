// Package bundler generates Taskfile.yml and Makefile from resolved task configurations.
package bundler

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
	"tasker.jsas.dev/internal/config"
	"tasker.jsas.dev/internal/constants"
	"tasker.jsas.dev/internal/resolver"
)

// TaskfileTask represents a task in Taskfile.yml format.
type TaskfileTask struct {
	Desc          string              `yaml:"desc,omitempty"`
	Summary       string              `yaml:"summary,omitempty"`
	Cmds          []string            `yaml:"cmds,omitempty"`
	Dir           string              `yaml:"dir,omitempty"`
	Deps          []string            `yaml:"deps,omitempty"`
	Silent        bool                `yaml:"silent,omitempty"`
	Preconditions []map[string]string `yaml:"preconditions,omitempty"`
}

// TaskfileOutput represents the generated Taskfile.yml structure.
type TaskfileOutput struct {
	Version string                  `yaml:"version"`
	Dotenv  []string                `yaml:"dotenv,omitempty"`
	Vars    map[string]string       `yaml:"vars,omitempty"`
	Tasks   map[string]TaskfileTask `yaml:"tasks"`
}

// GenerateTaskfile creates Taskfile.yml from a resolved project.
func GenerateTaskfile(project *resolver.ResolvedProject) ([]byte, error) {
	output := TaskfileOutput{
		Version: project.Config.Version,
		Vars:    project.Config.Vars,
		Tasks:   make(map[string]TaskfileTask),
	}

	// Collect all dotenv files from all environments (deduplicated, ordered)
	output.Dotenv = collectDotenv(project.Config.Environments)

	// Default task: show help
	output.Tasks["default"] = TaskfileTask{
		Desc:   "Show available tasks",
		Silent: true,
		Cmds:   []string{"task --list"},
	}

	// Generate tasks per group
	groupKeys := resolver.SortedGroupKeys(project.Groups)
	for _, groupKey := range groupKeys {
		tasks := project.Groups[groupKey]
		for _, rt := range tasks {
			tt := TaskfileTask{
				Desc:    rt.Description,
				Summary: fmt.Sprintf("%s - %s", rt.Name, rt.Description),
				Cmds:    rt.Cmds,
				Dir:     rt.Dir,
				Deps:    rt.Deps,
				Silent:  rt.Silent,
			}

			if rt.Environment != "" {
				tt.Preconditions = []map[string]string{
					{
						"sh":  fmt.Sprintf(`test -z "$ENV" || test "$ENV" = "%s"`, rt.Environment),
						"msg": fmt.Sprintf("Task %s requires ENV=%s (current: $ENV)", rt.FullKey, rt.Environment),
					},
				}
			}

			output.Tasks[rt.FullKey] = tt
		}
	}

	data, err := yaml.Marshal(output)
	if err != nil {
		return nil, fmt.Errorf("marshaling Taskfile.yml: %w", err)
	}

	header := constants.HeaderGenerated + "\n" + constants.HeaderSource + "\n" + constants.HeaderRegenerate + "\n\n"
	return append([]byte(header), data...), nil
}

// WriteTaskfile writes the generated Taskfile.yml to disk.
func WriteTaskfile(project *resolver.ResolvedProject, dir string) error {
	data, err := GenerateTaskfile(project)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, constants.TaskfileOutput), data, constants.FilePermissions)
}

func collectDotenv(envs map[string]config.Environment) []string {
	seen := make(map[string]bool)
	var result []string

	keys := make([]string, 0, len(envs))
	for k := range envs {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		env := envs[k]
		for _, d := range env.Dotenv {
			if !seen[d] {
				seen[d] = true
				result = append(result, d)
			}
		}
	}
	return result
}

// TaskColonToDash converts task key colons to dashes for Makefile targets.
func TaskColonToDash(key string) string {
	return strings.ReplaceAll(key, ":", "-")
}
