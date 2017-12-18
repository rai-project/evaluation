package cmd

import (
	"github.com/rai-project/evaluation"
	"github.com/spf13/cobra"
)

var latencyCmd = &cobra.Command{
	Use: "latency",
	Aliases: []string{
		"throughput",
	},
	Short: "Get evaluation latency or throughput information from CarML",
	RunE: func(cmd *cobra.Command, args []string) error {
		durs, err := predictDurationInformationSummary()
		if err != nil {
			return err
		}

		lats, err := durs.ThroughputLatencySummary()

		writer := NewWriter(evaluation.SummaryThroughputLatency{})
		defer writer.Close()

		for _, lat := range lats {
			writer.Row(lat)
		}

		return nil
	},
}

func init() {
}
