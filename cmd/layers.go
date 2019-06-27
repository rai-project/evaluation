package cmd

import (
	"os"
	"path/filepath"

	"github.com/rai-project/evaluation"
	"github.com/spf13/cobra"
)

var (
	mergeLayerInformationAcrossRuns bool
)

var layersCmd = &cobra.Command{
	Use: "layers",
	Aliases: []string{
		"layer",
	},
	Short: "Get evaluation layer information from database",
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

			if mergeLayerInformationAcrossRuns {
				summary, err := evals.AcrossEvaluationLayerInformationSummary(performanceCollection)
				if err != nil {
					return err
				}
				writer := NewWriter(evaluation.MeanLayerInformation{})
				defer writer.Close()

				for _, lyr := range summary[0].LayerInformations {
					writer.Row(evaluation.MeanLayerInformation{LayerInformation: lyr})
				}
				return nil
			}

			summaries, err := evals.LayerInformationSummary(performanceCollection)

			writer := NewWriter(evaluation.LayerInformation{})
			defer writer.Close()

			for _, summary := range summaries {
				writer.Row(summary)
			}

			return nil
		}

		return forallmodels(run)
	},
}

func init() {
	layersCmd.PersistentFlags().BoolVar(&mergeLayerInformationAcrossRuns, "merge_evaluations", false, "merges layer evaluations across runs")
}
