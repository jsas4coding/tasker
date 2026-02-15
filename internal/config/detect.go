package config

import (
	"os"
	"path/filepath"
)

// DetectedManager represents a detected package manager.
type DetectedManager struct {
	Key         string
	Name        string
	Description string
	File        string
}

var detectors = []struct {
	file        string
	key         string
	name        string
	description string
}{
	{"package.json", "npm", "Node.js", "Node.js package management and scripts"},
	{"composer.json", "composer", "Composer", "PHP dependency management"},
	{"go.mod", "go", "Go", "Go build and dependency management"},
	{"Cargo.toml", "cargo", "Cargo", "Rust build and dependency management"},
	{"pyproject.toml", "python", "Python", "Python package management"},
}

// DetectManagers scans a directory for known package manager files.
func DetectManagers(dir string) []DetectedManager {
	var found []DetectedManager
	for _, d := range detectors {
		path := filepath.Join(dir, d.file)
		if _, err := os.Stat(path); err == nil {
			found = append(found, DetectedManager{
				Key:         d.key,
				Name:        d.name,
				Description: d.description,
				File:        d.file,
			})
		}
	}
	return found
}
