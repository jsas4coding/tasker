// Package config provides configuration parsing and validation for Tasker.
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"gopkg.in/yaml.v3"
	"tasker.jsas.dev/internal/constants"
)

// identifierPattern validates keys used as identifiers (group keys, environment
// keys). These become Makefile targets and are interpolated into shell commands,
// so they must contain only safe characters. This mirrors the JSON Schema
// propertyNames pattern as defense in depth.
var identifierPattern = regexp.MustCompile(`^[a-z][a-z0-9_-]*$`)

// EcosystemEntry represents a child project in ecosystem mode.
type EcosystemEntry struct {
	Path string `yaml:"path"`
}

// Config represents the top-level config.yml configuration.
type Config struct {
	Module       string                    `yaml:"module,omitempty"`
	Name         string                    `yaml:"name"`
	Description  string                    `yaml:"description"`
	Version      string                    `yaml:"version"`
	Environments map[string]Environment    `yaml:"environments"`
	Vars         map[string]string         `yaml:"vars,omitempty"`
	Ecosystem    map[string]EcosystemEntry `yaml:"ecosystem,omitempty"`
	Groups       map[string]Group          `yaml:"groups"`
}

// Project holds the fully loaded configuration including group task files.
type Project struct {
	Config     *Config
	GroupFiles map[string]*GroupFile
	RootDir    string
}

// Load reads .tasker/config.yml from the given directory and loads all group files.
func Load(dir string) (*Project, error) {
	cfgPath := filepath.Join(dir, constants.ConfigDir, constants.ConfigFile)

	data, err := os.ReadFile(cfgPath)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", cfgPath, err)
	}

	if err := ValidateSchema(data, constants.TaskSchemaFile); err != nil {
		return nil, fmt.Errorf("config.yml %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing %s: %w", cfgPath, err)
	}

	if cfg.Version == "" {
		cfg.Version = "3"
	}

	project := &Project{
		Config:     &cfg,
		GroupFiles: make(map[string]*GroupFile),
		RootDir:    dir,
	}

	for key := range cfg.Groups {
		filename := GroupFileName(key)
		groupPath := filepath.Join(dir, constants.TasksDir, filename)

		gf, err := LoadGroupFile(groupPath)
		if err != nil {
			if os.IsNotExist(err) {
				project.GroupFiles[key] = &GroupFile{Tasks: make(map[string]Task)}
				continue
			}
			return nil, err
		}
		project.GroupFiles[key] = gf
	}

	return project, nil
}

// Validate checks the project configuration for errors.
func (p *Project) Validate() []error {
	var errs []error

	if p.Config.Name == "" {
		errs = append(errs, fmt.Errorf("config: 'name' is required"))
	}

	if len(p.Config.Groups) == 0 {
		errs = append(errs, fmt.Errorf("config: at least one group is required"))
	}

	// Validate identifier keys (defense in depth — schema also enforces this)
	for key := range p.Config.Groups {
		if key == constants.BuiltinGroupKey {
			errs = append(errs, fmt.Errorf("config: group key %q is reserved for built-in commands", key))
		}
		if !identifierPattern.MatchString(key) {
			errs = append(errs, fmt.Errorf("config: invalid group key %q (must match %s)", key, identifierPattern))
		}
	}

	envKeys := make(map[string]bool)
	for key := range p.Config.Environments {
		if !identifierPattern.MatchString(key) {
			errs = append(errs, fmt.Errorf("config: invalid environment key %q (must match %s)", key, identifierPattern))
		}
		envKeys[key] = true
	}

	for groupKey, gf := range p.GroupFiles {
		for taskKey, task := range gf.Tasks {
			fullKey := groupKey + ":" + taskKey
			if !taskKeyPattern.MatchString(taskKey) {
				errs = append(errs, fmt.Errorf("task %s: invalid task key %q (must match %s)", fullKey, taskKey, taskKeyPattern))
			}
			if task.Name == "" {
				errs = append(errs, fmt.Errorf("task %s: 'name' is required", fullKey))
			}
			if task.Environment != "" {
				if !identifierPattern.MatchString(task.Environment) {
					errs = append(errs, fmt.Errorf("task %s: invalid environment %q (must match %s)", fullKey, task.Environment, identifierPattern))
				} else if !envKeys[task.Environment] {
					errs = append(errs, fmt.Errorf("task %s: unknown environment %q", fullKey, task.Environment))
				}
			}
			if len(task.Cmds) == 0 {
				errs = append(errs, fmt.Errorf("task %s: at least one command is required", fullKey))
			}
		}
	}

	return errs
}

// taskKeyPattern validates task keys within group files. Mirrors the JSON
// Schema pattern as defense in depth.
var taskKeyPattern = regexp.MustCompile(`^[a-z][a-z0-9_-]*(:[a-z][a-z0-9_-]*)*$`)
