package cmd

import (
	"os"
	"path/filepath"

	"github.com/rai-project/evaluation"
	"github.com/spf13/cobra"
)

var modelInfoCmd = &cobra.Command{
	Use: "info",
	Aliases: []string{
		"durations",
		"duration",
	},
	Short: "Get evaluation model information summary from model traces in a database",
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

			summary, err := evals.SummaryModelInformation(performanceCollection)
			if err != nil {
				return err
			}
			writer := NewWriter(evaluation.SummaryModelInformation{})
			defer writer.Close()

			writer.Row(summary)

			return nil
		}
		return forallmodels(run)
	},
}
