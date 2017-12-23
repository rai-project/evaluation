package cmd

import (
	"os"
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
		if databaseName == "" {
			databaseName = defaultDatabaseName[cmd.Name()]
		}
		rootSetup()
		if modelName == "all" && outputFormat == "json" && outputFileName == "" {
			outputFileName = filepath.Join(mlArcWebAssetsPath, "layers")
		}
		if overwrite && isExists(outputFileName) {
			os.RemoveAll(outputFileName)
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
