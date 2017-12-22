package cmd

import (
	"path/filepath"

	"github.com/rai-project/evaluation"
	"github.com/spf13/cobra"
)

var cudaLaunchCmd = &cobra.Command{
	Use: "cuda_launch",
	Aliases: []string{
		"kernel_launch",
		"kernels",
	},
	Short: "Get evaluation kernel launch information from CarML",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if modelName == "all" && outputFormat == "json" {
			outputFileName = filepath.Join(mlArcWebAssetsPath, "cuda_kernel_launch")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		run := func() error {

			evals, err := getEvaluations()
			if err != nil {
				return err
			}

			lst, err := evals.CUDALaunchInformationSummary(performanceCollection)

			writer := NewWriter(evaluation.SummaryCUDALaunchInformation{})
			defer writer.Close()

			for _, elem := range lst {
				writer.Row(elem)
			}

			return nil
		}

		return forallmodels(run)
	},
}
