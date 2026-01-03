package golang

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/holydocs/servicefile/pkg/servicefile"
)

type Parser struct {
	services      []service
	relationships []relationship
	firstGoFile   string
}

func NewParser() *Parser {
	return &Parser{
		services:      make([]service, 0),
		relationships: make([]relationship, 0),
		firstGoFile:   "",
	}
}

func (cp *Parser) Parse(dir string, recursive bool, detectRepository bool, analyzeGoMod bool) ([]*servicefile.ServiceFile, error) {
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("failed to walk the path: %w", err)
		}

		if info.IsDir() && !recursive && path != dir {
			return filepath.SkipDir
		}

		if info.IsDir() {
			return nil
		}

		if !strings.HasSuffix(path, ".go") {
			return nil
		}

		if cp.firstGoFile == "" {
			cp.firstGoFile = path
		}

		if err := cp.parseFile(path); err != nil {
			return fmt.Errorf("failed to parse %s: %w", path, err)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error walking the path: %w", err)
	}

	if strings.TrimSpace(cp.firstGoFile) == "" {
		return nil, fmt.Errorf("no go files found")
	}

	serviceFiles, err := cp.buildServiceFiles(dir, detectRepository)
	if err != nil {
		return nil, err
	}

	if detectRepository && isEmptyRepository(serviceFiles) {
		if err := cp.fillRepository(dir, serviceFiles); err != nil {
			return nil, fmt.Errorf("error detecting repositories: %w", err)
		}
	}

	if analyzeGoMod {
		if err := cp.enrichWithGoModDependencies(dir, serviceFiles); err != nil {
			return nil, fmt.Errorf("error analyzing go.mod dependencies: %w", err)
		}
	}

	return serviceFiles, nil
}

type service struct {
	name        string
	description string
	system      string
	owner       string
	repository  string
	tags        []string
	definedIn   string
}

func (s service) String() string {
	return fmt.Sprintf("name: %s, description: %s", s.name, s.description)
}

type relationship struct {
	serviceName string
	action      string
	targetName  string
	technology  string
	description string
	proto       string
	tags        []string
	external    bool
	person      bool
	definedIn   string
}

func (r relationship) String() string {
	return fmt.Sprintf("service_name: %s, action: %s, target_name: %s",
		r.serviceName,
		r.action,
		r.targetName,
	)
}

func (cp *Parser) parseFile(path string) error {
	fset := token.NewFileSet()

	f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("failed to parse %s: %w", path, err)
	}

	for _, cg := range f.Comments {
		var commentText strings.Builder
		for _, c := range cg.List {
			commentText.WriteString(c.Text)
			commentText.WriteString("\n")
		}
		cp.parseCommentGroup(commentText.String(), path)
	}

	ast.Inspect(f, func(n ast.Node) bool {
		x, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}

		if x.Doc == nil {
			return true
		}

		var commentText strings.Builder
		for _, c := range x.Doc.List {
			commentText.WriteString(c.Text)
			commentText.WriteString("\n")
		}
		cp.parseCommentGroup(commentText.String(), path)

		return true
	})

	return nil
}

func (cp *Parser) parseCommentGroup(commentGroup string, definedIn string) {
	if !strings.Contains(commentGroup, "service:") {
		return
	}

	lines := strings.Split(commentGroup, "\n")

	switch {
	case strings.Contains(commentGroup, "service:name"):
		cp.parseServiceDefinition(lines, definedIn)
	default:
		cp.parseRelationshipDefinition(lines, definedIn)
	}
}

func (cp *Parser) parseServiceDefinition(lines []string, definedIn string) {
	var s service

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		comment := cp.extractCommentText(line)
		if comment == "" {
			continue
		}

		if strings.HasPrefix(comment, "service:name") {
			parts := strings.SplitN(comment, " ", 2)
			if len(parts) == 2 {
				s.name = strings.TrimSpace(parts[1])
			}
			continue
		}

		if strings.HasPrefix(comment, "description:") {
			parts := strings.SplitN(comment, ":", 2)
			if len(parts) == 2 {
				s.description = strings.TrimSpace(parts[1])
			}
			continue
		}

		if strings.HasPrefix(comment, "system:") {
			parts := strings.SplitN(comment, ":", 2)
			if len(parts) == 2 {
				s.system = strings.TrimSpace(parts[1])
			}
			continue
		}

		if strings.HasPrefix(comment, "owner:") {
			parts := strings.SplitN(comment, ":", 2)
			if len(parts) == 2 {
				s.owner = strings.TrimSpace(parts[1])
			}
			continue
		}

		if strings.HasPrefix(comment, "repository:") {
			parts := strings.SplitN(comment, ":", 2)
			if len(parts) == 2 {
				s.repository = strings.TrimSpace(parts[1])
			}
			continue
		}

		if strings.HasPrefix(comment, "tags:") {
			parts := strings.SplitN(comment, ":", 2)
			if len(parts) == 2 {
				tagsStr := strings.TrimSpace(parts[1])
				if tagsStr != "" {
					// Split tags by comma and trim whitespace
					tags := strings.Split(tagsStr, ",")
					for i, tag := range tags {
						tags[i] = strings.TrimSpace(tag)
					}
					s.tags = tags
				}
			}
			continue
		}
	}

	if s.name != "" {
		s.definedIn = definedIn
		cp.services = append(cp.services, s)
	}
}

