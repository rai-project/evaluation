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

var layerDurationCmd = &cobra.Command{
	Use:     "duration",
	Aliases: []string{},
	Short:   "Get model layer aggregated duration information from framework traces in a database",
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
			plotPath = evaluation.TempFile("", "layer_duration_plot_*.html")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		run := func() error {
			evals, err := getEvaluations()
			if err != nil {
				return err
			}

			summary0, err := evals.SummaryLayerAggregatedInformation(performanceCollection)
			if err != nil {
				return err
			}
			summary := evaluation.SummaryLayerDruationInformations(summary0)

			if sortOutput {
				sort.Slice(summary, func(ii, jj int) bool {
					return summary[ii].Duration > summary[jj].Duration
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

			writer := NewWriter(evaluation.SummaryLayerAggregatedInformation{})
			defer writer.Close()

			for _, lyr := range summary {
				writer.Row(lyr)
			}
			return nil
		}

		return forallmodels(run)
	},
}
