package gojsonschematoolkit

import (
	"encoding/json"
	"errors"
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

// const BaseSchema = `
// {
//     "$schema": "http://json-schema.org/draft-07/schema#",
//     "type": "object",
//     "properties": {
//         "id": {
//             "type": "string"
//         }
//     },
//     "required": ["id"],
//     "additionalProperties": false
// }
// `
const BaseSchema = `
{
	"definitions": {
	  "group": {
		"type": "object",
		"properties": {
		  "id": {
			"type": "string"
		  },
		  "name": {
			"type": "string"
		  }
		},
		"required": ["id", "name"]
	  },
	  "social": {
		"type": "object",
		"properties": {
		  "network": {
			"type": "string"
		  },
		  "username": {
			"type": "string"
		  }
		},
		"required": ["network", "username"]
	  },
	  "attribute": {
		"type": "object",
		"additionalProperties": true
	  }
	},
	"type": "object",
	"properties": {
	  "user": {
		"type": "object",
		"properties": {
		  "id": {
			"type": "string"
		  },
		  "name": {
			"type": "string"
		  },
		  "groups": {
			"type": "array",
			"items": {
			  "$ref": "#/definitions/group"
			}
		  },
		  "socials": {
			"type": "array",
			"items": {
			  "$ref": "#/definitions/social"
			}
		  },
		  "attributes": {
			"$ref": "#/definitions/attribute"
		  }
		},
		"additionalProperties": false
	  }
	}
  }
`

// CreateSchema creates a new schema with the given ID and additional properties.
// The additional properties are merged with the base schema to form the new schema.
// It returns an error if the schema cannot be created.
func (sm *SchemaManagerImpl) CreateSchema(schemaID string, additionalSchema string) error {
	// Unmarshal the base schema and the additional schema
	var base, additional map[string]interface{}
	if err := json.Unmarshal([]byte(BaseSchema), &base); err != nil {
		return err
	}
	if err := json.Unmarshal([]byte(additionalSchema), &additional); err != nil {
		return err
	}

	// Retrieve the "properties" field from the base schema and the additional schema
	baseProperties := base["properties"].(map[string]interface{})
	additionalProperties := additional["properties"].(map[string]interface{})

	// Check if the additional schema's properties contain a "user" property and no other properties
	if _, ok := additionalProperties["user"]; !ok {
		return errors.New("additional schema 'properties' object does not contain a 'user' property")
	}
	if len(additionalProperties) > 1 {
		return errors.New("additional schema 'properties' object contains properties other than 'user'")
	}

	// Determine which entity to extend
	for k, v := range additionalProperties {
		if _, ok := baseProperties[k]; ok {
			// If the key exists in the base schema, extend that entity
			targetProperties := baseProperties[k].(map[string]interface{})["properties"].(map[string]interface{})
			for kk, vv := range v.(map[string]interface{})["properties"].(map[string]interface{}) {
				targetProperties[kk] = vv
			}
		}
	}

	// Merge the definitions
	baseDefinitions := base["definitions"].(map[string]interface{})
	additionalDefinitions, ok := additional["definitions"].(map[string]interface{})
	if ok {
		for k, v := range additionalDefinitions {
			baseDefinitions[k] = v
		}
	}

	// Convert the final merged schema to a JSON string
	mergedSchemaJson, err := json.MarshalIndent(base, "", "  ")
	if err != nil {
		fmt.Println("Error converting merged schema to JSON:", err)
	} else {
		fmt.Println("Final merged schema:", string(mergedSchemaJson))
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
