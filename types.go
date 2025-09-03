package main

type RelationshipType int
type EntityType int

const (
	Contains RelationshipType = iota
	IsA
	References
)

const (
	Foundational EntityType = iota
	Associative
)

var relationshipName = map[RelationshipType]string{
	Contains:   "contains",
	IsA:        "is_a",
	References: "references",
}

var entityTypeName = map[EntityType]string{
	Foundational: "foundational",
	Associative:  "associative",
}

func (r RelationshipType) String() string {
	return relationshipName[r]
}

func (e EntityType) String() string {
	return entityTypeName[e]
}

type Relationship struct {
	RelatedEntity    string `json:"related_entity"`
	RelationshipType string `json:"relationship_type"`
	RelatedAttribute string `json:"related_attribute"`
}

type Attribute struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	DataType    string         `json:"data_type"`
	DataFormat  any            `json:"data_format,omitempty"`
	Items       map[string]any `json:"items,omitempty"`
}

type FoundationalModel struct {
	Name           string         `json:"name"`
	Description    string         `json:"description"`
	Owner          string         `json:"owner"`
	Version        string         `json:"version"`
	SystemOfRecord string         `json:"system_of_record"`
	LifecycleState string         `json:"lifecyle_state"`
	PrimaryKey     any            `json:"primary_key"`
	EntityType     string         `json:"entity_type"`
	Attributes     []Attribute    `json:"attributes"`
	Relationships  []Relationship `json:"relationships,omitempty"`
}
