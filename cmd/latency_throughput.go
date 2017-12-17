package cmd

// import (
// 	"bytes"
// 	"encoding/csv"
// 	"os"
// 	"path"

// 	"github.com/olekukonko/tablewriter"
// 	"github.com/spf13/cobra"

// 	udb "upper.io/db.v3"
// )

// var latencyCmd = &cobra.Command{
// 	Use: "latency",
// 	Aliases: []string{
// 		"throughput",
// 	},
// 	Short: "Get evaluation latency or throughput information from CarML",
// 	RunE: func(cmd *cobra.Command, args []string) error {

// 		filter := udb.Cond{
// 			"model.name":    modelName,
// 			"model.version": modelVersion,
// 		}
// 		if machineArchitecture != "" {
// 			filter["machinearchitecture"] = machineArchitecture
// 		}
// 		if hostName != "" {
// 			filter["hostname"] = hostName
// 		}
// 		evals, err := evaluationCollection.Find(filter)
// 		if err != nil {
// 			return err
// 		}

// 		output := os.Stdout
// 		if outputFileName != nil {
// 			output = &bytes.Buffer{}
// 		}
// 		tableWriter := tablewriter.NewWriter(output)
// 		csvWriter := csv.NewWriter(output)

// 		writeHeader := func(header []string) {
// 			switch outputFormat {
// 			case "table":
// 				tableWriter.SetHeader(header)
// 			case "csv":
// 				csvWriter.Write(header)
// 			}
// 		}

// 		writeRecord := func(row []string) {
// 			switch outputFormat {
// 			case "table":
// 				tableWriter.Append(row)
// 			case "csv":
// 				csvWriter.Write(row)
// 			}
// 		}

// 		if outputFileName != nil {
// 			return comm.WriteFile(outputFileName, output.(&bytes.Buffer).Bytes())
// 		}

// 		return nil
// 	},
// }

// func init() {
// }
