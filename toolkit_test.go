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

func TestSchemaManager_1(t *testing.T) {
	sm := gojsonschematoolkit.NewSchemaManager(make(map[string]interface{}))

	// Test CreateSchema
	err := sm.CreateSchema("userSchema", `{
        "properties": {
          "user": {
            "properties": {
              "email": {
                "type": "string",
                "format": "email"
              }
            },
			"additionalProperties": false
          }
        },
        "definitions": {
          "attribute": {
            "properties": {
              "DOB": {
                "type": "string",
                "format": "date"
              }
            },
			"additionalProperties": false,
            "type": "object"
          },
          "newDefinition": {
            "properties": {
              "newProperty": {
                "type": "string"
              }
            },
            "type": "object"
          }
        }
      }`)
	assert.NoError(t, err)

	testCases := []struct {
		name        string
		schemaID    string
		data        string
		expectError bool
	}{
		{
			name:     "Valid user data",
			schemaID: "userSchema",
			data: `{
                "user": {
                    "id": "1",
                    "name": "John Doe",
                    "email": "john.doe@example.com",
                    "groups": [
                        {
                            "id": "group1",
                            "name": "Group 1"
                        }
                    ],
                    "socials": [
                        {
                            "network": "Twitter",
                            "username": "johndoe"
                        }
                    ],
                    "attributes": {
                        "DOB": "1990-01-01"
                    }
                }
            }`,
			expectError: false,
		},
		{
			name:     "Invalid user data - missing email",
			schemaID: "userSchema",
			data: `{
		        "user": {
		            "id": "1",
		            "name": "John Doe",
		            "DOB": "1990-01-01",
		            "groups": [
		                {
		                    "id": "group1",
		                    "name": "Group 1"
		                }
		            ],
		            "socials": [
		                {
		                    "network": "Twitter",
		                    "username": "johndoe"
		                }
		            ],
		            "attributes": {
		                "DOB": "1990-01-01",
		                "newProperty": "newValue"
		            }
		        }
		    }`,
			expectError: true,
		},
		// Add more test cases here
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := sm.ValidateSchema(tc.schemaID, tc.data)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSchemaManager_2(t *testing.T) {
	sm := gojsonschematoolkit.NewSchemaManager(make(map[string]interface{}))

	// Test CreateSchema with userSchema
	err := sm.CreateSchema("userSchema", `{
        "properties": {
          "user": {
            "properties": {
              "email": {
                "type": "string",
                "format": "email"
              }
            },
            "additionalProperties": false
          }
        },
        "definitions": {
          "attribute": {
            "properties": {
              "DOB": {
                "type": "string",
                "format": "date"
              }
            },
            "additionalProperties": false,
            "type": "object"
          },
          "newDefinition": {
            "properties": {
              "newProperty": {
                "type": "string"
              }
            },
            "type": "object"
          }
        }
      }`)
	assert.NoError(t, err)

	// Test CreateSchema with invalid schema (properties other than 'user')
	err = sm.CreateSchema("invalidSchema", `{
        "properties": {
          "item": {
            "properties": {
              "name": {
                "type": "string"
              },
              "price": {
                "type": "number"
              }
            },
            "additionalProperties": false
          }
        }
      }`)
	assert.Error(t, err)

	// Test CreateSchema with valid schema (only 'user' property)
	err = sm.CreateSchema("validSchema", `{
        "properties": {
          "user": {
            "properties": {
              "name": {
                "type": "string"
              },
              "price": {
                "type": "number"
              }
            },
            "additionalProperties": false
          }
        }
      }`)
	assert.NoError(t, err)

	testCases := []struct {
		name        string
		schemaID    string
		data        string
		expectError bool
	}{
		// {
		// 	name:     "Valid user data",
		// 	schemaID: "userSchema",
		// 	data: `{
		//         "user": {
		//             "id": "1",
		//             "name": "John Doe",
		//             "email": "john.doe@example.com",
		//             "groups": [
		//                 {
		//                     "id": "group1",
		//                     "name": "Group 1"
		//                 }
		//             ],
		//             "socials": [
		//                 {
		//                     "network": "Twitter",
		//                     "username": "johndoe"
		//                 }
		//             ],
		//             "attributes": {
		//                 "DOB": "1990-01-01"
		//             }
		//         }
		//     }`,
		// 	expectError: false,
		// },
		// {
		// 	name:     "Invalid user data - missing email",
		// 	schemaID: "userSchema",
		// 	data: `{
		//         "user": {
		//             "id": "1",
		//             "name": "John Doe",
		//             "DOB": "1990-01-01",
		//             "groups": [
		//                 {
		//                     "id": "group1",
		//                     "name": "Group 1"
		//                 }
		//             ],
		//             "socials": [
		//                 {
		//                     "network": "Twitter",
		//                     "username": "johndoe"
		//                 }
		//             ],
		//             "attributes": {
		//                 "DOB": "1990-01-01",
		//                 "newProperty": "newValue"
		//             }
		//         }
		//     }`,
		// 	expectError: true,
		// },
		// Test case for valid data with 'validSchema'
		// {
		// 	name:     "Valid user data with validSchema",
		// 	schemaID: "validSchema",
		// 	data: `{
		//         "user": {
		//             "name": "John Doe",
		//             "price": 100
		//         }
		//     }`,
		// 	expectError: false,
		// },
		// Test case for invalid data with 'validSchema' (missing 'name')
		{
			name:     "Invalid user data with validSchema - missing name",
			schemaID: "validSchema",
			data: `{
		        "user": {
		            "price": 100
		        }
		    }`,
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := sm.ValidateSchema(tc.schemaID, tc.data)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