func (cp *Parser) parseRelationshipDefinition(lines []string, definedIn string) {
	var r relationship

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		comment := cp.extractCommentText(line)
		if comment == "" {
			continue
		}

		switch {
		case strings.HasPrefix(comment, "service:"):
			r.serviceName, r.action, r.targetName = cp.extractRelationshipInfo(comment)
			continue
		case strings.HasPrefix(comment, "technology:"):
			parts := strings.SplitN(comment, ":", 2)
			if len(parts) == 2 {
				r.technology = strings.TrimSpace(parts[1])
			}
			continue
		case strings.HasPrefix(comment, "description:"):
			parts := strings.SplitN(comment, ":", 2)
			if len(parts) == 2 {
				r.description = strings.TrimSpace(parts[1])
			}
			continue
		case strings.HasPrefix(comment, "proto:"):
			parts := strings.SplitN(comment, ":", 2)
			if len(parts) == 2 {
				r.proto = strings.TrimSpace(parts[1])
			}
			continue
		case strings.HasPrefix(comment, "tags:"):
			parts := strings.SplitN(comment, ":", 2)
			if len(parts) == 2 {
				tagsStr := strings.TrimSpace(parts[1])
				if tagsStr != "" {
					tags := strings.Split(tagsStr, ",")
					for i, tag := range tags {
						tags[i] = strings.TrimSpace(tag)
					}
					r.tags = tags
				}
			}
			continue
		case strings.HasPrefix(comment, "external:"):
			parts := strings.SplitN(comment, ":", 2)
			if len(parts) == 2 {
				externalStr := strings.TrimSpace(parts[1])
				r.external = externalStr == "true" || externalStr == "yes" || externalStr == "1"
			}
			continue
		case strings.HasPrefix(comment, "person:"):
			parts := strings.SplitN(comment, ":", 2)
			if len(parts) == 2 {
				personStr := strings.TrimSpace(parts[1])
				r.person = personStr == "true" || personStr == "yes" || personStr == "1"
			}
			continue
		}
	}

	if r.action != "" {
		r.definedIn = definedIn
		cp.relationships = append(cp.relationships, r)
	}
}

func (cp *Parser) extractCommentText(line string) string {
	comment := strings.TrimSpace(line)
	comment = strings.TrimPrefix(comment, "//")
	comment = strings.TrimPrefix(comment, "/*")
	comment = strings.TrimSuffix(comment, "*/")
	return strings.TrimSpace(comment)
}

// extractRelationshipInfo extracts the service name, action, and target name from a comment.
// Format: service:{service_name}:{action} [target_service] or service:{action} [target_service]
// Example: service:database:uses PostgreSQL
// Example: service:uses PostgreSQL
func (cp *Parser) extractRelationshipInfo(comment string) (serviceName, action, targetName string) {
	parts := strings.SplitN(comment, " ", 2)
	serviceActionPart := parts[0]

	serviceActionParts := strings.Split(serviceActionPart, ":")
	if len(serviceActionParts) >= 3 {
		// Format: service:{service_name}:{action}
		serviceName = serviceActionParts[1]
		action = serviceActionParts[2]
	} else if len(serviceActionParts) == 2 {
		// Format: service:{action}
		action = serviceActionParts[1]
	}

	// Extract target name if present
	if len(parts) > 1 {
		targetName = strings.TrimSpace(parts[1])
	}

	return serviceName, action, targetName
}

func (cp *Parser) buildServiceFiles(rootDir string, detectRepository bool) ([]*servicefile.ServiceFile, error) {
	if err := cp.validateNoMixedUsage(); err != nil {
		return nil, err
	}

	serviceFiles := make(map[string]*servicefile.ServiceFile)

	for _, s := range cp.services {
		serviceFiles[s.name] = &servicefile.ServiceFile{
			Version: servicefile.Version,
			Info: servicefile.Info{
				Name:        s.name,
				Description: s.description,
				System:      s.system,
				Owner:       s.owner,
				Repository:  s.repository,
				Tags:        s.tags,
			},
			Relationships: []servicefile.Relationship{},
		}
	}

	for _, r := range cp.relationships {
		serviceName, err := cp.determineServiceName(r, serviceFiles)
		if err != nil {
			return nil, fmt.Errorf("failed to determine service name: %w", err)
		}

		if _, exists := serviceFiles[serviceName]; !exists {
			serviceFiles[serviceName] = &servicefile.ServiceFile{
				Version: servicefile.Version,
				Info: servicefile.Info{
					Name: serviceName,
				},
				Relationships: []servicefile.Relationship{},
			}
		}

		relationship := servicefile.Relationship{
			Action:      servicefile.RelationshipAction(r.action),
			Participant: r.targetName,
		}

		if r.technology != "" {
			relationship.Technology = r.technology
		}

		if r.description != "" {
			relationship.Description = r.description
		}

		if r.proto != "" {
			relationship.Proto = r.proto
		}

		if len(r.tags) > 0 {
			relationship.Tags = r.tags
		}

		if r.external {
			relationship.External = r.external
		}

		if r.person {
			relationship.Person = r.person
		}

		serviceFiles[serviceName].Relationships = append(serviceFiles[serviceName].Relationships, relationship)
	}

	if len(serviceFiles) == 0 {
		sf, err := cp.buildDefaultServiceFile(rootDir, detectRepository)
		if err != nil {
			return nil, err
		}
		return []*servicefile.ServiceFile{sf}, nil
	}

	result := make([]*servicefile.ServiceFile, 0, len(serviceFiles))
	for _, sf := range serviceFiles {
		sf.Sort()
		result = append(result, sf)
	}

	return result, nil
}

