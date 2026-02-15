package config

import (
	"embed"
	"encoding/json"
	"fmt"

	"github.com/santhosh-tekuri/jsonschema/v6"
	"gopkg.in/yaml.v3"
)

//go:embed schemas/*.json
var schemaFS embed.FS

// ValidateSchema validates raw YAML data against an embedded JSON Schema.
// schemaName should be "tasker.schema.json" or "tasks.schema.json".
func ValidateSchema(data []byte, schemaName string) error {
	schemaData, err := schemaFS.ReadFile("schemas/" + schemaName)
	if err != nil {
		return fmt.Errorf("loading schema %s: %w", schemaName, err)
	}

	var schemaDoc any
	if err := json.Unmarshal(schemaData, &schemaDoc); err != nil {
		return fmt.Errorf("parsing schema %s: %w", schemaName, err)
	}

	compiler := jsonschema.NewCompiler()
	if err := compiler.AddResource("schema.json", schemaDoc); err != nil {
		return fmt.Errorf("adding schema resource: %w", err)
	}

	schema, err := compiler.Compile("schema.json")
	if err != nil {
		return fmt.Errorf("compiling schema %s: %w", schemaName, err)
	}

	var doc any
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return fmt.Errorf("parsing YAML: %w", err)
	}

	doc = convertYAMLToJSON(doc)

	if err := schema.Validate(doc); err != nil {
		return fmt.Errorf("schema validation: %w", err)
	}

	return nil
}

// convertYAMLToJSON converts YAML-decoded data to JSON-compatible types.
// YAML decodes maps as map[string]any but may also use map[any]any.
func convertYAMLToJSON(v any) any {
	switch val := v.(type) {
	case map[string]any:
		result := make(map[string]any, len(val))
		for k, v := range val {
			result[k] = convertYAMLToJSON(v)
		}
		return result
	case map[any]any:
		result := make(map[string]any, len(val))
		for k, v := range val {
			result[fmt.Sprintf("%v", k)] = convertYAMLToJSON(v)
		}
		return result
	case []any:
		result := make([]any, len(val))
		for i, v := range val {
			result[i] = convertYAMLToJSON(v)
		}
		return result
	default:
		return v
	}
}

// SchemaFS returns the embedded schema filesystem for external use.
func SchemaFS() embed.FS {
	return schemaFS
}
