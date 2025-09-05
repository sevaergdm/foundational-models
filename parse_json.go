package main

import (
	"encoding/json"
	"fmt"
	"github.com/sevaergdm/foundational-models/model_types"
	"os"
)

func ParseJSON(filepath string) (model_types.FoundationalModel, error) {
	file, err := os.ReadFile(filepath)
	if err != nil {
		return model_types.FoundationalModel{}, fmt.Errorf("Unable to read file: %s", err)
	}
	if len(file) == 0 {
		return model_types.FoundationalModel{}, fmt.Errorf("File was empty")
	}

	var entity model_types.FoundationalModel
	err = json.Unmarshal(file, &entity)
	if err != nil {
		return model_types.FoundationalModel{}, fmt.Errorf("Invalid format: %s", err)
	}
	return entity, nil
}
