package main

import (
	"fmt"
	"log"
	"net/http"
)

type apiConfig struct {
	entitiesCache map[string]CoreEntity
}

func main() {
	const port = "8080"
	const filepathRoot = "."

	mux := http.NewServeMux()

	apiCfg := &apiConfig{
		entitiesCache: make(map[string]CoreEntity),
	}

	err := apiCfg.loadEntities("entities")
	if err != nil {
		log.Fatalf("Failed to load entities: %v", err)
	}

	opts := DefaultConverterOptions()
	workspace := apiCfg.entitiesCache["Workspace"]
	schemaBytes, err := apiCfg.EntityToJSONSchemaBytes(workspace, opts)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(schemaBytes))

	mux.HandleFunc("GET /api/entities", apiCfg.handlerGetEntities)
	mux.HandleFunc("GET /api/entities/{entityName}", apiCfg.handlerGetEntity)

	server := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
