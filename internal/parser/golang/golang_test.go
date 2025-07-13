package golang

import (
	"reflect"
	"testing"

	"github.com/denchenko/servicefile/pkg/servicefile"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name           string
		dir            string
		recursive      bool
		expectedResult *servicefile.ServiceFile
		expectError    bool
	}{
		{
			name:      "parse do example service",
			dir:       "testdata/do",
			recursive: true,
			expectedResult: &servicefile.ServiceFile{
				Name:        "Example",
				Description: "Example service for exampling stuff.",
				Relationships: []servicefile.Relationship{
					{
						Action:      servicefile.RelationshipActionReplies,
						Name:        "",
						Description: "Provides user management APIs to other services",
						Technology:  "grpc",
					},
					{
						Action:      servicefile.RelationshipActionRequests,
						Name:        "Firebase",
						Description: "Handles push notifications",
						Technology:  "http",
					},
					{
						Action:      servicefile.RelationshipActionUses,
						Name:        "PostgreSQL",
						Description: "Stores user data and authentication tokens",
						Technology:  "postgres",
					},
				},
			},
			expectError: false,
		},
		{
			name:        "parse non-existent directory",
			dir:         "testdata/nonexistent",
			recursive:   true,
			expectError: true,
		},
		{
			name:        "parse directory with no go files",
			dir:         "testdata",
			recursive:   false,
			expectError: false,
			expectedResult: &servicefile.ServiceFile{
				Relationships: []servicefile.Relationship{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewCommentParser()
			result, err := parser.Parse(tt.dir, tt.recursive)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if !reflect.DeepEqual(result, tt.expectedResult) {
				t.Errorf("Parse() = %+v, want %+v", result, tt.expectedResult)
			}
		})
	}
}

func TestParseFile(t *testing.T) {
	tests := []struct {
		name           string
		filePath       string
		expectedResult *servicefile.ServiceFile
		expectError    bool
	}{
		{
			name:     "parse service file with comments",
			filePath: "testdata/do/service/example/example.go",
			expectedResult: &servicefile.ServiceFile{
				Name:          "Example",
				Description:   "Example service for exampling stuff.",
				Relationships: []servicefile.Relationship{},
			},
			expectError: false,
		},
		{
			name:     "parse file with relationship comments",
			filePath: "testdata/do/database/postgres/postgres.go",
			expectedResult: &servicefile.ServiceFile{
				Relationships: []servicefile.Relationship{
					{
						Action:      servicefile.RelationshipActionUses,
						Name:        "PostgreSQL",
						Description: "Stores user data and authentication tokens",
						Technology:  "postgres",
					},
				},
			},
			expectError: false,
		},
		{
			name:        "parse non-existent file",
			filePath:    "testdata/nonexistent.go",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewCommentParser()
			err := parser.parseFile(tt.filePath)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// Reset the parser for each test to avoid state pollution
			parser = NewCommentParser()
			err = parser.parseFile(tt.filePath)
			if err != nil {
				t.Fatalf("Failed to parse file: %v", err)
			}

			if !reflect.DeepEqual(parser.serviceFile, tt.expectedResult) {
				t.Errorf("parseFile() = %+v, want %+v", parser.serviceFile, tt.expectedResult)
			}
		})
	}
}

