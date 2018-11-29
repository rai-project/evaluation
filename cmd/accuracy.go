package cmd

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/rai-project/evaluation"
)

var accuracyCmd = &cobra.Command{
	Use: "accuracy",
	Aliases: []string{
		"top_accuracy",
	},
	Short: "Get accuracy summary from MLModelScope",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if databaseName == "" {
			databaseName = defaultDatabaseName[cmd.Name()]
		}
		rootSetup()
		if modelName == "all" && outputFormat == "json" && outputFileName == "" {
			outputFileName = filepath.Join(mlArcWebAssetsPath, "accuracy")
		}
		if overwrite && isExists(outputFileName) {
			os.RemoveAll(outputFileName)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		run := func() error {
			accs, err := predictAccuracyInformationSummary()
			if err != nil {
				return err
			}

			writer := NewWriter(evaluation.SummaryPredictAccuracyInformation{})
			defer writer.Close()

			for _, acc := range accs {
				writer.Row(acc)
			}

			return nil
		}
		return forallmodels(run)
	},
}

func predictAccuracyInformationSummary() (evaluation.SummaryPredictAccuracyInformations, error) {
	evals, err := getEvaluations()
	if err != nil {
		return nil, err
	}

	accs, err := evals.PredictAccuracyInformationSummary(modelAccuracyCollection)
	if err != nil {
		return nil, err
	}

	return accs.Group()
}
