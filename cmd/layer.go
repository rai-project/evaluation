package cmd

import (
	"github.com/spf13/cobra"
)

var layerCmd = &cobra.Command{
	Use: "layer",
	Aliases: []string{
		"layers",
	},
	Short: "Get evaluation model layer analysis from framework traces in a database",
}

func init() {
	layerCmd.AddCommand(layerInfoCmd)
	layerCmd.AddCommand(layerDurationCmd)
	layerCmd.AddCommand(layerOcurrenceCmd)
}
