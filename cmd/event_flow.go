package cmd

import (
	"os"
	"path/filepath"

	"github.com/rai-project/evaluation"
	"github.com/spf13/cobra"
)

var eventflowCmd = &cobra.Command{
	Use: "eventflow",
	Aliases: []string{
		"flow",
		"event_flow",
	},
	Short: "Get evaluation trace in event_flow format from CarML",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if databaseName == "" {
			databaseName = defaultDatabaseName[cmd.Name()]
		}
		rootSetup()
		if modelName == "all" && outputFormat == "json" && outputFileName == "" {
			outputFileName = filepath.Join(mlArcWebAssetsPath, "event_flow")
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

			flows, err := evals.EventFlowSummary(performanceCollection)

			writer := NewWriter(evaluation.SummaryEventFlow{})
			defer writer.Close()

			for _, flow := range flows {
				writer.Row(flow)
			}

			return nil
		}
		return forallmodels(run)
	},
}
