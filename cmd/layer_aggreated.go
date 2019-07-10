package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/rai-project/evaluation"
	"github.com/spf13/cobra"
)

var (
	piePlot bool
)

var layerAggregatedCmd = &cobra.Command{
	Use:     "aggregated",
	Aliases: []string{},
	Short:   "Get model layer aggregated information from framework traces in a database",
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

			summary0, err := evals.LayerAggregatedInformationSummary(performanceCollection)
			if err != nil {
				return err
			}

			summary := evaluation.SummaryLayerAggregatedInformationOccurrences{summary0}

			layerInfos := summary.LayerAggregatedInformations

			if sortLayer {
				sort.Slice(layerInfos, func(ii, jj int) bool {
					return layerInfos[ii].Duration > layerInfos[jj].Duration
				})
			}

			if openPlot {
				return summary.OpenPiePlot()
			}

			if piePlot {
				err := summary.WritePiePlot(plotPath)
				if err != nil {
					return err
				}
				fmt.Println("Created plot in " + plotPath)
				return nil
			}

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
	layerAggregatedCmd.PersistentFlags().BoolVar(&piePlot, "pie_plot", false, "generates a pie plot of the layers")
}
