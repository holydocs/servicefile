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
    name: "database"
    description: "Uses PostgreSQL database"
    technology: "postgresql"
  - action: "sends"
    name: "notifications"
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
						Name:        "database",
						Description: "Uses PostgreSQL database",
						Technology:  "postgresql",
					},
					{
						Action:      "sends",
						Name:        "notifications",
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
    name: "database"
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
						Action:     "uses",
						Name:       "database",
						Technology: "postgresql",
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
    name: "database"
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
						Action:     "uses",
						Name:       "database",
						Technology: "postgresql",
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
    name: "database"
    description: "Uses PostgreSQL database"
    technology: "postgresql"
    proto: "tcp"
  - action: "requests"
    name: "auth-service"
    description: "Makes HTTP requests to authentication service"
    technology: "auth-service"
    proto: "http"
  - action: "replies"
    description: "Provides gRPC APIs"
    technology: "grpc-server"
    proto: "grpc"
  - action: "sends"
    name: "events"
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
						Name:        "database",
						Description: "Uses PostgreSQL database",
						Technology:  "postgresql",
						Proto:       "tcp",
					},
					{
						Action:      "requests",
						Name:        "auth-service",
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
						Name:        "events",
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
    name: "database"
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
						Action:     "uses",
						Name:       "database",
						Technology: "postgresql",
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
					{Action: "sends", Name: "b", Technology: "kafka", Proto: "tcp", Description: "second"},
					{Action: "uses", Name: "a", Technology: "postgres", Proto: "tcp", Description: "first"},
					{Action: "sends", Name: "a", Technology: "kafka", Proto: "udp", Description: "first"},
					{Action: "uses", Name: "a", Technology: "redis", Proto: "tcp", Description: "first"},
					{Action: "sends", Name: "a", Technology: "kafka", Proto: "tcp", Description: "first"},
				},
			},
			expected: &ServiceFile{
				Version: Version,
				Info: Info{
					Name: "test-service",
				},
				Relationships: []Relationship{
					{Action: "sends", Name: "a", Technology: "kafka", Proto: "tcp", Description: "first"},
					{Action: "sends", Name: "a", Technology: "kafka", Proto: "udp", Description: "first"},
					{Action: "sends", Name: "b", Technology: "kafka", Proto: "tcp", Description: "second"},
					{Action: "uses", Name: "a", Technology: "postgres", Proto: "tcp", Description: "first"},
					{Action: "uses", Name: "a", Technology: "redis", Proto: "tcp", Description: "first"},
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
					{Action: "uses", Name: "database", Technology: "postgres"},
				},
			},
			expected: &ServiceFile{
				Version: Version,
				Info: Info{
					Name: "single-service",
				},
				Relationships: []Relationship{
					{Action: "uses", Name: "database", Technology: "postgres"},
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
