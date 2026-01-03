package golang

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/holydocs/servicefile/pkg/servicefile"
	"golang.org/x/mod/modfile"
)

type inferredDependency struct {
	modulePrefix string
	rel          servicefile.Relationship
}

var inferredDependencies = []inferredDependency{
	// PostgreSQL
	{
		modulePrefix: "github.com/jackc/pgx",
		rel: servicefile.Relationship{
			Action:      servicefile.RelationshipActionUses,
			Participant: "PostgreSQL",
			Technology:  "postgresql",
			Proto:       "tcp",
		},
	},
	{
		modulePrefix: "github.com/lib/pq",
		rel: servicefile.Relationship{
			Action:      servicefile.RelationshipActionUses,
			Participant: "PostgreSQL",
			Technology:  "postgresql",
			Proto:       "tcp",
		},
	},
	{
		modulePrefix: "gorm.io/driver/postgres",
		rel: servicefile.Relationship{
			Action:      servicefile.RelationshipActionUses,
			Participant: "PostgreSQL",
			Technology:  "postgresql",
			Proto:       "tcp",
		},
	},

	// MySQL
	{
		modulePrefix: "github.com/go-sql-driver/mysql",
		rel: servicefile.Relationship{
			Action:      servicefile.RelationshipActionUses,
			Participant: "MySQL",
			Technology:  "mysql",
			Proto:       "tcp",
		},
	},
	{
		modulePrefix: "gorm.io/driver/mysql",
		rel: servicefile.Relationship{
			Action:      servicefile.RelationshipActionUses,
			Participant: "MySQL",
			Technology:  "mysql",
			Proto:       "tcp",
		},
	},

	// SQLite
	{
		modulePrefix: "github.com/mattn/go-sqlite3",
		rel: servicefile.Relationship{
			Action:      servicefile.RelationshipActionUses,
			Participant: "SQLite",
			Technology:  "sqlite",
			Proto:       "file",
		},
	},
	{
		modulePrefix: "modernc.org/sqlite",
		rel: servicefile.Relationship{
			Action:      servicefile.RelationshipActionUses,
			Participant: "SQLite",
			Technology:  "sqlite",
			Proto:       "file",
		},
	},
	{
		modulePrefix: "gorm.io/driver/sqlite",
		rel: servicefile.Relationship{
			Action:      servicefile.RelationshipActionUses,
			Participant: "SQLite",
			Technology:  "sqlite",
			Proto:       "file",
		},
	},

	// Redis
	{
		modulePrefix: "github.com/redis/go-redis",
		rel: servicefile.Relationship{
			Action:      servicefile.RelationshipActionUses,
			Participant: "Redis",
			Technology:  "redis",
			Proto:       "tcp",
		},
	},

	// MongoDB
	{
		modulePrefix: "go.mongodb.org/mongo-driver",
		rel: servicefile.Relationship{
			Action:      servicefile.RelationshipActionUses,
			Participant: "MongoDB",
			Technology:  "mongodb",
			Proto:       "tcp",
		},
	},

	// Cassandra / ScyllaDB
	{
		modulePrefix: "github.com/gocql/gocql",
		rel: servicefile.Relationship{
			Action:      servicefile.RelationshipActionUses,
			Participant: "Cassandra",
			Technology:  "cassandra",
			Proto:       "tcp",
		},
	},
	{
		// Treat gocqlx as ScyllaDB by default (requested).
		modulePrefix: "github.com/scylladb/gocqlx",
		rel: servicefile.Relationship{
			Action:      servicefile.RelationshipActionUses,
			Participant: "ScyllaDB",
			Technology:  "scylladb",
			Proto:       "tcp",
		},
	},

	// Kafka
	{
		modulePrefix: "github.com/segmentio/kafka-go",
		rel: servicefile.Relationship{
			Action:      servicefile.RelationshipActionUses,
			Participant: "Kafka",
			Technology:  "kafka",
			Proto:       "tcp",
		},
	},
	{
		modulePrefix: "github.com/IBM/sarama",
		rel: servicefile.Relationship{
			Action:      servicefile.RelationshipActionUses,
			Participant: "Kafka",
			Technology:  "kafka",
			Proto:       "tcp",
		},
	},

	// RabbitMQ
	{
		modulePrefix: "github.com/rabbitmq/amqp091-go",
		rel: servicefile.Relationship{
			Action:      servicefile.RelationshipActionUses,
			Participant: "RabbitMQ",
			Technology:  "rabbitmq",
			Proto:       "amqp",
		},
	},

	// NATS
	{
		modulePrefix: "github.com/nats-io/nats.go",
		rel: servicefile.Relationship{
			Action:      servicefile.RelationshipActionUses,
			Participant: "NATS",
			Technology:  "nats",
			Proto:       "tcp",
		},
	},

	// Elasticsearch / OpenSearch
	{
		modulePrefix: "github.com/elastic/go-elasticsearch",
		rel: servicefile.Relationship{
			Action:      servicefile.RelationshipActionUses,
			Participant: "Elasticsearch",
			Technology:  "elasticsearch",
			Proto:       "http",
		},
	},
	{
		modulePrefix: "github.com/opensearch-project/opensearch-go",
		rel: servicefile.Relationship{
			Action:      servicefile.RelationshipActionUses,
			Participant: "OpenSearch",
			Technology:  "opensearch",
			Proto:       "http",
		},
	},

	// ClickHouse
	{
		modulePrefix: "github.com/ClickHouse/clickhouse-go/v2",
		rel: servicefile.Relationship{
			Action:      servicefile.RelationshipActionUses,
			Participant: "ClickHouse",
			Technology:  "clickhouse",
			Proto:       "tcp",
		},
	},
	{
		modulePrefix: "github.com/ClickHouse/clickhouse-go",
		rel: servicefile.Relationship{
			Action:      servicefile.RelationshipActionUses,
			Participant: "ClickHouse",
			Technology:  "clickhouse",
			Proto:       "tcp",
		},
	},

	// Memcached
	{
		modulePrefix: "github.com/bradfitz/gomemcache",
		rel: servicefile.Relationship{
			Action:      servicefile.RelationshipActionUses,
			Participant: "Memcached",
			Technology:  "memcached",
			Proto:       "tcp",
		},
	},

	// etcd
	{
		modulePrefix: "go.etcd.io/etcd/client/v3",
		rel: servicefile.Relationship{
			Action:      servicefile.RelationshipActionUses,
			Participant: "etcd",
			Technology:  "etcd",
			Proto:       "grpc",
		},
	},

	// BadgerDB (embedded KV)
	{
		modulePrefix: "github.com/dgraph-io/badger",
		rel: servicefile.Relationship{
			Action:      servicefile.RelationshipActionUses,
			Participant: "BadgerDB",
			Technology:  "badgerdb",
			Proto:       "file",
		},
	},

	// InfluxDB
	{
		modulePrefix: "github.com/influxdata/influxdb-client-go",
		rel: servicefile.Relationship{
			Action:      servicefile.RelationshipActionUses,
			Participant: "InfluxDB",
			Technology:  "influxdb",
			Proto:       "http",
		},
	},

	// MinIO / S3-compatible storage
	{
		modulePrefix: "github.com/minio/minio-go",
		rel: servicefile.Relationship{
			Action:      servicefile.RelationshipActionUses,
			Participant: "MinIO",
			Technology:  "minio",
			Proto:       "http",
		},
	},

	// Pulsar
	{
		modulePrefix: "github.com/apache/pulsar-client-go",
		rel: servicefile.Relationship{
			Action:      servicefile.RelationshipActionUses,
			Participant: "Pulsar",
			Technology:  "pulsar",
			Proto:       "tcp",
		},
	},

	// ORMs / SQL helpers (generic SQL database usage)
	{
		modulePrefix: "gorm.io/gorm",
		rel: servicefile.Relationship{
			Action:      servicefile.RelationshipActionUses,
			Participant: "SQL Database",
			Technology:  "sql",
			Proto:       "tcp",
		},
	},
	{
		modulePrefix: "github.com/jmoiron/sqlx",
		rel: servicefile.Relationship{
			Action:      servicefile.RelationshipActionUses,
			Participant: "SQL Database",
			Technology:  "sql",
			Proto:       "tcp",
		},
	},
	{
		modulePrefix: "github.com/uptrace/bun",
		rel: servicefile.Relationship{
			Action:      servicefile.RelationshipActionUses,
			Participant: "SQL Database",
			Technology:  "sql",
			Proto:       "tcp",
		},
	},

	// AWS SDK v2 (selected common managed services)
	{
		modulePrefix: "github.com/aws/aws-sdk-go-v2/service/s3",
		rel: servicefile.Relationship{
			Action:      servicefile.RelationshipActionUses,
			Participant: "Amazon S3",
			Technology:  "aws-s3",
			Proto:       "http",
		},
	},
	{
		modulePrefix: "github.com/aws/aws-sdk-go-v2/service/sqs",
		rel: servicefile.Relationship{
			Action:      servicefile.RelationshipActionUses,
			Participant: "Amazon SQS",
			Technology:  "aws-sqs",
			Proto:       "http",
		},
	},
	{
		modulePrefix: "github.com/aws/aws-sdk-go-v2/service/sns",
		rel: servicefile.Relationship{
			Action:      servicefile.RelationshipActionUses,
			Participant: "Amazon SNS",
			Technology:  "aws-sns",
			Proto:       "http",
		},
	},
	{
		modulePrefix: "github.com/aws/aws-sdk-go-v2/service/dynamodb",
		rel: servicefile.Relationship{
			Action:      servicefile.RelationshipActionUses,
			Participant: "Amazon DynamoDB",
			Technology:  "aws-dynamodb",
			Proto:       "http",
		},
	},
	{
		modulePrefix: "github.com/aws/aws-sdk-go-v2/service/kinesis",
		rel: servicefile.Relationship{
			Action:      servicefile.RelationshipActionUses,
			Participant: "Amazon Kinesis",
			Technology:  "aws-kinesis",
			Proto:       "http",
		},
	},
	{
		modulePrefix: "github.com/aws/aws-sdk-go-v2/service/secretsmanager",
		rel: servicefile.Relationship{
			Action:      servicefile.RelationshipActionUses,
			Participant: "AWS Secrets Manager",
			Technology:  "aws-secrets-manager",
			Proto:       "http",
		},
	},

	// Google Cloud
	{
		modulePrefix: "cloud.google.com/go/storage",
		rel: servicefile.Relationship{
			Action:      servicefile.RelationshipActionUses,
			Participant: "Google Cloud Storage",
			Technology:  "gcp-storage",
			Proto:       "http",
		},
	},
	{
		modulePrefix: "cloud.google.com/go/pubsub",
		rel: servicefile.Relationship{
			Action:      servicefile.RelationshipActionUses,
			Participant: "Google Cloud Pub/Sub",
			Technology:  "gcp-pubsub",
			Proto:       "http",
		},
	},
	{
		modulePrefix: "cloud.google.com/go/firestore",
		rel: servicefile.Relationship{
			Action:      servicefile.RelationshipActionUses,
			Participant: "Google Cloud Firestore",
			Technology:  "gcp-firestore",
			Proto:       "http",
		},
	},
	{
		modulePrefix: "cloud.google.com/go/secretmanager",
		rel: servicefile.Relationship{
			Action:      servicefile.RelationshipActionUses,
			Participant: "Google Secret Manager",
			Technology:  "gcp-secret-manager",
			Proto:       "http",
		},
	},

	// Azure (selected common managed services)
	{
		modulePrefix: "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob",
		rel: servicefile.Relationship{
			Action:      servicefile.RelationshipActionUses,
			Participant: "Azure Blob Storage",
			Technology:  "azure-blob-storage",
			Proto:       "http",
		},
	},
	{
		modulePrefix: "github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus",
		rel: servicefile.Relationship{
			Action:      servicefile.RelationshipActionUses,
			Participant: "Azure Service Bus",
			Technology:  "azure-service-bus",
			Proto:       "http",
		},
	},
	{
		modulePrefix: "github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets",
		rel: servicefile.Relationship{
			Action:      servicefile.RelationshipActionUses,
			Participant: "Azure Key Vault",
			Technology:  "azure-key-vault",
			Proto:       "http",
		},
	},
}

