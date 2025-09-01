package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
)

type RelationshipType int
const (
	Contains RelationshipType = iota
	IsA 
	References
)

var relationshipName = map[RelationshipType]string {
	Contains: "contains",
	IsA: "is_a",
	References: "references",
}

func (r RelationshipType) String() string {
	return relationshipName[r]
}

type CoreEntityRelationship struct {
	RelatedEntity string `json:"related_entity" yaml:"related_entity"`
	RelationshipType string `json:"relationship_type" yaml:"relationship_type"`
	Description string `json:"description" yaml:"description"`
	RelatedAttribute string `json:"related_attribute" yaml:"related_attribute"`
}

type CoreEntityAttribute struct {
	ID	uuid.UUID `json:"id" yaml:"id"`
	Name string `json:"name" yaml:"name"`
	Description string `json:"description" yaml:"description"`
	Owner string `json:"owner" yaml:"owner"`
	Version string `json:"version" yaml:"version"`
	DataType string `json:"data_type" yaml:"data_type"`
	LifecycleState string `json:"lifecyle_state" yaml:"lifecycle_state"`
}

type CoreEntity struct {
	ID uuid.UUID `json:"id" yaml:"id"`
	Name string `json:"name" yaml:"name"`
	Description string `json:"description" yaml:"description"`
	Owner string `json:"owner" yaml:"owner"`
	Version string `json:"version" yaml:"version"`
	SystemOfRecord string `json:"system_of_record" yaml:"system_of_record"`
	LifecycleState string `json:"lifecyle_state" yaml:"lifecycle_state"`
	Attributes []CoreEntityAttribute `json:"attributes" yaml:"attributes"`
	Relationships []CoreEntityRelationship `json:"relationships" yaml:"relationships"` 
}

func main() {
	file, err := os.ReadFile("entities/workspace.yaml")
	if err != nil {
		log.Printf("Unable to read yaml file: %s", err)
	}
	if len(file) == 0 {
		log.Printf("File read with 0 bytes")
	}

	var entityMap map[string]CoreEntity
	err = yaml.Unmarshal(file, &entityMap)
	if err != nil {
		log.Printf("Unable to parse yaml: %s", err)
	}

	for _, value := range entityMap {
		value.ID = uuid.New()
		for i := range value.Attributes {
			 value.Attributes[i].ID = uuid.New()
		}
		pretty, err := json.MarshalIndent(value, "", "\t")
		if err != nil {
			log.Printf("Unable to convert to json: %s", err)
		}
		fmt.Println(string(pretty))
	}
}
