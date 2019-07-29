package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/rai-project/evaluation"
	"github.com/spf13/cobra"
)

var layerAggreMemoryCmd = &cobra.Command{
	Use:     "aggre_memory",
	Aliases: []string{},
	Short:   "Get model layer aggregated allocated memory information from framework traces in a database",
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

			summary0, err := evals.SummaryLayerAggreInformations(performanceCollection)
			if err != nil {
				return err
			}
			summary := evaluation.SummaryLayerAggreAllocatedMemoryInformations(summary0)

			if sortOutput {
				sort.Slice(summary, func(ii, jj int) bool {
					return summary[ii].AllocatedMemory > summary[jj].AllocatedMemory
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

			writer := NewWriter(evaluation.SummaryLayerAggreInformation{})
			defer writer.Close()

			for _, v := range summary {
				writer.Row(v)
			}
			return nil
		}

		return forallmodels(run)
	},
}
