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

var gpuKernelLayerAggreInfoCmd = &cobra.Command{
	Use:     "layer_aggre_info",
	Aliases: []string{},
	Short:   "Get gpu information aggregated within each layer from system library traces in a database. Specify model name as `all` to list information of all the models.",
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

			if sortOutput || topKernels != -1 {
				sort.Sort(summary0)
				if topKernels != -1 {
					if topKernels >= len(summary0) {
						topKernels = len(summary0)
					}
					summary0 = summary0[:topKernels]
				}
			}

			if plotAll {
				plotPath = outputFileName + "_flops.html"
				summary1 := evaluation.SummaryGPUKernelLayerFlopsInformations(summary0)
				err := summary1.WriteBarPlot(plotPath)
				if err != nil {
					return err
				}
				fmt.Println("Created plot in " + plotPath)

				plotPath = outputFileName + "_dram_read.html"
				summary2 := evaluation.SummaryGPUKernelLayerDramReadInformations(summary0)
				err = summary2.WriteBarPlot(plotPath)
				if err != nil {
					return err
				}
				fmt.Println("Created plot in " + plotPath)

				plotPath = outputFileName + "_dram_write.html"
				summary3 := evaluation.SummaryGPUKernelLayerDramWriteInformations(summary0)
				err = summary3.WriteBarPlot(plotPath)
				if err != nil {
					return err
				}
				fmt.Println("Created plot in " + plotPath)

				plotPath = outputFileName + "_achieved_occupancy.html"
				summary4 := evaluation.SummaryGPUKernelLayerAchievedOccupancyInformations(summary0)
				err = summary4.WriteBarPlot(plotPath)
				if err != nil {
					return err
				}
				fmt.Println("Created plot in " + plotPath)

				plotPath = outputFileName + "_gpu_cpu.html"
				summary5 := evaluation.SummaryGPUKernelLayerGPUCPUInformations(summary0)
				err = summary5.WriteBarPlot(plotPath)
				if err != nil {
					return err
				}
				fmt.Println("Created plot in " + plotPath)
			}

			var writer *Writer
			if len(summary0) == 0 {
				writer = NewWriter(evaluation.SummaryGPUKernelLayerAggreInformation{})
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
