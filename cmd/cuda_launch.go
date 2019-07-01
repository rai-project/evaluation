package cmd

import (
	"os"
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
	Short: "Get evaluation kernel launch information from system library traces in a database. Specify model name as `all` to list information of all the models.",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if databaseName == "" {
			databaseName = defaultDatabaseName[cmd.Name()]
		}
		rootSetup()
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

			lst, err := evals.CUDALaunchInformationSummary(performanceCollection)

			writer := NewWriter(evaluation.SummaryCUDALaunchInformation{})
			defer writer.Close()

			for _, elem := range lst {
				writer.Rows(elem)
			}

			return nil
		}

		return forallmodels(run)
	},
}
