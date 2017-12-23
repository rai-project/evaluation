package cmd

import (
	"github.com/spf13/cobra"
)

var allCmd = &cobra.Command{
	Use: "all",
	Aliases: []string{
		"eval_all",
	},
	Short: "Get all evaluation information from CarML",
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error

		argWoFlags := latencyCmd.Flags().Args()
		err = latencyCmd.PreRunE(latencyCmd, argWoFlags)
		if err != nil {
			return nil
		}
		err = latencyCmd.RunE(latencyCmd, argWoFlags)
		if err != nil {
			return nil
		}

		argWoFlags = layersCmd.Flags().Args()
		err = layersCmd.PreRunE(layersCmd, argWoFlags)
		if err != nil {
			return nil
		}
		err = layersCmd.RunE(layersCmd, argWoFlags)
		if err != nil {
			return nil
		}

		argWoFlags = layersTreeCmd.Flags().Args()
		err = layersTreeCmd.PreRunE(layersTreeCmd, argWoFlags)
		if err != nil {
			return nil
		}
		err = layersTreeCmd.RunE(layersTreeCmd, argWoFlags)
		if err != nil {
			return nil
		}

		argWoFlags = cudaLaunchCmd.Flags().Args()
		err = cudaLaunchCmd.PreRunE(cudaLaunchCmd, argWoFlags)
		if err != nil {
			return nil
		}
		err = cudaLaunchCmd.RunE(cudaLaunchCmd, argWoFlags)
		if err != nil {
			return nil
		}

		return nil
	},
}
