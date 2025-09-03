package main

import (
	"log"
	"os"
	"path/filepath"
)

func (cfg *apiConfig) loadEntities(path string) error {
	files := make([]string, 0)

	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if !f.IsDir() && filepath.Ext(path) == ".json" {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return err
	}

	for _, file := range files {
		log.Printf("Processing file: %s", file)
		entity, err := ParseJSON(file)
		if err != nil {
			return err
		}
		cfg.entitiesCache[entity.Name] = entity
	}
	return nil
}
