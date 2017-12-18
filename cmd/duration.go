package cmd

import (
	"github.com/spf13/cobra"

	"github.com/rai-project/evaluation"
)

var durationCmd = &cobra.Command{
	Use: "duration",
	Aliases: []string{
		"durations",
	},
	Short: "Get evaluation duration summary from CarML",
	RunE: func(cmd *cobra.Command, args []string) error {
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
	},
}

func predictDurationInformationSummary() (evaluation.SummaryPredictDurationInformations, error) {
	evals, err := getEvaluations()
	if err != nil {
		return nil, err
	}
	return evals.PredictDurationInformationSummary(performanceCollection)
}
