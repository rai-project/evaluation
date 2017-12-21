package cmd

import (
	"path/filepath"

	"github.com/rai-project/evaluation"
	"github.com/spf13/cobra"
)

var latencyCmd = &cobra.Command{
	Use: "latency",
	Aliases: []string{
		"throughput",
	},
	Short: "Get evaluation latency or throughput information from CarML",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if modelName == "all" && outputFormat == "json" {
			outputFileName = filepath.Join(raiSrcPath, "ml-arc-web", "src", "assets", "latency_throughput")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		run := func() error {
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
		}

		return forallmodels(run)
	},
}

func init() {
}
