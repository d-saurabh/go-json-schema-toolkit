package gojsonschematoolkit

import (
	"encoding/json"
	"fmt"

	"github.com/xeipuuv/gojsonschema"
)

// SchemaManager is an interface for managing JSON schemas.
// It provides methods for creating and validating schemas.
type SchemaManager interface {
	// CreateSchema creates a new schema with the given ID and additional properties.
	// The additional properties are merged with the base schema to form the new schema.
	// It returns an error if the schema cannot be created.
	CreateSchema(schemaID string, additionalSchema string) error

	// ValidateSchema validates the given data against the schema with the given ID.
	// It returns an error if the data is not valid according to the schema.
	ValidateSchema(schemaID string, data string) error
}

// SchemaManagerImpl is a concrete implementation of the SchemaManager interface.
// It holds a map of schemas identified by a string ID.
type SchemaManagerImpl struct {
	schemaMap map[string]interface{}
}

// NewSchemaManager creates a new SchemaManagerImpl
func NewSchemaManager(schemaMap map[string]interface{}) SchemaManager {
	return &SchemaManagerImpl{
		schemaMap: schemaMap,
	}
}

const BaseSchema = `
{
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "properties": {
        "id": {
            "type": "string"
        }
    },
    "required": ["id"],
    "additionalProperties": false
}
`

// CreateSchema creates a new schema with the given ID and additional properties.
// The additional properties are merged with the base schema to form the new schema.
// It returns an error if the schema cannot be created.
func (sm *SchemaManagerImpl) CreateSchema(schemaID string, additionalSchema string) error {
	var base map[string]interface{}
	var props map[string]interface{}

	// Unmarshal the base schema and the additional schema
	if err := json.Unmarshal([]byte(BaseSchema), &base); err != nil {
		return err
	}
	if err := json.Unmarshal([]byte(additionalSchema), &props); err != nil {
		return err
	}

	// Merge the base schema with the additional properties
	baseProperties := base["properties"].(map[string]interface{})
	additionalProperties := props["properties"].(map[string]interface{})
	for k, v := range additionalProperties {
		baseProperties[k] = v
	}

	// Add the new schema to the schema map
	sm.schemaMap[schemaID] = base

	return nil
}

// ValidateSchema validates the given data against the schema with the given ID.
// It returns an error if the data is not valid according to the schema.
func (sm *SchemaManagerImpl) ValidateSchema(schemaID string, data string) error {
	// Get the schema from the schema map
	schema, ok := sm.schemaMap[schemaID]
	if !ok {
		return fmt.Errorf("schema not found: %s", schemaID)
	}

	// Convert the schema to a JSON string
	schemaLoader := gojsonschema.NewGoLoader(schema)

	// Load the data
	dataLoader := gojsonschema.NewStringLoader(data)

	// Validate the data against the schema
	result, err := gojsonschema.Validate(schemaLoader, dataLoader)
	if err != nil {
		return err
	}

	if !result.Valid() {
		return fmt.Errorf("invalid data: %v", result.Errors())
	}

	return nil
}
