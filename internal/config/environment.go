package config

// Environment represents a deployment environment (dev, test, prod, etc.).
type Environment struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Dotenv      []string `yaml:"dotenv"`
}
