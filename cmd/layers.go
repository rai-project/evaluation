package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/rai-project/evaluation"
	"github.com/spf13/cobra"
)

var (
	listRuns      bool
	sortByLatency bool
	plotLayers    bool
	openPlot      bool
	topLayers     int
	plotPath      string
)

var layersCmd = &cobra.Command{
	Use: "layers",
	Aliases: []string{
		"layer",
	},
	Short: "Get model layer information from framework traces in a database",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if databaseName == "" {
			databaseName = defaultDatabaseName[cmd.Name()]
		}
		rootSetup()
		if modelName == "all" && outputFormat == "json" && outputFileName == "" {
			outputFileName = filepath.Join(mlArcWebAssetsPath, "layers")
		}
		if overwrite && isExists(outputFileName) {
			os.RemoveAll(outputFileName)
		}
		if plotLayers == true && plotPath == "" {
			plotPath = evaluation.TempFile("", "layer_plot_*.html")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		run := func() error {
			evals, err := getEvaluations()
			if err != nil {
				return err
			}

			if listRuns {
				summaries, err := evals.LayerInformationSummary(performanceCollection)
				if err != nil {
					return err
				}

				writer := NewWriter(evaluation.LayerInformation{})
				defer writer.Close()
				for _, summary := range summaries {
					writer.Row(summary)
				}

				return nil
			}

			summary, err := evals.AcrossEvaluationLayerInformationSummary(performanceCollection)
			if err != nil {
				return err
			}

			layers := summary[0].LayerInformations
			meanLayers := make(evaluation.MeanLayerInformations, len(layers))
			for ii, layer := range layers {
				meanLayers[ii] = evaluation.MeanLayerInformation{LayerInformation: layer}
			}
			if sortByLatency || topLayers != -1 {
				sort.Slice(meanLayers, func(ii, jj int) bool {
					return evaluation.TrimmedMean(meanLayers[ii].Durations, 0) > evaluation.TrimmedMean(meanLayers[jj].Durations, 0)
				})
				if topLayers != -1 {
					if topLayers >= len(meanLayers) {
						topLayers = len(meanLayers)
					}
					meanLayers = meanLayers[:topLayers]
				}
			}

			if openPlot {
				return meanLayers.OpenPlot()
			}
			if plotLayers {
				err := meanLayers.WritePlot(plotPath)
				if err == nil {
					fmt.Println("Created plot in " + plotPath)
				}
				return err
			}

			writer := NewWriter(evaluation.MeanLayerInformation{})
			defer writer.Close()

			for _, lyr := range meanLayers {
				writer.Row(lyr)
			}
			return nil
		}

		return forallmodels(run)
	},
}

func init() {
	layersCmd.PersistentFlags().BoolVar(&listRuns, "list_runs", false, "list evaluations")
	layersCmd.PersistentFlags().BoolVar(&sortByLatency, "sort_by_latency", false, "sort layer information by layer latency")
	layersCmd.PersistentFlags().BoolVar(&plotLayers, "plot", false, "generates a plot of the layers")
	layersCmd.PersistentFlags().BoolVar(&openPlot, "open_plot", false, "opens the plot of the layers")
	layersCmd.PersistentFlags().IntVar(&topLayers, "top_layers", -1, "consider only the top k layers")
	layersCmd.PersistentFlags().StringVar(&plotPath, "plot_path", "", "output file for the layer plot")
}
