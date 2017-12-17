package cmd

import (
	"bytes"
	"encoding/csv"
	"io"
	"os"

	"github.com/Unknwon/com"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"github.com/rai-project/evaluation"
	udb "upper.io/db.v3"
)

var ()

var durationCmd = &cobra.Command{
	Use: "duration",
	Aliases: []string{
		"durations",
	},
	Short: "Get evaluation duration summary from CarML",
	RunE: func(cmd *cobra.Command, args []string) error {
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
			return err
		}

		durs, err := evaluation.Evaluations(evals).PredictDurationInformationSummary(performanceCollection)
		if err != nil {
			return err
		}

		var output io.Writer = os.Stdout
		if outputFileName != "" {
			buf := &bytes.Buffer{}
			defer func() {
				com.WriteFile(outputFileName, buf.Bytes())
			}()
			output = buf
		}
		tableWriter := tablewriter.NewWriter(output)
		csvWriter := csv.NewWriter(output)

		writeHeader := func(header []string) {
			switch outputFormat {
			case "table":
				tableWriter.SetHeader(header)
			case "csv":
				csvWriter.Write(header)
			}
		}

		writeRecord := func(row []string) {
			switch outputFormat {
			case "table":
				tableWriter.Append(row)
			case "csv":
				csvWriter.Write(row)
			}
		}

		flush := func() {
			switch outputFormat {
			case "table":
				tableWriter.Render()
			case "csv":
				csvWriter.Flush()
			}
		}

		writeHeader(evaluation.SummaryPredictDurationInformation{}.Header())
		for _, dur := range durs {
			writeRecord(dur.Row())
		}

		flush()

		return nil
	},
}

func init() {
}