func (cp *CommentParser) enrichWithGoModDependencies(rootDir string, serviceFiles []*servicefile.ServiceFile) error {
	serviceToSource := cp.serviceDefinedInByName()

	rootAbs, err := filepath.Abs(rootDir)
	if err != nil {
		return fmt.Errorf("resolve root dir: %w", err)
	}

	// Cache parsed go.mod results by path to avoid reparsing for multiple services in the same module.
	modCache := make(map[string][]string)

	for _, sf := range serviceFiles {
		src, ok := serviceToSource[sf.Info.Name]
		if !ok || strings.TrimSpace(src) == "" {
			continue
		}

		goModPath, found, err := findNearestGoMod(src, rootAbs)
		if err != nil {
			return err
		}
		if !found {
			continue
		}

		requiredModules, ok := modCache[goModPath]
		if !ok {
			mods, err := parseDirectRequires(goModPath)
			if err != nil {
				return err
			}
			requiredModules = mods
			modCache[goModPath] = mods
		}

		inferred := inferRelationshipsFromModules(requiredModules)
		if len(inferred) == 0 {
			continue
		}

		mergeRelationships(sf, inferred)
		sf.Sort()
	}

	return nil
}

func (cp *CommentParser) serviceDefinedInByName() map[string]string {
	out := make(map[string]string)

	for _, s := range cp.services {
		if s.name == "" || strings.TrimSpace(s.definedIn) == "" {
			continue
		}
		if _, exists := out[s.name]; !exists {
			out[s.name] = s.definedIn
		}
	}

	// Fallback: if relationships explicitly name a service, use the relationship file location to map it.
	for _, r := range cp.relationships {
		if r.serviceName == "" || strings.TrimSpace(r.definedIn) == "" {
			continue
		}
		if _, exists := out[r.serviceName]; !exists {
			out[r.serviceName] = r.definedIn
		}
	}

	return out
}

