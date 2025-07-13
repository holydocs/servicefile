package golang

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/denchenko/servicefile/pkg/servicefile"
)

type CommentParser struct {
	serviceFile *servicefile.ServiceFile
}

func NewCommentParser() *CommentParser {
	return &CommentParser{
		serviceFile: &servicefile.ServiceFile{
			Relationships: []servicefile.Relationship{},
		},
	}
}

func (cp *CommentParser) Parse(dir string, recursive bool) (*servicefile.ServiceFile, error) {
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

		if err := cp.parseFile(path); err != nil {
			return fmt.Errorf("failed to parse %s: %w", path, err)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error walking the path: %w", err)
	}

	cp.serviceFile.Sort()

	return cp.serviceFile, nil
}

func (cp *CommentParser) parseFile(path string) error {
	fset := token.NewFileSet()

	f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("failed to parse %s: %w", path, err)
	}

	// Parse package-level comments
	for _, cg := range f.Comments {
		// Parse all comments in a group together to handle multi-line comments
		var commentText strings.Builder
		for _, c := range cg.List {
			commentText.WriteString(c.Text)
			commentText.WriteString("\n")
		}
		cp.parseCommentGroup(commentText.String())
	}

	// Parse type-level comments
	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.TypeSpec:
			if x.Doc != nil {
				var commentText strings.Builder
				for _, c := range x.Doc.List {
					commentText.WriteString(c.Text)
					commentText.WriteString("\n")
				}
				cp.parseCommentGroup(commentText.String())
			}
		case *ast.StructType:
			// For struct types, we need to look at the parent GenDecl or TypeSpec
			// The documentation is attached to the declaration, not the struct type itself
		}
		return true
	})

	return nil
}

func (cp *CommentParser) parseCommentGroup(commentGroup string) {
	lines := strings.Split(commentGroup, "\n")

	var currentRel *servicefile.Relationship

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		comment := strings.TrimSpace(line)
		comment = strings.TrimPrefix(comment, "//")
		comment = strings.TrimPrefix(comment, "/*")
		comment = strings.TrimSuffix(comment, "*/")
		comment = strings.TrimSpace(comment)

		if comment == "" {
			continue
		}

		if strings.HasPrefix(comment, "service:name") {
			parts := strings.SplitN(comment, " ", 2)
			if len(parts) == 2 {
				cp.serviceFile.Name = strings.TrimSpace(parts[1])
			}
			continue
		}

		// Parse service description (only if we're not currently building a relationship)
		if strings.HasPrefix(comment, "description:") && currentRel == nil {
			parts := strings.SplitN(comment, ":", 2)
			if len(parts) == 2 {
				cp.serviceFile.Description = strings.TrimSpace(parts[1])
			}
			continue
		}

		if strings.HasPrefix(comment, "service:") && !strings.HasPrefix(comment, "service:name") {
			// If there is a current relationship being built, add it to the list
			if currentRel != nil {
				cp.serviceFile.Relationships = append(cp.serviceFile.Relationships, *currentRel)
			}

			// Start a new relationship
			parts := strings.SplitN(comment, " ", 2)
			action := strings.TrimPrefix(parts[0], "service:")
			name := ""
			if len(parts) > 1 {
				name = strings.TrimSpace(parts[1])
			}
			currentRel = &servicefile.Relationship{
				Action: servicefile.RelationshipAction(action),
				Name:   name,
			}
			continue
		}

		if strings.HasPrefix(comment, "technology:") {
			parts := strings.SplitN(comment, ":", 2)
			if len(parts) == 2 && currentRel != nil {
				currentRel.Technology = strings.TrimSpace(parts[1])
			}
			continue
		}

		if strings.HasPrefix(comment, "description:") {
			parts := strings.SplitN(comment, ":", 2)
			if len(parts) == 2 && currentRel != nil {
				currentRel.Description = strings.TrimSpace(parts[1])
			}
			continue
		}
	}

	if currentRel != nil {
		cp.serviceFile.Relationships = append(cp.serviceFile.Relationships, *currentRel)
	}
}
