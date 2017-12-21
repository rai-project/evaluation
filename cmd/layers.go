package cmd

import (
	"path/filepath"

	"github.com/rai-project/evaluation"
	"github.com/spf13/cobra"
)

var layersCmd = &cobra.Command{
	Use: "layers",
	Aliases: []string{
		"layer",
	},
	Short: "Get evaluation layer  information from CarML",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if modelName == "all" && outputFormat == "json" {
			outputFileName = filepath.Join(raiSrcPath, "ml-arc-web", "src", "assets", "layers")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		run := func() error {
			evals, err := getEvaluations()
			if err != nil {
				return err
			}

			lyrs, err := evals.LayerInformationSummary(performanceCollection)

			writer := NewWriter(evaluation.SummaryLayerInformation{})
			defer writer.Close()

			for _, lyr := range lyrs {
				writer.Row(lyr)
			}

			return nil
		}

		return forallmodels(run)
	},
}
