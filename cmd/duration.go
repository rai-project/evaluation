package cmd

import (
	"github.com/spf13/cobra"

	"github.com/rai-project/evaluation"
	udb "upper.io/db.v3"
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

	filter := udb.Cond{}
	if modelName != "" {
		filter["model.name"] = modelName
	}
	if modelVersion != "" {
		filter["model.version"] = modelVersion
	}
	if frameworkName != "" {
		filter["framework.name"] = frameworkName
	}
	if frameworkVersion != "" {
		filter["framework.version"] = frameworkVersion
	}
	if machineArchitecture != "" {
		filter["machinearchitecture"] = machineArchitecture
	}
	if hostName != "" {
		filter["hostname"] = hostName
	}
	evals, err := evaluationCollection.Find(filter)
	if err != nil {
		return nil, err
	}

	return evaluation.Evaluations(evals).PredictDurationInformationSummary(performanceCollection)
}
