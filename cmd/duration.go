package cmd

import (
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/rai-project/evaluation"
)

var durationCmd = &cobra.Command{
	Use: "duration",
	Aliases: []string{
		"durations",
	},
	Short: "Get evaluation duration summary from CarML",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if modelName == "all" && outputFormat == "json" {
			outputFileName = filepath.Join(mlArcAssetsPath, "duration")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		run := func() error {
			durs, err := predictDurationInformationSummary()
			if err != nil {
				return err
			}

			writer := NewWriter(evaluation.SummaryPredictDurationInformation{})
			defer writer.Close()

			for _, dur := range durs {
				writer.Row(dur)
			}

			return nil
		}
		return forallmodels(run)
	},
}

func predictDurationInformationSummary() (evaluation.SummaryPredictDurationInformations, error) {
	evals, err := getEvaluations()
	if err != nil {
		return nil, err
	}
	return evals.PredictDurationInformationSummary(performanceCollection)
}
