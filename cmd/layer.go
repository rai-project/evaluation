package cmd

import (
	"github.com/spf13/cobra"
)

var (
	topLayers int
)

var layerCmd = &cobra.Command{
	Use: "layer",
	Aliases: []string{
		"layers",
	},
	Short: "Get evaluation model layer analysis from framework traces in a database",
}

func init() {
	layerCmd.PersistentFlags().IntVar(&topLayers, "top_layers", -1, "consider only the top k layers ranked by duration")

	layerCmd.AddCommand(layerInfoCmd)
	layerCmd.AddCommand(layerLatencyCmd)
	layerCmd.AddCommand(layerAllocatedMemoryCmd)
	layerCmd.AddCommand(layerAggreInfoCmd)
	layerCmd.AddCommand(layerAggreLatencyCmd)
	layerCmd.AddCommand(layerAggreOcurrenceCmd)
	layerCmd.AddCommand(layerAggreMemoryCmd)
}
