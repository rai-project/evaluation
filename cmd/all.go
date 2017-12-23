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
			return err
		}
		err = latencyCmd.RunE(latencyCmd, argWoFlags)
		if err != nil {
			return err
		}

		argWoFlags = layersCmd.Flags().Args()
		err = layersCmd.PreRunE(layersCmd, argWoFlags)
		if err != nil {
			return err
		}
		err = layersCmd.RunE(layersCmd, argWoFlags)
		if err != nil {
			return err
		}

		argWoFlags = layersTreeCmd.Flags().Args()
		err = layersTreeCmd.PreRunE(layersTreeCmd, argWoFlags)
		if err != nil {
			return err
		}
		err = layersTreeCmd.RunE(layersTreeCmd, argWoFlags)
		if err != nil {
			return err
		}

		argWoFlags = cudaLaunchCmd.Flags().Args()
		err = cudaLaunchCmd.PreRunE(cudaLaunchCmd, argWoFlags)
		if err != nil {
			return err
		}
		err = cudaLaunchCmd.RunE(cudaLaunchCmd, argWoFlags)
		if err != nil {
			return err
		}

		return nil
	},
}
