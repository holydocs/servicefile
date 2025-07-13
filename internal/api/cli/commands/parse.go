package commands

import (
	"fmt"
	"os"

	"github.com/denchenko/servicefile/internal/parser/golang"
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
		Short: "Parse servicefile from source",
		RunE: func(_ *cobra.Command, _ []string) error {
			return parseServiceFile(dir, recursive, output)
		},
	}

	cmd.Flags().StringVarP(&dir, "dir", "d", ".", "Directory to analyze")
	cmd.Flags().BoolVarP(&recursive, "recursive", "r", true, "Recursively analyze subdirectories")
	cmd.Flags().StringVarP(&output, "output", "o", "servicefile.yaml", "Output file path for YAML")

	return cmd
}

func parseServiceFile(dir string, recursive bool, output string) error {
	parser := golang.NewCommentParser()

	serviceFile, err := parser.Parse(dir, recursive)
	if err != nil {
		return fmt.Errorf("error parsing service file: %w", err)
	}

	yamlData, err := yaml.Marshal(serviceFile)
	if err != nil {
		return fmt.Errorf("error marshaling to YAML: %w", err)
	}

	err = os.WriteFile(output, yamlData, 0644)
	if err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}

	fmt.Printf("ServiceFile generated and saved to: %s\n", output)

	return nil
}
