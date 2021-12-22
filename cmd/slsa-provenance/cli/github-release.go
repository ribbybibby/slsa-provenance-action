package cli

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/philips-labs/slsa-provenance-action/cmd/slsa-provenance/cli/options"
	"github.com/philips-labs/slsa-provenance-action/lib/github"
	"github.com/philips-labs/slsa-provenance-action/lib/transport"
)

// GitHubRelease creates an instance of *cobra.Command to manage GitHub release provenance
func GitHubRelease() *cobra.Command {
	o := options.GitHubReleaseOptions{}

	cmd := &cobra.Command{
		Use:   "github-release",
		Short: "Generate provenance on GitHub release assets",
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

			tagName, err := o.GetTagName()
			if err != nil {
				return err
			}

			ghToken := os.Getenv("GITHUB_TOKEN")
			if ghToken == "" {
				return errors.New("GITHUB_TOKEN environment variable not set")
			}
			tc := github.NewOAuth2Client(cmd.Context(), func() string { return ghToken })
			tc.Transport = transport.TeeRoundTripper{
				RoundTripper: tc.Transport,
				Writer:       cmd.OutOrStdout(),
			}
			rc := github.NewReleaseClient(tc)
			env, err := github.NewEnvironment()
			if err != nil {
				return err
			}
			releaseEnv := github.NewReleaseEnvironment(env, tagName, rc)

			stmt, err := releaseEnv.GenerateProvenanceStatement(cmd.Context(), artifactPath, materials...)
			if err != nil {
				return fmt.Errorf("failed to generate provenance: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Saving provenance to %s\n", outputPath)

			return releaseEnv.PersistProvenanceStatement(cmd.Context(), stmt, outputPath)
		},
	}

	o.AddFlags(cmd)

	return cmd
}
