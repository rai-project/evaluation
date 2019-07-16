package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/rai-project/evaluation"
	"github.com/spf13/cobra"
)

var layerMemoryCmd = &cobra.Command{
	Use:     "memory",
	Aliases: []string{},
	Short:   "Get model layer memory information from framework traces in a database",
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
			plotPath = evaluation.TempFile("", "layer_memory_plot_*.html")
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
			summary := evaluation.SummaryLayerMemoryInformations(summary0)

			if sortLayer {
				sort.Slice(summary, func(ii, jj int) bool {
					return evaluation.TrimmedMeanInt64Slice(summary[ii].AllocatedBytes, 0) > evaluation.TrimmedMeanInt64Slice(summary[jj].AllocatedBytes, 0)
				})
			}

			if openPlot {
				return summary.OpenBarPlot()
			}

			if barPlot {
				err := summary.WriteBarPlot(plotPath)
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
