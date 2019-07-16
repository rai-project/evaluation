package cmd

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/rai-project/evaluation"
	"github.com/spf13/cobra"
)

var cudaKernelDurationCmd = &cobra.Command{
	Use:     "info",
	Aliases: []string{},
	Short:   "Get model cuda kernel information from system library traces in a database. Specify model name as `all` to list information of all the models.",
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

			cudaKernelInfos, err := evals.CUDAKernelInformationSummary(performanceCollection)
			if err != nil {
				return err
			}

			if sortOutput || topLayers != -1 {
				sort.Sort(cudaKernelInfos)
				if topKernels != -1 {
					if topKernels >= len(cudaKernelInfos) {
						topKernels = len(cudaKernelInfos)
					}
					cudaKernelInfos = cudaKernelInfos[:topKernels]
				}
			}

			var writer *Writer
			if len(cudaKernelInfos) == 0 {
				writer = NewWriter(evaluation.SummaryCUDAKernelInformation{})
				defer writer.Close()
			}
			writer = NewWriter(cudaKernelInfos[0])
			defer writer.Close()

			for _, elem := range cudaKernelInfos {
				writer.Row(elem)
			}

			return nil
		}

		return forallmodels(run)
	},
}

func init() {
	cudaKernelDurationCmd.PersistentFlags().BoolVar(&sortOutput, "sort_output", false, "sort cuda kernel information by kernel duration")
	cudaKernelDurationCmd.PersistentFlags().StringVar(&kernelNameFilterString, "kernel_names", "", "filter out certain kernel (input must be mangled and is comma seperated)")
	cudaKernelDurationCmd.PersistentFlags().IntVar(&topKernels, "top_kernels", -1, "consider only the top k kernel ranked by duration")
}
