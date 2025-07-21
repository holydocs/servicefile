package golang

import (
	"testing"

	"github.com/holydocs/servicefile/pkg/servicefile"
)

func TestParse(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		dir            string
		recursive      bool
		expectedResult []*servicefile.ServiceFile
		expectError    bool
	}{
		{
			name:      "parse default example service",
			dir:       "testdata/default",
			recursive: true,
			expectedResult: []*servicefile.ServiceFile{
				{
					Version: servicefile.Version,
					Info: servicefile.Info{
						Name:        "Example",
						Description: "Example service for exampling stuff.",
					},
					Relationships: []servicefile.Relationship{
						{
							Action:      servicefile.RelationshipActionReplies,
							Name:        "",
							Description: "Provides user management APIs to other services",
							Technology:  "grpc-server",
							Proto:       "grpc",
						},
						{
							Action:      servicefile.RelationshipActionRequests,
							Name:        "Firebase",
							Description: "Handles push notifications",
							Technology:  "firebase",
							Proto:       "http",
						},
						{
							Action:      servicefile.RelationshipActionUses,
							Name:        "PostgreSQL",
							Description: "Stores user data and authentication tokens",
							Technology:  "postgresql",
							Proto:       "tcp",
						},
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
			name:           "parse directory with no go files",
			dir:            "testdata",
			recursive:      false,
			expectError:    true,
			expectedResult: []*servicefile.ServiceFile{},
		},
		{
			name:      "parse explicit service relationships",
			dir:       "testdata/explicit",
			recursive: true,
			expectedResult: []*servicefile.ServiceFile{
				{
					Version: servicefile.Version,
					Info: servicefile.Info{
						Name:        "auth",
						Description: "Authentication service that handles user authentication and authorization",
					},
					Relationships: []servicefile.Relationship{
						{
							Action:      servicefile.RelationshipActionReplies,
							Name:        "user",
							Description: "Provides authentication responses to user service",
							Technology:  "jwt",
						},
						{
							Action:      servicefile.RelationshipActionReplies,
							Name:        "notification",
							Description: "Provides authentication status to notification service",
							Technology:  "grpc",
						},
					},
				},
				{
					Version: servicefile.Version,
					Info: servicefile.Info{
						Name:        "user",
						Description: "User management service that handles user profiles and data",
					},
					Relationships: []servicefile.Relationship{
						{
							Action:      servicefile.RelationshipActionRequests,
							Name:        "auth",
							Description: "Requests authentication from auth service",
							Technology:  "jwt",
						},
						{
							Action:      servicefile.RelationshipActionSends,
							Name:        "notification",
							Description: "Sends user events to notification service",
							Technology:  "grpc",
						},
					},
				},
				{
					Version: servicefile.Version,
					Info: servicefile.Info{
						Name:        "notification",
						Description: "Notification service that handles sending notifications to users",
					},
					Relationships: []servicefile.Relationship{
						{
							Action:      servicefile.RelationshipActionRequests,
							Name:        "auth",
							Description: "Requests authentication status from auth service",
							Technology:  "grpc",
						},
						{
							Action:      servicefile.RelationshipActionReceives,
							Name:        "user",
							Description: "Receives user events from user service",
							Technology:  "grpc",
						},
					},
				},
			},
			expectError: false,
		},
		{
			name:        "parse mixed service relationships",
			dir:         "testdata/mixed",
			recursive:   true,
			expectError: true,
		},
		{
			name:        "parse mixed service relationships with error message",
			dir:         "testdata/mixed",
			recursive:   true,
			expectError: true,
		},
	}

	compareServiceFileSlices := func(actual, expected []*servicefile.ServiceFile) bool {
		if len(actual) != len(expected) {
			return false
		}
		// Compare by name
		actualMap := make(map[string]*servicefile.ServiceFile)
		for _, sf := range actual {
			actualMap[sf.Info.Name] = sf
		}
		expectedMap := make(map[string]*servicefile.ServiceFile)
		for _, sf := range expected {
			expectedMap[sf.Info.Name] = sf
		}
		return compareServiceFiles(actualMap, expectedMap)
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

			if !compareServiceFileSlices(result, tt.expectedResult) {
				t.Errorf("Parse() = %+v, want %+v", result, tt.expectedResult)
				for _, v := range result {
					t.Logf("actual: Name=%s, Desc=%s, Rels=%+v", v.Info.Name, v.Info.Description, v.Relationships)
				}
				for _, v := range tt.expectedResult {
					t.Logf("expected: Name=%s, Desc=%s, Rels=%+v", v.Info.Name, v.Info.Description, v.Relationships)
				}
			}
		})
	}
}

