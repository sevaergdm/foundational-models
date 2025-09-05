package main

import (
	"log"
	"net/http"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

const canonicalSchemaPath = "schema/foundational_model_schema.json"

type apiConfig struct {
	entitiesCache       map[string]FoundationalModel
	port                string
	compiledCanonicalSchema	*jsonschema.Schema
}

func createCompiledSchema() (*jsonschema.Schema, error) {
	compiler := jsonschema.NewCompiler()
	return compiler.Compile(canonicalSchemaPath)
}

func main() {
	apiCfg := &apiConfig{
		entitiesCache: make(map[string]FoundationalModel),
	}

	err := apiCfg.loadEntities("entities")
	if err != nil {
		log.Fatalf("Failed to load entities: %v", err)
	}

	err = apiCfg.graphBuilder()
	if err != nil {
		log.Fatalf("Failed to build graph: %v", err)
	}
}

func oldMain() {
	const filepathRoot = "."
	mux := http.NewServeMux()

	compiledSchema, err := createCompiledSchema()
	if err != nil {
		log.Fatalf("Unable to compile canonical schema: %v", err)
	}

	apiCfg := &apiConfig{
		entitiesCache:       make(map[string]FoundationalModel),
		port:                "8080",
		compiledCanonicalSchema: compiledSchema,
	}

	err = apiCfg.loadEntities("entities")
	if err != nil {
		log.Fatalf("Failed to load entities: %v", err)
	}

	mux.HandleFunc("GET /api/entities", apiCfg.handlerGetEntities)
	mux.HandleFunc("GET /api/entities/{entityName}", apiCfg.handlerGetEntity)
	mux.HandleFunc("POST /api/validate", apiCfg.handlerValidateEntitySchema)
	mux.HandleFunc("POST /api/entities", apiCfg.handlerCreateEntity)
	mux.HandleFunc("PUT /api/entities/{entityName}", apiCfg.handlerUpdateEntity)

	server := &http.Server{
		Handler: mux,
		Addr:    ":" + apiCfg.port,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, apiCfg.port)
	log.Fatal(server.ListenAndServe())
}
