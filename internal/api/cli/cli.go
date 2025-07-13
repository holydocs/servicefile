package cli

import (
	"github.com/denchenko/servicefile/internal/api/cli/commands"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Long: `A CLI tool for managing ServiceFile.`,
	}

	cmd.AddCommand(
		commands.Parse(),
	)

	return cmd
}
