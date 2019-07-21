package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/rai-project/evaluation"
	"github.com/spf13/cobra"
)

var gpuKernelLayerDramReadCmd = &cobra.Command{
	Use:     "layer_dram_read",
	Aliases: []string{},
	Short:   "Get the total dram read of all GPU kernels  within each layer from system library traces in a database. Specify model name as `all` to list information of all the models.",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if databaseName == "" {
			databaseName = defaultDatabaseName["cuda_kernel"]
		}
		err := rootSetup()
		if err != nil {
			return err
		}
		if modelName == "all" && outputFormat == "json" && outputFileName == "" {
			outputFileName = filepath.Join(mlArcWebAssetsPath, "cuda_kernel_launch")
		}
		if overwrite && isExists(outputFileName) {
			os.RemoveAll(outputFileName)
		}

		if kernelNameFilterString != "" {
			kernelNameFilterList = strings.Split(kernelNameFilterString, ",")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		run := func() error {
			evals, err := getEvaluations()
			if err != nil {
				return err
			}

			summary0, err := evals.SummaryGPUKernelLayerAggreInformations(performanceCollection)
			if err != nil {
				return err
			}
			summary := evaluation.SummaryGPUKernelLayerDramReadInformations(summary0)

			if sortOutput {
				sort.Slice(summary, func(ii, jj int) bool {
					return summary[ii].Index > summary[jj].Index
				})
			}

			if barPlot {
				err := summary.WriteBarPlot(plotPath)
				if err != nil {
					return err
				}
				fmt.Println("Created plot in " + plotPath)
				return nil
			}

			if openPlot {
				return summary.OpenBarPlot()
			}

			var writer *Writer
			if len(summary0) == 0 {
				writer = NewWriter(evaluation.SummaryGPUKernelModelAggreInformation{})
				defer writer.Close()
			}
			writer = NewWriter(summary0[0])
			defer writer.Close()

			for _, elem := range summary0 {
				writer.Row(elem)
			}

			return nil
		}

		return forallmodels(run)
	},
}