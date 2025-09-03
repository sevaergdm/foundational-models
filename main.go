package main

import (
	"log"
	"net/http"
)

type apiConfig struct {
	entitiesCache       map[string]FoundationalModel
	canonicalSchemaPath string
	port                string
}

func main() {
	const filepathRoot = "."
	mux := http.NewServeMux()

	apiCfg := &apiConfig{
		canonicalSchemaPath: "schema/foundational_model_schema.json",
		entitiesCache:       make(map[string]FoundationalModel),
		port:                "8080",
	}

	err := apiCfg.loadEntities("entities")
	if err != nil {
		log.Fatalf("Failed to load entities: %v", err)
	}

	mux.HandleFunc("GET /api/entities", apiCfg.handlerGetEntities)
	mux.HandleFunc("GET /api/entities/{entityName}", apiCfg.handlerGetEntity)
	mux.HandleFunc("POST /api/validate", apiCfg.handlerValidateEntitySchema)

	server := &http.Server{
		Handler: mux,
		Addr:    ":" + apiCfg.port,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, apiCfg.port)
	log.Fatal(server.ListenAndServe())
}