func findNearestGoMod(sourceFile string, rootAbs string) (goModPath string, found bool, err error) {
	srcAbs, err := filepath.Abs(sourceFile)
	if err != nil {
		return "", false, fmt.Errorf("resolve source path %q: %w", sourceFile, err)
	}

	dir := filepath.Dir(srcAbs)
	for {
		// Ensure we don't walk above rootAbs.
		rel, relErr := filepath.Rel(rootAbs, dir)
		if relErr != nil {
			return "", false, fmt.Errorf("rel %q to %q: %w", rootAbs, dir, relErr)
		}
		if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
			return "", false, nil
		}

		candidate := filepath.Join(dir, "go.mod")
		if _, statErr := os.Stat(candidate); statErr == nil {
			return candidate, true, nil
		} else if !os.IsNotExist(statErr) {
			return "", false, fmt.Errorf("stat %s: %w", candidate, statErr)
		}

		if dir == rootAbs {
			return "", false, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", false, nil
		}
		dir = parent
	}
}

func parseDirectRequires(goModPath string) ([]string, error) {
	data, err := os.ReadFile(goModPath)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", goModPath, err)
	}

	f, err := modfile.Parse(goModPath, data, nil)
	if err != nil {
		return nil, fmt.Errorf("parse %s: %w", goModPath, err)
	}

	mods := make([]string, 0, len(f.Require))
	for _, req := range f.Require {
		if req == nil || req.Mod.Path == "" {
			continue
		}
		if req.Indirect {
			continue
		}
		mods = append(mods, req.Mod.Path)
	}

	return mods, nil
}

