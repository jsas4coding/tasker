package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Group defines a task category (go, lint, test, etc.).
type Group struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

// Task represents a single task within a group file.
type Task struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Environment string   `yaml:"environment"`
	Cmds        []string `yaml:"cmds"`
	Dir         string   `yaml:"dir"`
	Deps        []string `yaml:"deps"`
	Silent      bool     `yaml:"silent"`
}

// GroupFile represents the contents of a .tasks/*.yml file.
type GroupFile struct {
	Tasks map[string]Task `yaml:"tasks"`
}

// LoadGroupFile reads and parses a .tasks/*.yml file.
func LoadGroupFile(path string) (*GroupFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading group file %s: %w", path, err)
	}

	if err := ValidateSchema(data, "tasks.schema.json"); err != nil {
		return nil, fmt.Errorf("%s %w", path, err)
	}

	var gf GroupFile
	if err := yaml.Unmarshal(data, &gf); err != nil {
		return nil, fmt.Errorf("parsing group file %s: %w", path, err)
	}

	return &gf, nil
}

// GroupFileName returns the lowercase filename for a group key.
// e.g., "go" -> "go.yml", "lint" -> "lint.yml"
func GroupFileName(key string) string {
	return key + ".yml"
}
