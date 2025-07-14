# ServiceFile

[![Run Tests](https://github.com/denchenko/servicefile/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/denchenko/servicefile/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/denchenko/servicefile)](https://goreportcard.com/report/github.com/denchenko/servicefile)
[![GoDoc](https://godoc.org/github.com/denchenko/servicefile?status.svg)](https://godoc.org/github.com/denchenko/servicefile)

A standardized approach to describing microservices and their relationships through structured specification file and automated parsing tools.

## Installation

```bash
# Install directly with Go
go install github.com/denchenko/servicefile/cmd/servicefile@latest
```

## Quick Start

### 1. Document Your Service

Add structured comments to your Go code to describe your service:

```go
/*
service:name UserService
description: Handles user authentication and profile management
*/
package main

// Service that uses PostgreSQL for data storage
/*
service:uses PostgreSQL
description: Stores user data and authentication tokens
technology:postgres
*/
type UserRepository struct {
    db *sql.DB
}

// Service that provides gRPC APIs
/*
service:replies
description: Provides user management APIs to other services
technology:grpc
*/
type UserServer struct {
    repo *UserRepository
}

// Service that makes HTTP requests to external service
/*
service:requests NotificationService
description: Sends user notifications via email and SMS
technology:http
*/
type NotificationClient struct {
    httpClient *http.Client
}
```

### 2. Parse Your Service

Use the CLI tool to parse your Go code and generate a service file:

```bash
# Parse current directory
servicefile parse

# Parse specific directory
servicefile parse --dir ./my-service

# Parse recursively (default)
servicefile parse --recursive

# Specify output file
servicefile parse --output my-service.yaml
```

### 3. Generated Output

The tool generates a `servicefile.yaml` with your service description:

```yaml
name: UserService
description: Handles user authentication and profile management
relationships:
  - action: uses
    name: PostgreSQL
    description: Stores user data and authentication tokens
    technology: postgres
  - action: replies
    description: Provides user management APIs to other services
    technology: grpc
  - action: requests
    name: NotificationService
    description: Sends user notifications via email and SMS
    technology: http
```

## ServiceFile Specification

### Service Metadata

- **`service:name`**: The name of your service
- **`description:`**: A description of what your service does

### Relationship Actions

ServiceFile supports several relationship types:

- **`service:uses`**: Service depends on another service/database
- **`service:requests`**: Service makes requests to another service
- **`service:replies`**: Service provides APIs for other services
- **`service:sends`**: Service sends messages/events
- **`service:receives`**: Service receives messages/events

### Relationship Properties

Each relationship can have:

- **`name`**: The name of the related service/resource
- **`description`**: Description of the relationship
- **`technology`**: Technology used (e.g., `grpc`, `http`, `postgres`, `redis`)

## Parsing Go Code

### Comment Format

ServiceFile parses Go comments to extract service information:

```go
/*
service:name MyService
description: This is my service description
*/

//service:uses DatabaseName
//description: Uses database for data storage
//technology:postgres
```

I recommend `/*` style, this way `go doc` will omit these declarations.

## Examples

See the `internal/parser/golang/testdata/` directory for complete examples of how to document services using ServiceFile comments.
