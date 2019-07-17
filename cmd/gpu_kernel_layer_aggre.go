package cmd

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/rai-project/evaluation"
	"github.com/spf13/cobra"
)

var gpuKernelLayerAggreCmd = &cobra.Command{
	Use:     "layer_aggre",
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

			gpuInfos, err := evals.SummaryGPUKernelLayerAggreInformations(performanceCollection)
			if err != nil {
				return err
			}

			if sortOutput || topKernels != -1 {
				sort.Sort(gpuInfos)
				if topKernels != -1 {
					if topKernels >= len(gpuInfos) {
						topKernels = len(gpuInfos)
					}
					gpuInfos = gpuInfos[:topKernels]
				}
			}

			var writer *Writer
			if len(gpuInfos) == 0 {
				writer = NewWriter(evaluation.SummaryGPUKernelModelAggreInformation{})
				defer writer.Close()
			}
			writer = NewWriter(gpuInfos[0])
			defer writer.Close()

			for _, elem := range gpuInfos {
				writer.Row(elem)
			}

			return nil
		}

		return forallmodels(run)
	},
}

func init() {
	gpuKernelLayerAggreCmd.PersistentFlags().BoolVar(&sortOutput, "sort_output", false, "sort cuda kernel information by kernel duration")
	gpuKernelLayerAggreCmd.PersistentFlags().StringVar(&kernelNameFilterString, "kernel_names", "", "filter out certain kernel (input must be mangled and is comma seperated)")
	gpuKernelLayerAggreCmd.PersistentFlags().IntVar(&topKernels, "top_kernels", -1, "consider only the top k kernel ranked by duration")
}
