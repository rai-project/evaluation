package cmd

import (
	"os"
	"path/filepath"
	"sort"

	"github.com/rai-project/evaluation"
	"github.com/spf13/cobra"
)

var (
	listRuns bool
)

var layerInfoCmd = &cobra.Command{
	Use:     "info",
	Aliases: []string{},
	Short:   "Get model layer information from framework traces in a database",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if databaseName == "" {
			databaseName = defaultDatabaseName["layer"]
		}
		err := rootSetup()
		if err != nil {
			return err
		}
		if modelName == "all" && outputFormat == "json" && outputFileName == "" {
			outputFileName = filepath.Join(mlArcWebAssetsPath, "layers")
		}
		if overwrite && isExists(outputFileName) {
			os.RemoveAll(outputFileName)
		}
		if plotPath == "" {
			plotPath = evaluation.TempFile("", "layer_plot_*.html")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		run := func() error {
			evals, err := getEvaluations()
			if err != nil {
				return err
			}

			summary, err := evals.SummaryLayerInformations(performanceCollection)
			if err != nil {
				return err
			}

			if listRuns {
				writer := NewWriter(evaluation.SummaryLayerInformation{})
				defer writer.Close()
				for _, lyr := range summary {
					writer.Row(lyr)
				}
				return nil
			}

			if sortOutput || topLayers != -1 {
				sort.Slice(summary, func(ii, jj int) bool {
					return evaluation.TrimmedMeanInt64Slice(summary[ii].Durations, 0) > evaluation.TrimmedMeanInt64Slice(summary[jj].Durations, 0)
				})
				if topLayers != -1 {
					if topLayers >= len(summary) {
						topLayers = len(summary)
					}
					summary = summary[:topLayers]
				}
			}

			writer := NewWriter(evaluation.SummaryMeanLayerInformation{})
			defer writer.Close()
			for _, lyr := range summary {
				writer.Row(evaluation.SummaryMeanLayerInformation(lyr))
			}
			return nil
		}

		return forallmodels(run)
	},
}

func init() {
	layerInfoCmd.PersistentFlags().BoolVar(&listRuns, "list_runs", false, "list evaluations")
	layerInfoCmd.PersistentFlags().IntVar(&topLayers, "top_layers", -1, "consider only the top k layers")
}
