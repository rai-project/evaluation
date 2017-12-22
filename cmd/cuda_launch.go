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
			durs, err := predictDurationInformationSummary()
			if err != nil {
				return err
			}

			lats, err := durs.CUDALaunchInformationSummary()

			writer := NewWriter(evaluation.SummaryCUDALaunchInformation{})
			defer writer.Close()

			for _, lat := range lats {
				writer.Row(lat)
			}

			return nil
		}

		return forallmodels(run)
	},
}
