package main

import (
	"log"
	"net/http"
)

type RelationshipType int

const (
	Contains RelationshipType = iota
	IsA
	References
)

var relationshipName = map[RelationshipType]string{
	Contains:   "contains",
	IsA:        "is_a",
	References: "references",
}

func (r RelationshipType) String() string {
	return relationshipName[r]
}

type CoreEntityRelationship struct {
	RelatedEntity    string `json:"related_entity" yaml:"related_entity"`
	RelationshipType string `json:"relationship_type" yaml:"relationship_type"`
	Description      string `json:"description" yaml:"description"`
	RelatedAttribute string `json:"related_attribute" yaml:"related_attribute"`
}

type CoreEntityAttribute struct {
	//	ID	uuid.UUID `json:"id" yaml:"id"`
	Name           string                 `json:"name" yaml:"name"`
	Description    string                 `json:"description" yaml:"description"`
	Owner          string                 `json:"owner" yaml:"owner"`
	Version        string                 `json:"version" yaml:"version"`
	DataType       string                 `json:"data_type" yaml:"data_type"`
	LifecycleState string                 `json:"lifecyle_state" yaml:"lifecycle_state"`
	Items          map[string]interface{} `json:"items,omitempty" yaml:"items,omitempty"`
}

type CoreEntity struct {
	//	ID uuid.UUID `json:"id" yaml:"id"`
	Name           string                   `json:"name" yaml:"name"`
	Description    string                   `json:"description" yaml:"description"`
	Owner          string                   `json:"owner" yaml:"owner"`
	Version        string                   `json:"version" yaml:"version"`
	SystemOfRecord string                   `json:"system_of_record" yaml:"system_of_record"`
	LifecycleState string                   `json:"lifecyle_state" yaml:"lifecycle_state"`
	Attributes     []CoreEntityAttribute    `json:"attributes" yaml:"attributes"`
	Relationships  []CoreEntityRelationship `json:"relationships,omitempty" yaml:"relationships,omitempty"`
}

func main() {
	const port = "8080"
	const filepathRoot = "."

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/entities", handlerGetEntities)

	server := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
