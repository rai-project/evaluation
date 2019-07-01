package cmd

import (
	"os"
	"path/filepath"

	"github.com/rai-project/evaluation"
	"github.com/spf13/cobra"
)

var latencyCmd = &cobra.Command{
	Use: "latency",
	Aliases: []string{
		"throughput",
	},
	Short: "Get model inference latency or throughput information from model traces in a database. Specify model name as `all` to list information of all the models.",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if databaseName == "" {
			databaseName = defaultDatabaseName[cmd.Name()]
		}
		rootSetup()
		if modelName == "all" && outputFormat == "json" && outputFileName == "" {
			outputFileName = filepath.Join(mlArcWebAssetsPath, "latency_throughput")
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

			lats, err := durs.ThroughputLatencySummary()
			if err != nil {
				return err
			}

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
