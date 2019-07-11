package cmd

import (
	"github.com/spf13/cobra"
)

var modelCmd = &cobra.Command{
	Use:     "model",
	Aliases: []string{},
	Short:   "Get evaluation model analysis from model traces in a database",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	modelCmd.AddCommand(modelInfoCmd)
	modelCmd.AddCommand(modelLatencyCmd)
}
