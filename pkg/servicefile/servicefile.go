package servicefile

import (
	"fmt"
	"os"
	"sort"

	"gopkg.in/yaml.v3"
)

const Version string = "0.1.0"

// ServiceFile represents a service file.
type ServiceFile struct {
	Version       string         `yaml:"servicefile"`
	Info          Info           `yaml:"info"`
	Relationships []Relationship `yaml:"relationships"`
}

// Info represents a info about service.
type Info struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	System      string `yaml:"system,omitempty"`
	Owner       string `yaml:"owner,omitempty"`
}

// Relationship represents a relationship between current service and external components.
type Relationship struct {
	Action      RelationshipAction `yaml:"action"`
	Name        string             `yaml:"name,omitempty"`
	Description string             `yaml:"description,omitempty"`
	Technology  string             `yaml:"technology"`
	Proto       string             `yaml:"proto,omitempty"`
}

// RelationshipAction represents an action between services.
type RelationshipAction string

const (
	RelationshipActionUses     = "uses"
	RelationshipActionRequests = "requests"
	RelationshipActionReplies  = "replies"
	RelationshipActionSends    = "sends"
	RelationshipActionReceives = "receives"
)

// Sort sorts the relationships in the service file.
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

		if rel1.Proto != rel2.Proto {
			return rel1.Proto < rel2.Proto
		}

		return rel1.Description < rel2.Description
	})
}

// Load reads and parses a ServiceFile from a YAML file at the given path.
func Load(path string) (*ServiceFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", path, err)
	}

	var sf ServiceFile
	if err := yaml.Unmarshal(data, &sf); err != nil {
		return nil, fmt.Errorf("failed to parse file %s: %w", path, err)
	}

	return &sf, nil
}
