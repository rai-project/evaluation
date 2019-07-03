package cmd

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

var (
	sortLayer              bool
	trimOutput             int
	sortOutput             bool
	kernelNameFilterString string
	kernelNameFilterList   = []string{}
	topLayers              int
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

			summary, err := evals.LayerCUDAKernelInformationSummary(performanceCollection)
			if err != nil {
				return err
			}

			layerCUDAKernelInfos := summary.LayerCUDAKernelInformations

			if sortOutput || topLayers != -1 {
				sort.Sort(layerCUDAKernelInfos)

				if topLayers != -1 {
					if topLayers >= len(layerCUDAKernelInfos) {
						topLayers = len(layerCUDAKernelInfos)
					}
					layerCUDAKernelInfos = layerCUDAKernelInfos[:topLayers]
				}

				for ii := range layerCUDAKernelInfos {
					kernelInfo := layerCUDAKernelInfos[ii]
					sort.Sort(kernelInfo)
					layerCUDAKernelInfos[ii] = kernelInfo
				}
			}

			writer := NewWriter(layerCUDAKernelInfos)
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
	cudaKernelCmd.PersistentFlags().BoolVar(&sortOutput, "sort_output", false, "sort cuda kernel information by layer index and then kernel duration")
	cudaKernelCmd.PersistentFlags().StringVar(&kernelNameFilterString, "kernel_names", "", "filter out certain kernel (input must be mangled and is comma seperated)")
	cudaKernelCmd.PersistentFlags().IntVar(&topLayers, "top_layers", -1, "consider only the top k layers")
	cudaKernelCmd.PersistentFlags().IntVar(&trimOutput, "trim_output", "-1", "trim the output to a certain number of charaters")
}
