package cmd

import (
	"path/filepath"

	"github.com/rai-project/evaluation"
	"github.com/spf13/cobra"
)

var layersTreeCmd = &cobra.Command{
	Use: "layer_tree",
	Aliases: []string{
		"layertree",
		"treemap",
	},
	Short: "Get evaluation layer tree information from CarML",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if modelName == "all" && outputFormat == "json" {
			outputFileName = filepath.Join(mlArcAssetsPath, "layer_tree")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		run := func() error {
			evals, err := getEvaluations()
			if err != nil {
				return err
			}

			lyrs, err := evals.LayerInformationTreeSummary(performanceCollection)

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
