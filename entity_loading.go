package main

import (
	"log"
	"maps"
	"os"
	"path/filepath"
)

func (cfg *apiConfig) loadEntities(path string) error {
	files := make([]string, 0)

	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if !f.IsDir() && (filepath.Ext(path) == ".yaml" || filepath.Ext(path) == ".yml") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return err
	}

	for _, file := range files {
		log.Printf("Processing file: %s", file)
		entityMap, err := ParseYaml(file)
		if err != nil {
			return err
		}
		maps.Copy(cfg.entitiesCache, entityMap)
	}
	return nil
}
