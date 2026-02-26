package bundler

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"time"

	"tasker.jsas.dev/internal/config"
	"tasker.jsas.dev/internal/constants"
	"tasker.jsas.dev/internal/resolver"
)

// TaskerJSONOutput represents the full exported Tasker.json structure.
type TaskerJSONOutput struct {
	Meta         TaskerJSONMeta               `json:"meta"`
	Module       string                       `json:"module,omitempty"`
	Name         string                       `json:"name"`
	Description  string                       `json:"description"`
	Version      string                       `json:"version"`
	Vars         map[string]string            `json:"vars,omitempty"`
	Environments map[string]TaskerJSONEnv     `json:"environments,omitempty"`
	Groups       []TaskerJSONGroup            `json:"groups"`
}

// TaskerJSONMeta holds generation metadata.
type TaskerJSONMeta struct {
	GeneratedAt  string `json:"generatedAt"`
	GeneratedBy  string `json:"generatedBy"`
	SourceConfig string `json:"sourceConfig"`
}

// TaskerJSONEnv represents an environment in the JSON output.
type TaskerJSONEnv struct {
	Key         string   `json:"key"`
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	Dotenv      []string `json:"dotenv,omitempty"`
}

// TaskerJSONGroup represents a task group in the JSON output.
type TaskerJSONGroup struct {
	Key         string           `json:"key"`
	Name        string           `json:"name"`
	Description string           `json:"description,omitempty"`
	Builtin     bool             `json:"builtin"`
	Tasks       []TaskerJSONTask `json:"tasks"`
}

// TaskerJSONTask represents a single task in the JSON output.
type TaskerJSONTask struct {
	FullKey     string   `json:"fullKey"`
	GroupKey    string   `json:"groupKey"`
	TaskKey     string   `json:"taskKey"`
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	Environment string   `json:"environment,omitempty"`
	Cmds        []string `json:"cmds"`
	Dir         string   `json:"dir,omitempty"`
	Deps        []string `json:"deps,omitempty"`
	Silent      bool     `json:"silent,omitempty"`
	Builtin     bool     `json:"builtin"`
}

// GenerateTaskerJSON creates the Tasker.json content from a resolved project.
func GenerateTaskerJSON(project *resolver.ResolvedProject) ([]byte, error) {
	out := TaskerJSONOutput{
		Meta: TaskerJSONMeta{
			GeneratedAt:  time.Now().UTC().Format(time.RFC3339),
			GeneratedBy:  "tasker",
			SourceConfig: filepath.Join(constants.ConfigDir, constants.ConfigFile),
		},
		Module:      project.Config.Module,
		Name:        project.Config.Name,
		Description: project.Config.Description,
		Version:     project.Config.Version,
		Vars:        project.Config.Vars,
	}

	// Environments as map
	if len(project.Config.Environments) > 0 {
		out.Environments = make(map[string]TaskerJSONEnv, len(project.Config.Environments))
		for key, env := range project.Config.Environments {
			out.Environments[key] = TaskerJSONEnv{
				Key:         key,
				Name:        env.Name,
				Description: env.Description,
				Dotenv:      env.Dotenv,
			}
		}
	}

	// Groups as sorted array
	groupKeys := resolver.SortedGroupKeys(project.Groups)
	for _, groupKey := range groupKeys {
		tasks := project.Groups[groupKey]
		builtin := groupKey == constants.BuiltinGroupKey

		groupName, groupDesc := GroupMetadata(project.Config, groupKey)

		jg := TaskerJSONGroup{
			Key:         groupKey,
			Name:        groupName,
			Description: groupDesc,
			Builtin:     builtin,
			Tasks:       make([]TaskerJSONTask, 0, len(tasks)),
		}

		sorted := make([]resolver.ResolvedTask, len(tasks))
		copy(sorted, tasks)
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].FullKey < sorted[j].FullKey
		})

		for _, rt := range sorted {
			jt := TaskerJSONTask{
				FullKey:     rt.FullKey,
				GroupKey:    rt.GroupKey,
				TaskKey:     rt.TaskKey,
				Name:        rt.Name,
				Description: rt.Description,
				Environment: rt.Environment,
				Cmds:        rt.Cmds,
				Dir:         rt.Dir,
				Deps:        rt.Deps,
				Silent:      rt.Silent,
				Builtin:     builtin,
			}
			jg.Tasks = append(jg.Tasks, jt)
		}

		out.Groups = append(out.Groups, jg)
	}

	return json.MarshalIndent(out, "", "  ")
}

// WriteTaskerJSON writes Tasker.json to disk.
func WriteTaskerJSON(project *resolver.ResolvedProject, dir string) error {
	data, err := GenerateTaskerJSON(project)
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(filepath.Join(dir, constants.TaskerJSONOutput), data, constants.FilePermissions)
}

// GroupMetadata returns name and description for a group key, handling the
// built-in tasker group that has no config entry.
func GroupMetadata(cfg *config.Config, groupKey string) (string, string) {
	if groupKey == constants.BuiltinGroupKey {
		return "Tasker", "Built-in management commands"
	}
	g := cfg.Groups[groupKey]
	return g.Name, g.Description
}
