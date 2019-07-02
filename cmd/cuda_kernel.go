package cmd

import (
	"os"
	"path/filepath"

	"github.com/rai-project/evaluation"
	"github.com/spf13/cobra"
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

			writer := NewWriter(evaluation.LayerCUDAKernelInformation{})
			defer writer.Close()

			for _, elem := range layerCUDAKernelInfos {
				writer.Rows(elem)
			}

			return nil
		}

		return forallmodels(run)
	},
}