func TestParseFile(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                  string
		filePath              string
		expectedServices      []service
		expectedRelationships []relationship
		expectError           bool
	}{
		{
			name:     "parse service file with comments",
			filePath: "testdata/default/service/example/example.go",
			expectedServices: []service{
				{
					name:        "Example",
					description: "Example service for exampling stuff.",
				},
			},
			expectedRelationships: []relationship{},
			expectError:           false,
		},
		{
			name:             "parse file with relationship comments",
			filePath:         "testdata/default/database/postgres/postgres.go",
			expectedServices: []service{},
			expectedRelationships: []relationship{
				{
					serviceName: "",
					action:      "uses",
					targetName:  "PostgreSQL",
					description: "Stores user data and authentication tokens",
					technology:  "postgresql",
					proto:       "tcp",
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

			if !compareServices(parser.services, tt.expectedServices) {
				t.Errorf("parseFile() services = %+v, want %+v", parser.services, tt.expectedServices)
			}

			if !compareRelationships(parser.relationships, tt.expectedRelationships) {
				t.Errorf("parseFile() relationships = %+v, want %+v", parser.relationships, tt.expectedRelationships)
			}
		})
	}
}

func TestParseCommentGroup(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                  string
		commentGroup          string
		expectedServices      []service
		expectedRelationships []relationship
	}{
		{
			name: "parse service name and description",
			commentGroup: `/*
service:name Example
description: Example service for exampling stuff.
*/`,
			expectedServices: []service{
				{
					name:        "Example",
					description: "Example service for exampling stuff.",
				},
			},
			expectedRelationships: []relationship{},
		},
		{
			name: "parse service name, description, and system",
			commentGroup: `/*
service:name UserService
description: Handles user authentication and profiles
system: e-commerce-platform
*/`,
			expectedServices: []service{
				{
					name:        "UserService",
					description: "Handles user authentication and profiles",
					system:      "e-commerce-platform",
				},
			},
			expectedRelationships: []relationship{},
		},
		{
			name: "parse relationship with all fields",
			commentGroup: `/*
service:uses PostgreSQL
description: Stores user data and authentication tokens
technology:postgresql
proto:tcp
*/`,
			expectedServices: []service{},
			expectedRelationships: []relationship{
				{
					serviceName: "",
					action:      "uses",
					targetName:  "PostgreSQL",
					description: "Stores user data and authentication tokens",
					technology:  "postgresql",
					proto:       "tcp",
				},
			},
		},
		{
			name:                  "parse empty comment group",
			commentGroup:          `/* */`,
			expectedServices:      []service{},
			expectedRelationships: []relationship{},
		},
		{
			name: "parse comments starting with //",
			commentGroup: `// service:name Example
// description: Example service for exampling stuff.`,
			expectedServices: []service{
				{
					name:        "Example",
					description: "Example service for exampling stuff.",
				},
			},
			expectedRelationships: []relationship{},
		},
		{
			name: "parse mixed comments: regular golang comments first, then service comments with /* */",
			commentGroup: `// User represents a user in the system
// This struct contains all user-related fields
/*
service:uses PostgreSQL
description: Stores user data and authentication tokens
technology:postgresql
proto:tcp
*/`,
			expectedServices: []service{},
			expectedRelationships: []relationship{
				{
					serviceName: "",
					action:      "uses",
					targetName:  "PostgreSQL",
					description: "Stores user data and authentication tokens",
					technology:  "postgresql",
					proto:       "tcp",
				},
			},
		},
		{
			name: "parse mixed comments: regular golang comments first, then service comments with //",
			commentGroup: `// User represents a user in the system
// This struct contains all user-related fields
// service:uses PostgreSQL
// description: Stores user data and authentication tokens
// technology:postgresql
// proto:tcp`,
			expectedServices: []service{},
			expectedRelationships: []relationship{
				{
					serviceName: "",
					action:      "uses",
					targetName:  "PostgreSQL",
					description: "Stores user data and authentication tokens",
					technology:  "postgresql",
					proto:       "tcp",
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
			expectedServices: []service{
				{
					name:        "Example",
					description: "Example service for exampling stuff.",
				},
			},
			expectedRelationships: []relationship{},
		},
		{
			name: "parse mixed comments: service comments with // first, then regular golang comments",
			commentGroup: `// service:name Example
// description: Example service for exampling stuff.
// User represents a user in the system
// This struct contains all user-related fields`,
			expectedServices: []service{
				{
					name:        "Example",
					description: "Example service for exampling stuff.",
				},
			},
			expectedRelationships: []relationship{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewCommentParser()
			parser.parseCommentGroup(tt.commentGroup)

			if !compareServices(parser.services, tt.expectedServices) {
				t.Errorf("parseCommentGroup() services = %+v, want %+v", parser.services, tt.expectedServices)
			}

			if !compareRelationships(parser.relationships, tt.expectedRelationships) {
				t.Errorf("parseCommentGroup() relationships = %+v, want %+v", parser.relationships, tt.expectedRelationships)
			}
		})
	}
}

// compareServiceFiles compares two service file maps for equality
func compareServiceFiles(actual, expected map[string]*servicefile.ServiceFile) bool {
	if len(actual) != len(expected) {
		return false
	}

	for serviceName, expectedService := range expected {
		actualService, exists := actual[serviceName]
		if !exists {
			return false
		}

		if actualService.Info.Name != expectedService.Info.Name {
			return false
		}

		if actualService.Info.Description != expectedService.Info.Description {
			return false
		}

		if actualService.Info.System != expectedService.Info.System {
			return false
		}

		if len(actualService.Relationships) != len(expectedService.Relationships) {
			return false
		}

		// Compare relationships (order doesn't matter for this test)
		for _, expectedRel := range expectedService.Relationships {
			found := false
			for _, actualRel := range actualService.Relationships {
				if actualRel.Action == expectedRel.Action &&
					actualRel.Name == expectedRel.Name &&
					actualRel.Description == expectedRel.Description &&
					actualRel.Technology == expectedRel.Technology &&
					actualRel.Proto == expectedRel.Proto {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
	}

	return true
}

// compareServices compares two service slices for equality
func compareServices(actual, expected []service) bool {
	if len(actual) != len(expected) {
		return false
	}

	// Compare services (order doesn't matter for this test)
	for _, expectedService := range expected {
		found := false
		for _, actualService := range actual {
			if actualService.name == expectedService.name &&
				actualService.description == expectedService.description &&
				actualService.system == expectedService.system {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	return true
}

// compareRelationships compares two relationship slices for equality
func compareRelationships(actual, expected []relationship) bool {
	if len(actual) != len(expected) {
		return false
	}

	// Compare relationships (order doesn't matter for this test)
	for _, expectedRel := range expected {
		found := false
		for _, actualRel := range actual {
			if actualRel.serviceName == expectedRel.serviceName &&
				actualRel.action == expectedRel.action &&
				actualRel.targetName == expectedRel.targetName &&
				actualRel.technology == expectedRel.technology &&
				actualRel.description == expectedRel.description &&
				actualRel.proto == expectedRel.proto {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	return true
}
