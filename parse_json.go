package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func ParseJSON(filepath string) (FoundationalModel, error) {
	file, err := os.ReadFile(filepath)
	if err != nil {
		return FoundationalModel{}, fmt.Errorf("Unable to read file: %s", err)
	}
	if len(file) == 0 {
		return FoundationalModel{}, fmt.Errorf("File was empty")
	}

	var entity FoundationalModel
	err = json.Unmarshal(file, &entity)
	if err != nil {
		return FoundationalModel{}, fmt.Errorf("Invalid format: %s", err)
	}
	return entity, nil
}
