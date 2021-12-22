package options

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/philips-labs/slsa-provenance-action/lib/intoto"
)

// GenerateOptions Commandline flags used for the generate command.
type GenerateOptions struct {
	OutputPath     string
	ExtraMaterials []string
}

// GetOutputPath The location to write the provenance file.
func (o *GenerateOptions) GetOutputPath() (string, error) {
	if o.OutputPath == "" {
		return "", RequiredFlagError("output-path")
	}
	return o.OutputPath, nil
}

// GetExtraMaterials Additional material files to be used when generating provenance.
func (o *GenerateOptions) GetExtraMaterials() ([]intoto.Item, error) {
	var materials []intoto.Item

	for _, extra := range o.ExtraMaterials {
		file, err := os.Open(extra)
		if err != nil {
			return nil, fmt.Errorf("failed retrieving extra materials: %w", err)
		}
		defer file.Close()

		m, err := intoto.ReadMaterials(file)
		if err != nil {
			return nil, fmt.Errorf("failed retrieving extra materials for %s: %w", extra, err)
		}
		materials = append(materials, m...)
	}

	return materials, nil
}

// AddFlags Registers the flags with the cobra.Command.
func (o *GenerateOptions) AddFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&o.OutputPath, "output-path", "provenance.json", "The path to which the generated provenance should be written.")
	cmd.PersistentFlags().StringSliceVarP(&o.ExtraMaterials, "extra-materials", "m", nil, "The '${runner}' context value.")
}
