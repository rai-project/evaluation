package cmd

import (
	"os"
	"path/filepath"
	"sort"

	"github.com/rai-project/evaluation"
	"github.com/spf13/cobra"
)

var (
	sortByLatency  bool
	trimKernelName bool
	topLayers      int
)

var cudaKernelCmd = &cobra.Command{
	Use: "cuda_kernel",
	Aliases: []string{
		"cuda",
		"kernel",
		"kernels",
		"gpu_kernel",
		"gpu_kernels",
	},
	Short: "Get evaluation kernel launch information from system library traces in a database. Specify model name as `all` to list information of all the models.",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if databaseName == "" {
			databaseName = defaultDatabaseName[cmd.Name()]
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
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		run := func() error {
			evals, err := getEvaluations()
			if err != nil {
				return err
			}

			summary, err := evals.LayerCUDAKernelInformationSummary(performanceCollection)
			if err != nil {
				return err
			}

			layerCUDAKernelInfos := summary.LayerCUDAKernelInformations

			if sortByLatency || topLayers != -1 {
				sort.Slice(layerCUDAKernelInfos, func(ii, jj int) bool {
					return evaluation.TrimmedMean(layerCUDAKernelInfos[ii].Durations, 0) > evaluation.TrimmedMean(layerCUDAKernelInfos[jj].Durations, 0)
				})
				if topLayers != -1 {
					if topLayers >= len(layerCUDAKernelInfos) {
						topLayers = len(layerCUDAKernelInfos)
					}
					layerCUDAKernelInfos = layerCUDAKernelInfos[:topLayers]
				}
			}

			writer := NewWriter(layerCUDAKernelInfos[0])
			defer writer.Close()

			for _, elem := range layerCUDAKernelInfos {
				writer.Rows(elem)
			}

			return nil
		}

		return forallmodels(run)
	},
}

func init() {
	cudaKernelCmd.PersistentFlags().BoolVar(&sortByLatency, "sort_by_latency", false, "sort layer information by layer latency")
	cudaKernelCmd.PersistentFlags().IntVar(&topLayers, "top_layers", -1, "consider only the top k layers")
	cudaKernelCmd.PersistentFlags().BoolVar(&trimKernelName, "trim_name", true, "trim kernel names to the first `<`")
}
