// Package resolver validates and resolves task configurations into a flat structure.
package resolver

import (
	"fmt"
	"sort"

	"tasker.jsas.dev/internal/config"
	"tasker.jsas.dev/internal/constants"
)

// ResolvedTask is a fully qualified task ready for output generation.
type ResolvedTask struct {
	FullKey     string // e.g., "go:dev:build"
	GroupKey    string // e.g., "go"
	TaskKey     string // e.g., "dev:build"
	Name        string
	Description string
	Environment string
	Cmds        []string
	Dir         string
	Deps        []string
	Silent      bool
}

// ResolvedProject holds all resolved tasks grouped by group key.
type ResolvedProject struct {
	Config *config.Config
	Groups map[string][]ResolvedTask // key = group key
	RootDir string
}

// Resolve processes a loaded project into resolved tasks with validation.
func Resolve(project *config.Project) (*ResolvedProject, error) {
	resolved := &ResolvedProject{
		Config:  project.Config,
		Groups:  make(map[string][]ResolvedTask),
		RootDir: project.RootDir,
	}

	for groupKey, gf := range project.GroupFiles {
		if _, ok := project.Config.Groups[groupKey]; !ok {
			return nil, fmt.Errorf("group file for %q exists but group not declared in config", groupKey)
		}

		var tasks []ResolvedTask
		for taskKey, task := range gf.Tasks {
			fullKey := groupKey + ":" + taskKey

			rt := ResolvedTask{
				FullKey:     fullKey,
				GroupKey:    groupKey,
				TaskKey:     taskKey,
				Name:        task.Name,
				Description: task.Description,
				Environment: task.Environment,
				Cmds:        task.Cmds,
				Dir:         task.Dir,
				Deps:        task.Deps,
				Silent:      task.Silent,
			}
			tasks = append(tasks, rt)
		}

		sort.Slice(tasks, func(i, j int) bool {
			return tasks[i].FullKey < tasks[j].FullKey
		})

		resolved.Groups[groupKey] = tasks
	}

	return resolved, nil
}

// FilterByEnv returns tasks filtered to a specific environment.
// Tasks with no environment set are always included.
func FilterByEnv(tasks []ResolvedTask, env string) []ResolvedTask {
	if env == "" {
		return tasks
	}
	var filtered []ResolvedTask
	for _, t := range tasks {
		if t.Environment == "" || t.Environment == env {
			filtered = append(filtered, t)
		}
	}
	return filtered
}

// InjectBuiltins appends the built-in tasker:* tasks to the resolved project.
func InjectBuiltins(project *ResolvedProject) {
	builtins := []ResolvedTask{
		{FullKey: "tasker:export", GroupKey: constants.BuiltinGroupKey, TaskKey: "export", Name: "Export", Description: "Export resolved config as Tasker.json", Cmds: []string{"tasker export"}},
		{FullKey: "tasker:generate", GroupKey: constants.BuiltinGroupKey, TaskKey: "generate", Name: "Generate", Description: "Generate Taskfile.yml, Makefile, and Tasker.json", Cmds: []string{"tasker generate"}},
		{FullKey: "tasker:init", GroupKey: constants.BuiltinGroupKey, TaskKey: "init", Name: "Init", Description: "Scaffold a new Tasker project", Cmds: []string{"tasker init"}},
		{FullKey: "tasker:list", GroupKey: constants.BuiltinGroupKey, TaskKey: "list", Name: "List", Description: "Show structured task list", Cmds: []string{"tasker list"}},
		{FullKey: "tasker:validate", GroupKey: constants.BuiltinGroupKey, TaskKey: "validate", Name: "Validate", Description: "Validate Tasker configuration", Cmds: []string{"tasker validate"}},
		{FullKey: "tasker:version", GroupKey: constants.BuiltinGroupKey, TaskKey: "version", Name: "Version", Description: "Print version information", Cmds: []string{"tasker version"}},
	}
	project.Groups[constants.BuiltinGroupKey] = builtins
}

// SortedGroupKeys returns group keys in alphabetical order.
func SortedGroupKeys(groups map[string][]ResolvedTask) []string {
	keys := make([]string, 0, len(groups))
	for k := range groups {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
