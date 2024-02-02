# go-json-schema-toolkit
Go JSON Schema Toolkit is a Go library for creating and validating JSON schemas. It supports dynamic schema creation, extension, and validation against these schemas, ensuring data integrity and compliance with predefined structures.

## Installation

To install this project, you need to have Go installed on your machine. You can download it from the [official website](https://golang.org/dl/).

Once you have Go installed, you can clone this repository using the following command:

```bash
git clone https://github.com/d-saurabh/go-json-schema-toolkit.git
```

Then, navigate to the project directory:
```bash
cd go-json-schema-toolkit
```

And install the dependencies:
```bash
go mod tidy
```

## Usage

First, import the `gojsonschematoolkit` package in your Go file:

```bash
import (
    gojsonschematoolkit "go-json-scheam-toolkit"
)
```

Then, you can create a new `SchemaManager` and use it to create a schema:

```bash
sm := gojsonschematoolkit.NewSchemaManager(make(map[string]interface{}))

err := sm.CreateSchema("testSchema", `{
    "properties": {
        "age": {
            "type": "integer"
        }
    }
}`)
```

Here's a complete example:

```go
package main

import (
    "fmt"
    gojsonschematoolkit "go-json-scheam-toolkit"
)

func main() {
    sm := gojsonschematoolkit.NewSchemaManager(make(map[string]interface{}))

    err := sm.CreateSchema("testSchema", `{
        "properties": {
            "age": {
                "type": "integer"
            }
        }
    }`)

    if err != nil {
        fmt.Println("Error creating schema:", err)
        return
    }

    fmt.Println("Schema created successfully")
}
```

## Testing

```bash
go test ./...
```

## Contributing

If you want to contribute to this project, please create a new issue or a pull request.

## License
This project is licensed under the MIT License - see the LICENSE file for details.