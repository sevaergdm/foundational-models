package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/sevaergdm/foundational-models/model_types"
)

func main() {
	allModels, err := loadEntities("./entities")
	if err != nil {
		log.Fatal(err)
	}

	var sb strings.Builder

	for _, v := range allModels {
		table := FoundationalModelToDBMLTable(v)
		sb.WriteString(fmt.Sprintf("%s\n", table))
	}

	for _, v := range allModels {
		rel := FoundationalModelToDBMLRelationship(v)
		sb.WriteString(fmt.Sprintf("%s\n", rel))
	}

	stringBytes := []byte(sb.String())
	err = os.WriteFile("./DBML/foundational_models.dbml", stringBytes, 0644)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully generated dbml")

	cmd := exec.Command("dbml-renderer",
		"-i", 
		"./DBML/foundational_models.dbml", 
		"-o", 
		"./DBML/foundational_models.svg",
		)

	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Println("Successfully generated svg")
}

func FoundationalModelToDBMLTable(fm model_types.FoundationalModel) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Table %s {\n", fm.Name))

	for _, attr := range fm.Attributes {
		field_name := attr.Name
		field_type := attr.DataType
		field_note := attr.Description
		
		sb.WriteString(fmt.Sprintf("  %s %s [note: '%s']\n", field_name, field_type, field_note))
	}

	sb.WriteString("}\n")
	return sb.String()
}

func FoundationalModelToDBMLRelationship(fm model_types.FoundationalModel) string {
	var sb strings.Builder
	
	for _, rel := range fm.Relationships {
		var dbmlRel string
		switch rel.RelationshipCardinality {
		case "one_to_many":
			dbmlRel = "<"
		case "many_to_one":
			dbmlRel = ">"
		case "one_to_one":
			dbmlRel = "-"
		case "many_to_many":
			dbmlRel = "<>"
		default:
			dbmlRel = "-"
		}	

		source := fmt.Sprintf("%s.%s", fm.Name, rel.RelatedAttribute)
		target := fmt.Sprintf("%s.%s", rel.RelatedEntity, rel.RelatedAttribute)

		sb.WriteString(fmt.Sprintf("Ref %s: %s %s %s\n", rel.RelationshipType, source, dbmlRel, target))
	}

	return sb.String()
}

func loadEntities(path string) (map[string]model_types.FoundationalModel, error) {
	output := make(map[string]model_types.FoundationalModel)
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
		return nil, err
	}

	for _, file := range files {
		log.Printf("Processing file: %s", file)
		entity, err := ParseJSON(file)
		if err != nil {
			return nil, err
		}
		output[entity.Name] = entity
	}
	return output, nil
}

func ParseJSON(filepath string) (model_types.FoundationalModel, error) {
	file, err := os.ReadFile(filepath)
	if err != nil {
		return model_types.FoundationalModel{}, fmt.Errorf("Unable to read file: %s", err)
	}
	if len(file) == 0 {
		return model_types.FoundationalModel{}, fmt.Errorf("File was empty")
	}

	var entity model_types.FoundationalModel
	err = json.Unmarshal(file, &entity)
	if err != nil {
		return model_types.FoundationalModel{}, fmt.Errorf("Invalid format: %s", err)
	}
	return entity, nil
}
