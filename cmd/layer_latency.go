package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/rai-project/evaluation"
	"github.com/spf13/cobra"
)

var layerLatencyCmd = &cobra.Command{
	Use:     "latency",
	Aliases: []string{},
	Short:   "Get model layer latency information from framework traces in a database",
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
			plotPath = evaluation.TempFile("", "layer_latency_plot_*.html")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		run := func() error {
			evals, err := getEvaluations()
			if err != nil {
				return err
			}

			summary0, err := evals.SummaryLayerInformations(performanceCollection)
			if err != nil {
				return err
			}
			summary := evaluation.SummaryLayerLatencyInformations(summary0)

			if sortLayer {
				sort.Slice(summary, func(ii, jj int) bool {
					return evaluation.TrimmedMeanInt64Slice(summary[ii].Durations, 0) > evaluation.TrimmedMeanInt64Slice(summary[jj].Durations, 0)
				})
			}

			if barPlot {
				if openPlot {
					return summary.OpenBarPlot()
				}
				err := summary.WriteBarPlot(plotPath)
				if err != nil {
					return err
				}
				fmt.Println("Created plot in " + plotPath)
				return nil
			}

			if boxPlot {
				if openPlot {
					return summary.OpenBoxPlot()
				}
				err := summary.WriteBoxPlot(plotPath)
				if err != nil {
					return err
				}
				fmt.Println("Created plot in " + plotPath)
				return nil
			}

			writer := NewWriter(evaluation.SummaryLayerInformation{})
			defer writer.Close()

			for _, lyr := range summary0 {
				writer.Row(lyr)
			}
			return nil
		}

		return forallmodels(run)
	},
}