func inferRelationshipsFromModules(modules []string) []servicefile.Relationship {
	var out []servicefile.Relationship

	for _, m := range modules {
		for _, dep := range inferredDependencies {
			if strings.HasPrefix(m, dep.modulePrefix) {
				out = append(out, dep.rel)
				break
			}
		}
	}

	return out
}

func mergeRelationships(sf *servicefile.ServiceFile, inferred []servicefile.Relationship) {
	// If the user already defined a relationship with the same (action, participant, technology),
	// do not add an auto-detected relationship even if proto/description differ.
	seenLoose := make(map[string]struct{}, len(sf.Relationships))
	for _, r := range sf.Relationships {
		seenLoose[relationshipKeyLoose(r)] = struct{}{}
	}

	hasSpecificSQL := serviceHasSpecificSQL(sf)

	for _, r := range inferred {
		// If we already have a specific SQL database (e.g., PostgreSQL/MySQL/SQLite),
		// avoid also adding the generic ORM-based SQL relationship.
		if hasSpecificSQL && isGenericSQLRelationship(r) {
			continue
		}

		kLoose := relationshipKeyLoose(r)
		if _, exists := seenLoose[kLoose]; exists {
			continue
		}

		sf.Relationships = append(sf.Relationships, r)
		seenLoose[kLoose] = struct{}{}
	}
}

func isGenericSQLRelationship(r servicefile.Relationship) bool {
	return r.Action == servicefile.RelationshipActionUses &&
		r.Technology == "sql" &&
		r.Participant == "SQL Database"
}

func serviceHasSpecificSQL(sf *servicefile.ServiceFile) bool {
	for _, r := range sf.Relationships {
		if r.Action != servicefile.RelationshipActionUses {
			continue
		}
		switch r.Technology {
		case "postgresql", "mysql", "sqlite", "mssql":
			return true
		}
	}
	return false
}

func relationshipKeyLoose(r servicefile.Relationship) string {
	return strings.Join([]string{
		normalizeAction(string(r.Action)),
		normalizeParticipant(r.Participant),
		normalizeTechnology(r.Technology),
	}, "|")
}

func normalizeAction(action string) string {
	return strings.ToLower(strings.TrimSpace(action))
}

func normalizeParticipant(p string) string {
	return strings.ToLower(strings.TrimSpace(p))
}

func normalizeTechnology(t string) string {
	raw := strings.ToLower(strings.TrimSpace(t))
	switch raw {
	case "postgres", "postgresql", "postgre", "postgresqldb", "postgresql-db", "postgres-db", "postgresql database", "postgres database", "postgresqlclient", "postgres client", "postgresql client":
		return "postgresql"
	case "mysql", "mariadb":
		return "mysql"
	case "sqlite", "sqlite3":
		return "sqlite"
	case "mssql", "sqlserver", "sql server":
		return "mssql"
	default:
		return raw
	}
}
