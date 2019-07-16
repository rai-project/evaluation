package cmd

import (
	"github.com/spf13/cobra"
)

var (
	barPlot    bool
	boxPlot    bool
	openPlot   bool
	plotPath   string
	sortOutput bool
	topLayers  int
)

var layerCmd = &cobra.Command{
	Use: "layer",
	Aliases: []string{
		"layers",
	},
	Short: "Get evaluation model layer analysis from framework traces in a database",
}

func init() {
	layerCmd.PersistentFlags().BoolVar(&barPlot, "bar_plot", false, "generates a bar plot of the layers")
	layerCmd.PersistentFlags().BoolVar(&boxPlot, "box_plot", false, "generates a box plot of the layers")
	layerCmd.PersistentFlags().BoolVar(&piePlot, "pie_plot", false, "generates a pie plot of the layers")
	layerCmd.PersistentFlags().BoolVar(&openPlot, "open_plot", false, "opens the plot of the layers")
	layerCmd.PersistentFlags().StringVar(&plotPath, "plot_path", "", "output file for the layer plot")
	layerCmd.PersistentFlags().IntVar(&topLayers, "top_layers", -1, "consider only the top k layers ranked by duration")

	layerCmd.AddCommand(layerInfoCmd)
	layerCmd.AddCommand(layerDurationCmd)
	layerCmd.AddCommand(layerOcurrenceCmd)
	layerCmd.AddCommand(layerLatencyCmd)
	layerCmd.AddCommand(layerMemoryCmd)
	layerCmd.AddCommand(layerCUDAKernelCmd)
}
