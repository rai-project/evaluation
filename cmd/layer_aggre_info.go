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

var layerAggreInfoCmd = &cobra.Command{
	Use:     "aggre_info",
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

			summary0, err := evals.SummaryLayerAggreInformations(performanceCollection)
			if err != nil {
				return err
			}

			if sortOutput {
				sort.Slice(summary0, func(ii, jj int) bool {
					return summary0[ii].Duration > summary0[jj].Duration
				})
			}

			if plotAll {
				plotPath = outputFileName + "_occurence.html"
				summary1 := evaluation.SummaryLayerAggreOccurrenceInformations(summary0)
				err := summary1.WritePiePlot(plotPath)
				if err != nil {
					return err
				}
				fmt.Println("Created plot in " + plotPath)

				plotPath = outputFileName + "_latency.html"
				summary2 := evaluation.SummaryLayerAggreDurationInformations(summary0)
				err = summary2.WritePiePlot(plotPath)
				if err != nil {
					return err
				}
				fmt.Println("Created plot in " + plotPath)

				plotPath = outputFileName + "_allocated_memory.html"
				summary3 := evaluation.SummaryLayerAggreAllocatedMemoryInformations(summary0)
				err = summary3.WritePiePlot(plotPath)
				if err != nil {
					return err
				}
				fmt.Println("Created plot in " + plotPath)
			}

			writer := NewWriter(evaluation.SummaryLayerAggreInformation{})
			defer writer.Close()
			for _, v := range summary0 {
				writer.Row(v)
			}
			return nil
		}

		return forallmodels(run)
	},
}
