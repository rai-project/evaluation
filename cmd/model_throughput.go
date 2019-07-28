package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/rai-project/evaluation"
	"github.com/spf13/cobra"
)

var modelThroughputCmd = &cobra.Command{
	Use:     "throughput",
	Aliases: []string{},
	Short:   "Get evaluation model throughput summary from model traces in a database",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if databaseName == "" {
			databaseName = defaultDatabaseName["model"]
		}
		err := rootSetup()
		if err != nil {
			return err
		}
		if modelName == "all" && outputFormat == "json" && outputFileName == "" {
			outputFileName = filepath.Join(mlArcWebAssetsPath, "duration")
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

			summary0, err := evals.SummaryModelInformations(performanceCollection)
			if err != nil {
				return err
			}

			sort.Sort(summary0)

			summary := evaluation.SummaryModelThroughputInformations(summary0)

			if openPlot {
				return summary.OpenBarPlot()
			}

			if barPlot {
				err := summary.WriteBarPlot(plotPath)
				if err != nil {
					return err
				}
				fmt.Println("Created plot in " + plotPath)
				return nil
			}

			writer := NewWriter(evaluation.SummaryModelInformation{})
			defer writer.Close()
			for _, v := range summary {
				writer.Row(v)
			}
			return nil
		}
		return forallmodels(run)
	},
}
