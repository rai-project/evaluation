package cmd

import (
	"os"
	"path/filepath"

	"github.com/rai-project/evaluation"
	"github.com/spf13/cobra"
)

var durationCmd = &cobra.Command{
	Use: "duration",
	Aliases: []string{
		"durations",
	},
	Short: "Get evaluation duration summary from MLModelScope",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if databaseName == "" {
			databaseName = defaultDatabaseName[cmd.Name()]
		}
		rootSetup()
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
			durs, err := predictDurationInformationSummary()
			if err != nil {
				return err
			}
			writer := NewWriter(evaluation.SummaryModelInformation{})
			defer writer.Close()

			for _, dur := range durs {
				writer.Row(dur)
			}

			return nil
		}
		return forallmodels(run)
	},
}

func predictDurationInformationSummary() (evaluation.SummaryModelInformations, error) {
	evals, err := getEvaluations()
	if err != nil {
		return nil, err
	}
	return evals.PredictDurationInformationSummary(performanceCollection)
}
