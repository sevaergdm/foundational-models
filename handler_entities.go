package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func handlerGetEntities(w http.ResponseWriter, r *http.Request) {
	files := make([]string, 0)

	err := filepath.Walk("entities", func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !f.IsDir() && (filepath.Ext(path) == ".yaml" || filepath.Ext(path) == ".yml") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to traverse entities directory", err)
		return
	}

	entities := []CoreEntity{}
	for _, file := range files {
		log.Printf("Processing file: %s", file)
		entityMap, err := ParseYaml(file)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Unable to parse yaml", err)
			return
		}
		for _, value := range entityMap {
			entities = append(entities, value)
		}
	}
	respondWithPrettyJSON(w, http.StatusOK, entities)
}