func TestParseCommentGroup(t *testing.T) {
	tests := []struct {
		name           string
		commentGroup   string
		expectedResult *servicefile.ServiceFile
	}{
		{
			name: "parse service name and description",
			commentGroup: `/*
service:name Example
description: Example service for exampling stuff.
*/`,
			expectedResult: &servicefile.ServiceFile{
				Name:          "Example",
				Description:   "Example service for exampling stuff.",
				Relationships: []servicefile.Relationship{},
			},
		},
		{
			name: "parse relationship with all fields",
			commentGroup: `/*
service:uses PostgreSQL
description: Stores user data and authentication tokens
technology:postgres
*/`,
			expectedResult: &servicefile.ServiceFile{
				Relationships: []servicefile.Relationship{
					{
						Action:      servicefile.RelationshipActionUses,
						Name:        "PostgreSQL",
						Description: "Stores user data and authentication tokens",
						Technology:  "postgres",
					},
				},
			},
		},
		{
			name: "parse multiple relationships",
			commentGroup: `/*
service:uses PostgreSQL
description: Stores user data and authentication tokens
technology:postgres

service:requests Firebase
description: Handles push notifications
technology:http
*/`,
			expectedResult: &servicefile.ServiceFile{
				Relationships: []servicefile.Relationship{
					{
						Action:      servicefile.RelationshipActionUses,
						Name:        "PostgreSQL",
						Description: "Stores user data and authentication tokens",
						Technology:  "postgres",
					},
					{
						Action:      servicefile.RelationshipActionRequests,
						Name:        "Firebase",
						Description: "Handles push notifications",
						Technology:  "http",
					},
				},
			},
		},
		{
			name: "parse empty comment group",
			commentGroup: `/*
*/`,
			expectedResult: &servicefile.ServiceFile{
				Relationships: []servicefile.Relationship{},
			},
		},
		{
			name: "parse comments starting with //",
			commentGroup: `// service:name Example
// description: Example service for exampling stuff.
// service:uses PostgreSQL
// description: Stores user data and authentication tokens
// technology:postgres`,
			expectedResult: &servicefile.ServiceFile{
				Name:        "Example",
				Description: "Example service for exampling stuff.",
				Relationships: []servicefile.Relationship{
					{
						Action:      servicefile.RelationshipActionUses,
						Name:        "PostgreSQL",
						Description: "Stores user data and authentication tokens",
						Technology:  "postgres",
					},
				},
			},
		},
		{
			name: "parse mixed comments: regular golang comments first, then service comments with /* */",
			commentGroup: `// User represents a user in the system
// This struct contains all user-related fields
/*
service:uses PostgreSQL
description: Stores user data and authentication tokens
technology:postgres
*/`,
			expectedResult: &servicefile.ServiceFile{
				Relationships: []servicefile.Relationship{
					{
						Action:      servicefile.RelationshipActionUses,
						Name:        "PostgreSQL",
						Description: "Stores user data and authentication tokens",
						Technology:  "postgres",
					},
				},
			},
		},
		{
			name: "parse mixed comments: regular golang comments first, then service comments with //",
			commentGroup: `// User represents a user in the system
// This struct contains all user-related fields
// service:name Example
// description: Example service for exampling stuff.
// service:uses PostgreSQL
// description: Stores user data and authentication tokens
// technology:postgres`,
			expectedResult: &servicefile.ServiceFile{
				Name:        "Example",
				Description: "Example service for exampling stuff.",
				Relationships: []servicefile.Relationship{
					{
						Action:      servicefile.RelationshipActionUses,
						Name:        "PostgreSQL",
						Description: "Stores user data and authentication tokens",
						Technology:  "postgres",
					},
				},
			},
		},
		{
			name: "parse mixed comments: service comments with /* */ first, then regular golang comments",
			commentGroup: `/*
service:name Example
description: Example service for exampling stuff.
*/
// User represents a user in the system
// This struct contains all user-related fields`,
			expectedResult: &servicefile.ServiceFile{
				Name:          "Example",
				Description:   "Example service for exampling stuff.",
				Relationships: []servicefile.Relationship{},
			},
		},
		{
			name: "parse mixed comments: service comments with // first, then regular golang comments",
			commentGroup: `// service:name Example
// description: Example service for exampling stuff.
// service:uses PostgreSQL
// description: Stores user data and authentication tokens
// technology:postgres
// User represents a user in the system
// This struct contains all user-related fields`,
			expectedResult: &servicefile.ServiceFile{
				Name:        "Example",
				Description: "Example service for exampling stuff.",
				Relationships: []servicefile.Relationship{
					{
						Action:      servicefile.RelationshipActionUses,
						Name:        "PostgreSQL",
						Description: "Stores user data and authentication tokens",
						Technology:  "postgres",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewCommentParser()
			parser.parseCommentGroup(tt.commentGroup)

			if !reflect.DeepEqual(parser.serviceFile, tt.expectedResult) {
				t.Errorf("parseCommentGroup() = %+v, want %+v", parser.serviceFile, tt.expectedResult)
			}
		})
	}
}
