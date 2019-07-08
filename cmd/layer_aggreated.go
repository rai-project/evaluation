package cmd

import (
	"os"
	"path/filepath"
	"sort"

	"github.com/rai-project/evaluation"
	"github.com/spf13/cobra"
)

var layerAggreCmd = &cobra.Command{
	Use: "layer_aggregated",
	Aliases: []string{
		"layers_aggregated",
	},
	Short: "Get model layer aggregated information from framework traces in a database",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if databaseName == "" {
			databaseName = defaultDatabaseName[cmd.Name()]
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
			plotPath = evaluation.TempFile("", "layer_aggre_plot_*.html")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		run := func() error {
			evals, err := getEvaluations()
			if err != nil {
				return err
			}

			summary, err := evals.LayerAggregatedInformationSummary(performanceCollection)
			if err != nil {
				return err
			}

			layerInfos := summary.LayerAggregatedInformations

			sort.Slice(layerInfos, func(ii, jj int) bool {
				return layerInfos[ii].Duration > layerInfos[jj].Duration
			})

			writer := NewWriter(evaluation.LayerAggregatedInformation{})
			defer writer.Close()

			for _, lyr := range layerInfos {
				writer.Row(lyr)
			}
			return nil
		}

		return forallmodels(run)
	},
}

func init() {
	// layerAggreCmd.PersistentFlags().BoolVar(&openPlot, "open_plot", false, "opens the plot of the layers")
	// layerAggreCmd.PersistentFlags().StringVar(&plotPath, "plot_path", "", "output file for the layer plot")
}
