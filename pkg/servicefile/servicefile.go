package servicefile

import "sort"

type ServiceFile struct {
	Name          string         `yaml:"name"`
	Description   string         `yaml:"description"`
	Relationships []Relationship `yaml:"relationships"`
}

type Relationship struct {
	Action      RelationshipAction `yaml:"action"`
	Name        string             `yaml:"name,omitempty"`
	Description string             `yaml:"description,omitempty"`
	Technology  string             `yaml:"technology"`
}

type RelationshipAction string

const (
	RelationshipActionUses     = "uses"
	RelationshipActionRequests = "requests"
	RelationshipActionReplies  = "replies"
	RelationshipActionSends    = "sends"
	RelationshipActionReceives = "receives"
)

func (sf *ServiceFile) Sort() {
	sort.Slice(sf.Relationships, func(i, j int) bool {
		rel1 := sf.Relationships[i]
		rel2 := sf.Relationships[j]

		if rel1.Action != rel2.Action {
			return string(rel1.Action) < string(rel2.Action)
		}

		if rel1.Name != rel2.Name {
			return rel1.Name < rel2.Name
		}

		if rel1.Technology != rel2.Technology {
			return rel1.Technology < rel2.Technology
		}

		return rel1.Description < rel2.Description
	})
}
