package main

import (
	"github.com/santhosh-tekuri/jsonschema/v6"
	"io"
)

func (cfg *apiConfig) validateFoundationalEntity(reader io.Reader) error {
	compiler := jsonschema.NewCompiler()
	canonicalSchema, err := compiler.Compile(cfg.canonicalSchemaPath)
	if err != nil {
		return err
	}

	data, err := jsonschema.UnmarshalJSON(reader)
	if err != nil {
		return err
	}

	err = canonicalSchema.Validate(data)
	if err != nil {
		return err
	}
	return nil
}