func (cp *Parser) buildDefaultServiceFile(rootDir string, detectRepository bool) (*servicefile.ServiceFile, error) {
	// 1) Prefer git remote basename for service name when available
	// 2) Otherwise fallback to current directory basename
	name := ""
	repoURL := ""

	url, err := detectGitRepository(rootDir)
	if err == nil && strings.TrimSpace(url) != "" {
		repoURL = url
		if base := basenameFromRepoURL(url); base != "" {
			name = base
		}
	}

	if name == "" {
		abs, absErr := filepath.Abs(rootDir)
		if absErr != nil {
			// Best-effort fallback
			name = filepath.Base(rootDir)
		} else {
			name = filepath.Base(abs)
		}
	}

	sf := &servicefile.ServiceFile{
		Version: servicefile.Version,
		Info: servicefile.Info{
			Name: name,
		},
		Relationships: []servicefile.Relationship{},
	}

	if detectRepository && repoURL != "" {
		sf.Info.Repository = repoURL
	}

	// If we parsed any relationship blocks without a service definition, attach them to the default service.
	for _, r := range cp.relationships {
		if r.serviceName != "" && r.serviceName != name {
			// Explicitly named relationships belong to another service; keep behavior strict.
			continue
		}

		rel := servicefile.Relationship{
			Action:      servicefile.RelationshipAction(r.action),
			Participant: r.targetName,
		}
		if r.technology != "" {
			rel.Technology = r.technology
		}
		if r.description != "" {
			rel.Description = r.description
		}
		if r.proto != "" {
			rel.Proto = r.proto
		}
		if len(r.tags) > 0 {
			rel.Tags = r.tags
		}
		if r.external {
			rel.External = r.external
		}
		if r.person {
			rel.Person = r.person
		}

		sf.Relationships = append(sf.Relationships, rel)
	}

	sf.Sort()
	return sf, nil
}

func basenameFromRepoURL(repoURL string) string {
	u := strings.TrimSpace(repoURL)
	u = strings.TrimSuffix(u, "/")
	if u == "" {
		return ""
	}
	parts := strings.Split(u, "/")
	if len(parts) == 0 {
		return ""
	}
	return parts[len(parts)-1]
}

func (cp *Parser) validateNoMixedUsage() error {
	var (
		hasExplicit bool
		hasImplicit bool
	)

	for _, r := range cp.relationships {
		if r.serviceName != "" {
			hasExplicit = true
		} else {
			hasImplicit = true
		}
	}

	if hasExplicit && hasImplicit {
		return fmt.Errorf("mixed relationship definition patterns detected: some relationships use explicit patterns (service:name:action) while others use implicit patterns (service:action)")
	}

	return nil
}

func (cp *Parser) determineServiceName(r relationship, serviceFiles map[string]*servicefile.ServiceFile) (string, error) {
	if r.serviceName != "" {
		return r.serviceName, nil
	}

	for name := range serviceFiles {
		return name, nil
	}

	return "", fmt.Errorf("no service name found for relationship: %s", r)
}

func isEmptyRepository(serviceFiles []*servicefile.ServiceFile) bool {
	for _, sf := range serviceFiles {
		if sf.Info.Repository == "" {
			return true
		}
	}
	return false
}

func (cp *Parser) fillRepository(dir string, serviceFiles []*servicefile.ServiceFile) error {
	repoURL, err := detectGitRepository(dir)
	if err != nil {
		fmt.Printf("Couldn't detect git repository: %v\n", err.Error())
		return nil
	}

	for _, sf := range serviceFiles {
		if sf.Info.Repository == "" {
			sf.Info.Repository = repoURL
		}
	}

	return nil
}

func detectGitRepository(dir string) (string, error) {
	cmd := exec.Command("git", "config", "--get", "remote.origin.url")
	cmd.Dir = dir
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get git remote URL: %w", err)
	}

	url := strings.TrimSpace(string(output))
	if url == "" {
		return "", nil
	}

	return makeGitRepositoryURL(url), nil
}

func makeGitRepositoryURL(url string) string {
	if strings.HasPrefix(url, "git@") {
		url = strings.TrimPrefix(url, "git@")

		parts := strings.SplitN(url, ":", 2)
		if len(parts) == 2 {
			url = "https://" + parts[0] + "/" + parts[1]
		}
	}

	url = strings.TrimSuffix(url, ".git")

	return url
}
