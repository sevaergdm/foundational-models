package model_types

type RelationshipType int
type EntityType int
type RelationshipCardinality int

const (
	Contains RelationshipType = iota
	IsA
	References
)

const (
	Foundational EntityType = iota
	Associative
)

const (
	OneToMany RelationshipCardinality = iota
	ManyToOne
	OneToOne
	ManyToMany
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

var relationshipCardinalityName = map[RelationshipCardinality]string{
	OneToMany:  "OneToMany",
	ManyToOne:  "ManyToOne",
	OneToOne:   "OneToOne",
	ManyToMany: "ManyToMany",
}

func (r RelationshipType) String() string {
	return relationshipName[r]
}

func (e EntityType) String() string {
	return entityTypeName[e]
}

func (r RelationshipCardinality) String() string {
	return relationshipCardinalityName[r]
}

type Relationship struct {
	RelatedEntity           string `json:"related_entity"`
	RelationshipType        string `json:"relationship_type"`
	RelatedAttribute        string `json:"related_attribute"`
	RelationshipCardinality string `json:"relationship_cardinality"`
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
