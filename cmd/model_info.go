package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/rai-project/evaluation"
	"github.com/spf13/cobra"
)

var modelInfoCmd = &cobra.Command{
	Use:     "info",
	Aliases: []string{},
	Short:   "Get evaluation model information summary from model traces in a database",
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

			if plotAll {
				plotPath = outputFileName + "_latency.html"
				summary1 := evaluation.SummaryModelLatencyInformations(summary0)
				err := summary1.WriteBarPlot(plotPath)
				if err != nil {
					return err
				}
				fmt.Println("Created plot in " + plotPath)

				plotPath = outputFileName + "_throughtput.html"
				summary2 := evaluation.SummaryModelThroughputInformations(summary0)
				err = summary2.WriteBarPlot(plotPath)
				if err != nil {
					return err
				}
				fmt.Println("Created plot in " + plotPath)
			}

			writer := NewWriter(evaluation.SummaryModelInformation{})
			defer writer.Close()
			for _, v := range summary0 {
				writer.Row(v)
			}
			return nil
		}
		return forallmodels(run)
	},
}
