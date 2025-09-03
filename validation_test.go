package main

import (
	"strings"
	"testing"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

func TestFormatValidationError(t *testing.T) {
	verr := &jsonschema.ValidationError{
		InstanceLocation: []string{"name"},
		ErrorKind: nil,
		Causes: nil,
	}

	errors := FormatValidationError(verr)

	if len(errors) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errors))
	}

	if errors[0].Path != "name" {
		t.Errorf("Expected path 'name', but got: %s", errors[0].Path)
	}

	if errors[0].Message != "(unspecified error)" {
		t.Errorf("Expected message '(unspecified error)', but got :%s", errors[0].Message)
	}
}

func addInMemorySchema(t *testing.T, schemaStr string) *apiConfig {
	t.Helper()

	compiler := jsonschema.NewCompiler()
	if err := compiler.AddResource("test-schema.json", strings.NewReader(schemaStr)); err != nil {
		t.Fatalf("Unable to add in-memory schema resource: %v", err)
	}

	compiledSchema, err := compiler.Compile("test-schema.json")
	if err != nil {
		t.Fatalf("Unable to compile schema: %v", err)
	}

	cfg := &apiConfig{
		compiledCanonicalSchema: compiledSchema,
	}
	return cfg
}

func TestValidateFoundationalEntityValid(t *testing.T) {
	schema := `{
	"$id": "test-schema.json",
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "required": ["name", "owner"],
  "properties": {
    "name": { "type": "string" },
    "owner": { "type": "string" }
  }
}`

	cfg := addInMemorySchema(t, schema)

	validJSON := `{"name":"TestEntity","owner":"client"}`
	err := cfg.validateFoundationalEntity(strings.NewReader(validJSON))
	if err != nil {
		t.Fatalf("expected no error for valid json, but got: %v", err)
	}
}
