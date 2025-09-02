package main

import (
	"encoding/json"
	"fmt"
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
	Name           string         `json:"name" yaml:"name"`
	Description    string         `json:"description" yaml:"description"`
	Owner          string         `json:"owner" yaml:"owner"`
	Version        string         `json:"version" yaml:"version"`
	DataType       string         `json:"data_type" yaml:"data_type"`
	DataFormat     any            `json:"data_format,omitempty" yaml:"data_format,omitempty"`
	LifecycleState string         `json:"lifecyle_state" yaml:"lifecycle_state"`
	Items          map[string]any `json:"items,omitempty" yaml:"items,omitempty"`
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

type ConverterOptions struct {
	AdditionalProperties bool
	FormatMap map[string]string
	AllAttributesRequired bool
}

func DefaultConverterOptions() ConverterOptions {
	return ConverterOptions{
		AdditionalProperties: true,
		AllAttributesRequired: true,
		FormatMap: map[string]string{
			"uuid": "uuid",
			"email": "email",
			"timestamp": "date-time",
			"date-time": "date-time",
		},
	}
}

func (cfg *apiConfig) EntityToJSONSchemaBytes(entity CoreEntity, opts ConverterOptions) ([]byte, error) {
	schema, err := cfg.EntityToJSONSchema(entity, opts)
	if err != nil {
		return nil, err
	}

	return json.MarshalIndent(schema, "", "\t")
}

func (cfg *apiConfig) EntityToJSONSchema(entity CoreEntity, opts ConverterOptions) (map[string]any, error) {
	properties := make(map[string]any, len(entity.Attributes))
	required := make([]string, 0, len(entity.Attributes))
	defs := make(map[string]any)

	for _, attr := range entity.Attributes {
		propertySchema := attributeToSchema(attr, opts)
		properties[attr.Name] = propertySchema

		if opts.AllAttributesRequired {
			required = append(required, attr.Name)
		}
	}

	for _, rel := range entity.Relationships {
		switch rel.RelationshipType {
		case "contains":
			properties[rel.RelatedAttribute] = map[string]any{
				"type": "array",
				"description": rel.Description,
				"items": map[string]any{
					"$ref": "#/$defs/" + rel.RelatedEntity,
				},
			}

			relatedEntity, ok := cfg.entitiesCache[rel.RelatedEntity]
			if ok {
				defSchema, err := cfg.EntityToJSONSchema(relatedEntity, opts)
				if err != nil {
					return nil, err
				}
				defs[rel.RelatedEntity] = defSchema
			}
			case "references":
			default:
				return nil, fmt.Errorf("Unsupported relationship type: %s", rel.RelationshipType)
		}
	}


	schema := map[string]any{
		"$schema": "https://json-schema.org/draft/2020-12/schema",
		"title": entity.Name,
		"description": entity.Description,
		"type": "object",
		"additionalProperties": opts.AdditionalProperties,
		"properties": properties,
	}

	if len(required) > 0 {
		schema["required"] = required
	}

	if len(defs) > 0 {
		schema["$defs"] = defs
	}

	return schema, nil
}

func attributeToSchema(attr CoreEntityAttribute, opts ConverterOptions) map[string]any {
	property := map[string]any{
		"description": attr.Description,
	}

	switch normalizeType(attr.DataType) {
	case "string":
		property["type"] = "string"
		addFormatIfAny(property, attr.DataFormat, opts)
	case "boolean":
		property["type"] = "boolean"
	case "integer":
		property["type"] = "integer"
	case "number":
		property["type"] = "number"
	case "array":
		property["type"] = "array"
		if attr.Items != nil {
			itemsNorm := normalizeYAMLValue(attr.Items)
			property["items"] = itemsNorm
		} else {
			property["items"] = map[string]any{}
		}
	case "object":
		property["type"] = "object"
		property["additionalProperties"] = true
	default:
		property["type"] = "string"
	}
	
	return property
}

func normalizeType(t string) string {
	switch t {
	case "string", "boolean", "integer", "number", "array", "object":
		return t
	}

	switch t {
	case "bool":
		return "boolean"
	case "int":
		return "integer"
	case "float", "double", "decimal":
		return "number"
	}
	return t
}

func addFormatIfAny(property map[string]any, dataFormat any, opts ConverterOptions) {
	if dataFormat == nil {
		return
	}

	s, ok := dataFormat.(string)
	if ok && s != "" {
		mapped, found := opts.FormatMap[s]
		if found {
			property["format"] = mapped
			return
		}
		property["format"] = s
	}
}

func normalizeYAMLValue(val any) any {
	switch typ := val.(type) {
	case map[string]any:
		out := make(map[string]any, len(typ))
		for k, v := range typ {
			out[k] = normalizeYAMLValue(v)
		}
		return out
	case map[any]any:
		out := make(map[string]any, len(typ))
		for k, v := range typ {
			out[fmt.Sprint(k)] = normalizeYAMLValue(v)
		}
		return out
	case []any:
		out := make([]any, len(typ))
		for i := range typ {
			out[i] = normalizeYAMLValue(typ[i])
		}
		return out
	default:
		return typ
	} 
}
