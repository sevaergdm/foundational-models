package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func ParseYaml(filepath string) (map[string]CoreEntity, error) {
	file, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("Unable to read file: %s", err)
	}
	if len(file) == 0 {
		return nil, fmt.Errorf("File was empty")
	}

	var entityMap map[string]CoreEntity
	err = yaml.Unmarshal(file, &entityMap)
	if err != nil {
		return nil, fmt.Errorf("Invalid format: %s", err)
	}

	return entityMap, nil
}
