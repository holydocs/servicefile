package servicefile

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		yamlContent string
		want        *ServiceFile
		wantErr     bool
		errContains string
	}{
		{
			name: "valid servicefile",
			yamlContent: `
servicefile: "0.1.0"
info:
    name: "test-service"
    description: "A test service"
relationships:
  - action: "uses"
    participant: "database"
    description: "Uses PostgreSQL database"
    technology: "postgresql"
  - action: "sends"
    participant: "notifications"
    description: "Sends email notifications"
    technology: "smtp"
`,
			want: &ServiceFile{
				Version: "0.1.0",
				Info: Info{
					Name:        "test-service",
					Description: "A test service",
				},
				Relationships: []Relationship{
					{
						Action:      "uses",
						Participant: "database",
						Description: "Uses PostgreSQL database",
						Technology:  "postgresql",
					},
					{
						Action:      "sends",
						Participant: "notifications",
						Description: "Sends email notifications",
						Technology:  "smtp",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "servicefile with system",
			yamlContent: `
servicefile: "0.1.0"
info:
    name: "user-service"
    description: "Handles user authentication and profiles"
    system: "e-commerce-platform"
relationships:
  - action: "uses"
    participant: "database"
    technology: "postgresql"
`,
			want: &ServiceFile{
				Version: "0.1.0",
				Info: Info{
					Name:        "user-service",
					Description: "Handles user authentication and profiles",
					System:      "e-commerce-platform",
				},
				Relationships: []Relationship{
					{
						Action:      "uses",
						Participant: "database",
						Technology:  "postgresql",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "servicefile with owner",
			yamlContent: `
servicefile: "0.1.0"
info:
    name: "user-service"
    description: "Handles user authentication and profiles"
    owner: "team-auth"
relationships:
  - action: "uses"
    participant: "database"
    technology: "postgresql"
`,
			want: &ServiceFile{
				Version: "0.1.0",
				Info: Info{
					Name:        "user-service",
					Description: "Handles user authentication and profiles",
					Owner:       "team-auth",
				},
				Relationships: []Relationship{
					{
						Action:      "uses",
						Participant: "database",
						Technology:  "postgresql",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "minimal servicefile",
			yamlContent: `
servicefile: 0.1.0
info:
    name: "minimal-service"
`,
			want: &ServiceFile{
				Version: "0.1.0",
				Info: Info{
					Name:        "minimal-service",
					Description: "",
				},
			},
			wantErr: false,
		},
		{
			name: "servicefile with all relationship actions",
			yamlContent: `
servicefile: 0.1.0
info:
    name: "complete-service"
    description: "Service with all relationship types"
relationships:
  - action: "uses"
    technology: "redis"
  - action: "requests"
    technology: "http"
  - action: "replies"
    technology: "grpc"
  - action: "sends"
    technology: "kafka"
  - action: "receives"
    technology: "rabbitmq"
`,
			want: &ServiceFile{
				Version: "0.1.0",
				Info: Info{
					Name:        "complete-service",
					Description: "Service with all relationship types",
				},
				Relationships: []Relationship{
					{Action: "uses", Technology: "redis"},
					{Action: "requests", Technology: "http"},
					{Action: "replies", Technology: "grpc"},
					{Action: "sends", Technology: "kafka"},
					{Action: "receives", Technology: "rabbitmq"},
				},
			},
			wantErr: false,
		},
		{
			name: "servicefile with proto field",
			yamlContent: `
servicefile: 0.1.0
info:
    name: "api-service"
    description: "API service with protocol specifications"
relationships:
  - action: "uses"
    participant: "database"
    description: "Uses PostgreSQL database"
    technology: "postgresql"
    proto: "tcp"
  - action: "requests"
    participant: "auth-service"
    description: "Makes HTTP requests to authentication service"
    technology: "auth-service"
    proto: "http"
  - action: "replies"
    description: "Provides gRPC APIs"
    technology: "grpc-server"
    proto: "grpc"
  - action: "sends"
    participant: "events"
    description: "Sends events to message queue"
    technology: "kafka"
    proto: "tcp"
`,
			want: &ServiceFile{
				Version: "0.1.0",
				Info: Info{
					Name:        "api-service",
					Description: "API service with protocol specifications",
				},
				Relationships: []Relationship{
					{
						Action:      "uses",
						Participant: "database",
						Description: "Uses PostgreSQL database",
						Technology:  "postgresql",
						Proto:       "tcp",
					},
					{
						Action:      "requests",
						Participant: "auth-service",
						Description: "Makes HTTP requests to authentication service",
						Technology:  "auth-service",
						Proto:       "http",
					},
					{
						Action:      "replies",
						Description: "Provides gRPC APIs",
						Technology:  "grpc-server",
						Proto:       "grpc",
					},
					{
						Action:      "sends",
						Participant: "events",
						Description: "Sends events to message queue",
						Technology:  "kafka",
						Proto:       "tcp",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "servicefile with tags",
			yamlContent: `
servicefile: 0.1.0
info:
    name: "tagged-service"
    description: "A service with tags"
    tags: ["auth", "user-management", "microservice"]
relationships:
  - action: "uses"
    participant: "database"
    technology: "postgresql"
`,
			want: &ServiceFile{
				Version: "0.1.0",
				Info: Info{
					Name:        "tagged-service",
					Description: "A service with tags",
					Tags:        []string{"auth", "user-management", "microservice"},
				},
				Relationships: []Relationship{
					{
						Action:      "uses",
						Participant: "database",
						Technology:  "postgresql",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "servicefile with relationship tags",
			yamlContent: `
servicefile: 0.1.0
info:
    name: "relationship-tagged-service"
    description: "A service with tagged relationships"
relationships:
  - action: "uses"
    participant: "database"
    description: "Uses PostgreSQL database"
    technology: "postgresql"
    proto: "tcp"
    tags: ["persistence", "data-store", "critical"]
  - action: "requests"
    participant: "auth-service"
    description: "Makes HTTP requests to authentication service"
    technology: "auth-service"
    proto: "http"
    tags: ["security", "authentication"]
  - action: "sends"
    participant: "events"
    description: "Sends events to message queue"
    technology: "kafka"
    proto: "tcp"
    tags: ["messaging", "async"]
`,
			want: &ServiceFile{
				Version: "0.1.0",
				Info: Info{
					Name:        "relationship-tagged-service",
					Description: "A service with tagged relationships",
				},
				Relationships: []Relationship{
					{
						Action:      "uses",
						Participant: "database",
						Description: "Uses PostgreSQL database",
						Technology:  "postgresql",
						Proto:       "tcp",
						Tags:        []string{"persistence", "data-store", "critical"},
					},
					{
						Action:      "requests",
						Participant: "auth-service",
						Description: "Makes HTTP requests to authentication service",
						Technology:  "auth-service",
						Proto:       "http",
						Tags:        []string{"security", "authentication"},
					},
					{
						Action:      "sends",
						Participant: "events",
						Description: "Sends events to message queue",
						Technology:  "kafka",
						Proto:       "tcp",
						Tags:        []string{"messaging", "async"},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "servicefile with external relationships",
			yamlContent: `
servicefile: "0.1.0"
info:
    name: "external-service"
    description: "A service with external relationships"
relationships:
  - action: "requests"
    participant: "ExternalAPI"
    description: "Requests data from external third-party API"
    technology: "http"
    external: true
  - action: "uses"
    participant: "InternalDatabase"
    description: "Uses internal database"
    technology: "postgresql"
    external: false
`,
			want: &ServiceFile{
				Version: "0.1.0",
				Info: Info{
					Name:        "external-service",
					Description: "A service with external relationships",
				},
				Relationships: []Relationship{
					{
						Action:      "requests",
						Participant: "ExternalAPI",
						Description: "Requests data from external third-party API",
						Technology:  "http",
						External:    true,
					},
					{
						Action:      "uses",
						Participant: "InternalDatabase",
						Description: "Uses internal database",
						Technology:  "postgresql",
						External:    false,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "servicefile with person relationships",
			yamlContent: `
servicefile: "0.1.0"
info:
    name: "user-service"
    description: "A service that replies to people"
relationships:
  - action: "replies"
    participant: "User"
    description: "Replies to user requests via web interface"
    technology: "http"
    person: true
  - action: "replies"
    participant: "Admin"
    description: "Replies to admin requests via API"
    technology: "grpc"
    person: true
  - action: "replies"
    participant: "OtherService"
    description: "Replies to other service requests"
    technology: "grpc"
    person: false
`,
			want: &ServiceFile{
				Version: "0.1.0",
				Info: Info{
					Name:        "user-service",
					Description: "A service that replies to people",
				},
				Relationships: []Relationship{
					{
						Action:      "replies",
						Participant: "User",
						Description: "Replies to user requests via web interface",
						Technology:  "http",
						Person:      true,
					},
					{
						Action:      "replies",
						Participant: "Admin",
						Description: "Replies to admin requests via API",
						Technology:  "grpc",
						Person:      true,
					},
					{
						Action:      "replies",
						Participant: "OtherService",
						Description: "Replies to other service requests",
						Technology:  "grpc",
						Person:      false,
					},
				},
			},
			wantErr: false,
		},
		{
			name:        "invalid yaml",
			yamlContent: `name: "test" invalid: yaml: content`,
			want:        nil,
			wantErr:     true,
			errContains: "failed to parse file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			tmpFile := filepath.Join(tmpDir, "servicefile.yaml")

			err := os.WriteFile(tmpFile, []byte(tt.yamlContent), 0644)
			require.NoError(t, err)

			got, err := Load(tmpFile)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Nil(t, got)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestSort(t *testing.T) {
	tests := []struct {
		name     string
		input    *ServiceFile
		expected *ServiceFile
	}{
		{
			name: "sort by action, name, technology, proto, description",
			input: &ServiceFile{
				Version: Version,
				Info: Info{
					Name: "test-service",
				},
				Relationships: []Relationship{
					{Action: "sends", Participant: "b", Technology: "kafka", Proto: "tcp", Description: "second"},
					{Action: "uses", Participant: "a", Technology: "postgres", Proto: "tcp", Description: "first"},
					{Action: "sends", Participant: "a", Technology: "kafka", Proto: "udp", Description: "first"},
					{Action: "uses", Participant: "a", Technology: "redis", Proto: "tcp", Description: "first"},
					{Action: "sends", Participant: "a", Technology: "kafka", Proto: "tcp", Description: "first"},
				},
			},
			expected: &ServiceFile{
				Version: Version,
				Info: Info{
					Name: "test-service",
				},
				Relationships: []Relationship{
					{Action: "sends", Participant: "a", Technology: "kafka", Proto: "tcp", Description: "first"},
					{Action: "sends", Participant: "a", Technology: "kafka", Proto: "udp", Description: "first"},
					{Action: "sends", Participant: "b", Technology: "kafka", Proto: "tcp", Description: "second"},
					{Action: "uses", Participant: "a", Technology: "postgres", Proto: "tcp", Description: "first"},
					{Action: "uses", Participant: "a", Technology: "redis", Proto: "tcp", Description: "first"},
				},
			},
		},
		{
			name: "empty relationships",
			input: &ServiceFile{
				Version: Version,
				Info: Info{
					Name: "empty-service",
				},
				Relationships: []Relationship{},
			},
			expected: &ServiceFile{
				Version: Version,
				Info: Info{
					Name: "empty-service",
				},
				Relationships: []Relationship{},
			},
		},
		{
			name: "single relationship",
			input: &ServiceFile{
				Version: Version,
				Info: Info{
					Name: "single-service",
				},
				Relationships: []Relationship{
					{Action: "uses", Participant: "database", Technology: "postgres"},
				},
			},
			expected: &ServiceFile{
				Version: Version,
				Info: Info{
					Name: "single-service",
				},
				Relationships: []Relationship{
					{Action: "uses", Participant: "database", Technology: "postgres"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputCopy := &ServiceFile{
				Version:       tt.input.Version,
				Info:          tt.input.Info,
				Relationships: make([]Relationship, len(tt.input.Relationships)),
			}
			copy(inputCopy.Relationships, tt.input.Relationships)

			inputCopy.Sort()

			assert.Equal(t, tt.expected, inputCopy)
		})
	}
}
