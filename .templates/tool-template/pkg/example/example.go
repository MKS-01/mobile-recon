// Package example provides core functionality for the {{TOOL_NAME_TITLE}}.
// This is where you implement the main logic for your tool.
package example

import (
	"fmt"
)

// ExampleFunction is a sample function showing the package structure.
// Replace this with your actual implementation.
func ExampleFunction(input string) (string, error) {
	if input == "" {
		return "", fmt.Errorf("input cannot be empty")
	}

	// TODO: Implement your logic here
	result := fmt.Sprintf("Processed: %s", input)

	return result, nil
}

// ExampleStruct demonstrates how to structure types in your package.
type ExampleStruct struct {
	Field1 string
	Field2 int
	Field3 bool
}

// NewExampleStruct creates a new instance of ExampleStruct.
func NewExampleStruct(field1 string, field2 int) *ExampleStruct {
	return &ExampleStruct{
		Field1: field1,
		Field2: field2,
		Field3: false,
	}
}

// Method is an example method on the struct.
func (e *ExampleStruct) Method() error {
	// TODO: Implement your method logic
	return nil
}
