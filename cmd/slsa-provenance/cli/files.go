package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/philips-labs/slsa-provenance-action/cmd/slsa-provenance/cli/options"
	"github.com/philips-labs/slsa-provenance-action/lib/github"
)

// Files creates an instance of *cobra.Command to manage file provenance
func Files() *cobra.Command {
	o := &options.FilesOptions{}

	cmd := &cobra.Command{
		Use:   "files",
		Short: "Generate provenance on file assets",
		RunE: func(cmd *cobra.Command, args []string) error {
			artifactPath, err := o.GetArtifactPath()
			if err != nil {
				return err
			}
			outputPath, err := o.GetOutputPath()
			if err != nil {
				return err
			}

			materials, err := o.GetExtraMaterials()
			if err != nil {
				return err
			}

			env, err := github.NewEnvironment()
			if err != nil {
				return err
			}

			stmt, err := env.GenerateProvenanceStatement(cmd.Context(), artifactPath, materials...)
			if err != nil {
				return fmt.Errorf("failed to generate provenance: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Saving provenance to %s\n", outputPath)

			return env.PersistProvenanceStatement(cmd.Context(), stmt, outputPath)
		},
	}

	o.AddFlags(cmd)

	return cmd
}
