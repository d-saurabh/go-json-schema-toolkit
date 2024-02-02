package gojsonschematoolkit_test

import (
	gojsonschematoolkit "go-json-scheam-toolkit"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestSchemaManager ...
func TestSchemaManager(t *testing.T) {
	sm := gojsonschematoolkit.NewSchemaManager(make(map[string]interface{}))

	// Test CreateSchema
	err := sm.CreateSchema("testSchema", `{
        "properties": {
            "age": {
                "type": "integer"
            }
        },
        "additionalProperties": false
    }`)
	assert.NoError(t, err)

	testCases := []struct {
		name      string
		schemaID  string
		data      string
		expectErr bool
	}{
		{
			name:      "Valid data",
			schemaID:  "testSchema",
			data:      `{"id": "123", "age": 30}`,
			expectErr: false,
		},
		{
			name:      "Invalid data (incorrect type)",
			schemaID:  "testSchema",
			data:      `{"id": "123", "age": "thirty"}`,
			expectErr: true,
		},
		{
			name:      "Invalid data (unexpected property)",
			schemaID:  "testSchema",
			data:      `{"id": "123", "age": 30, "name": "John"}`,
			expectErr: true,
		},
		{
			name:      "Non-existent schema",
			schemaID:  "nonExistentSchema",
			data:      `{"id": "123", "age": 30}`,
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := sm.ValidateSchema(tc.schemaID, tc.data)
			if tc.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
