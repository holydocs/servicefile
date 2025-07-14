package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/denchenko/servicefile/internal/parser/golang"
	"github.com/denchenko/servicefile/pkg/servicefile"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func Parse() *cobra.Command {
	var (
		dir       string
		recursive bool
		output    string
	)

	cmd := &cobra.Command{
		Use:   "parse",
		Short: "Parse servicefiles from source",
		RunE: func(_ *cobra.Command, _ []string) error {
			return parseServiceFiles(dir, recursive, output)
		},
	}

	cmd.Flags().StringVarP(&dir, "dir", "d", ".", "Directory to analyze")
	cmd.Flags().BoolVarP(&recursive, "recursive", "r", true, "Recursively analyze subdirectories")
	cmd.Flags().StringVarP(&output, "output", "o", "servicefile.yaml", "Output file path suffix for YAML")

	return cmd
}

func parseServiceFiles(dir string, recursive bool, output string) error {
	parser := golang.NewCommentParser()

	serviceFiles, err := parser.Parse(dir, recursive)
	if err != nil {
		return fmt.Errorf("error parsing service file: %w", err)
	}

	if len(serviceFiles) == 0 {
		return fmt.Errorf("no services found in the specified directory")
	}

	if len(serviceFiles) == 1 {
		sf := serviceFiles[0]

		if err := saveServiceFileToYAML(sf, output); err != nil {
			return fmt.Errorf("error saving service file to %s: %w", output, err)
		}

		fmt.Printf("ServiceFile generated and saved to: %s\n", output)

		return nil
	}

	for _, sf := range serviceFiles {
		filepath := fmt.Sprintf("%s.%s", strings.ToLower(sf.Name), output)

		if err := saveServiceFileToYAML(sf, filepath); err != nil {
			return fmt.Errorf("error saving service file to %s: %w", filepath, err)
		}

		fmt.Printf("ServiceFile for '%s' generated and saved to: %s\n", sf.Name, filepath)
	}

	return nil
}

func saveServiceFileToYAML(sf *servicefile.ServiceFile, filepath string) error {
	yamlData, err := yaml.Marshal(sf)
	if err != nil {
		return fmt.Errorf("error marshaling to YAML: %w", err)
	}

	err = os.WriteFile(filepath, yamlData, 0644)
	if err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}

	return nil
}
