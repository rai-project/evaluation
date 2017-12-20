package cmd

import (
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
