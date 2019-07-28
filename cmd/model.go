package cmd

import (
	"github.com/spf13/cobra"
)

var modelCmd = &cobra.Command{
	Use:     "model",
	Aliases: []string{},
	Short:   "Get evaluation model analysis from model traces in a database",
}

func init() {
	modelCmd.AddCommand(modelInfoCmd)
	modelCmd.AddCommand(modelThroughputCmd)
	modelCmd.AddCommand(modelLatencyCmd)
}
